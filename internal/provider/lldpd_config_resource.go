package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const lldpdConfigPath = "/lldpd/config"

var (
	_ resource.Resource                = &lldpdConfigResource{}
	_ resource.ResourceWithConfigure   = &lldpdConfigResource{}
	_ resource.ResourceWithImportState = &lldpdConfigResource{}
)

type lldpdConfigResource struct {
	client *client.Client
}

func NewLldpdConfigResource() resource.Resource {
	return &lldpdConfigResource{}
}

type lldpdConfigModel struct {
	ID               types.String `tfsdk:"id"`
	Managed          types.Bool   `tfsdk:"managed"`
	EnableCDP        types.Bool   `tfsdk:"enable_cdp"`
	EnableFDP        types.Bool   `tfsdk:"enable_fdp"`
	EnableSONMP      types.Bool   `tfsdk:"enable_sonmp"`
	EnableEDP        types.Bool   `tfsdk:"enable_edp"`
	EnableLLDPMED    types.Bool   `tfsdk:"enable_lldpmed"`
	LLDPClass        types.String `tfsdk:"lldp_class"`
	LLDPDescription  types.Bool   `tfsdk:"lldp_description"`
	LLDPCapabilities types.Bool   `tfsdk:"lldp_capabilities"`
	LLDPMgmtIP       types.String `tfsdk:"lldp_mgmt_ip"`
	Interface        types.List   `tfsdk:"interface"`
}

func (r *lldpdConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lldpd_config"
}

func (r *lldpdConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *lldpdConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Global lldpd daemon settings (uci lldpd.lldpd). This is a singleton: it cannot be " +
			"created or destroyed. `terraform destroy` only removes it from state; the underlying " +
			"settings are left as-is on the router.",
		Attributes: map[string]schema.Attribute{
			"id":                computedIDAttribute(),
			"managed":           managedAttribute(),
			"enable_cdp":        optionalComputedBool("Advertise via Cisco Discovery Protocol. Defaults to false."),
			"enable_fdp":        optionalComputedBool("Advertise via Foundry Discovery Protocol. Defaults to false."),
			"enable_sonmp":      optionalComputedBool("Advertise via Nortel SONMP. Defaults to false."),
			"enable_edp":        optionalComputedBool("Advertise via Extreme Discovery Protocol. Defaults to false."),
			"enable_lldpmed":    optionalComputedBool("Advertise via LLDP-MED. Defaults to false."),
			"lldp_class":        schema.StringAttribute{Optional: true, Description: "LLDP-MED device class (1-4)."},
			"lldp_description":  optionalComputedBool("Advertise the system description. Defaults to true."),
			"lldp_capabilities": optionalComputedBool("Advertise system capabilities. Defaults to true."),
			"lldp_mgmt_ip":      schema.StringAttribute{Optional: true, Description: "Management IP address to advertise."},
			"interface":         optionalComputedStringList("Network interfaces lldpd listens on. Empty means all interfaces."),
		},
	}
}

func (r *lldpdConfigResource) body(ctx context.Context, m lldpdConfigModel, diags *diagsink) map[string]any {
	out := map[string]any{}
	putBool(out, "enable_cdp", m.EnableCDP)
	putBool(out, "enable_fdp", m.EnableFDP)
	putBool(out, "enable_sonmp", m.EnableSONMP)
	putBool(out, "enable_edp", m.EnableEDP)
	putBool(out, "enable_lldpmed", m.EnableLLDPMED)
	putStr(out, "lldp_class", m.LLDPClass)
	putBool(out, "lldp_description", m.LLDPDescription)
	putBool(out, "lldp_capabilities", m.LLDPCapabilities)
	putStr(out, "lldp_mgmt_ip", m.LLDPMgmtIP)
	putList(ctx, out, "interface", m.Interface, diags.d)
	return out
}

func (r *lldpdConfigResource) read(ctx context.Context, obj map[string]any, m *lldpdConfigModel, diags *diagsink) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.EnableCDP = boolVal(obj, "enable_cdp")
	m.EnableFDP = boolVal(obj, "enable_fdp")
	m.EnableSONMP = boolVal(obj, "enable_sonmp")
	m.EnableEDP = boolVal(obj, "enable_edp")
	m.EnableLLDPMED = boolVal(obj, "enable_lldpmed")
	m.LLDPClass = strVal(obj, "lldp_class")
	m.LLDPDescription = boolVal(obj, "lldp_description")
	m.LLDPCapabilities = boolVal(obj, "lldp_capabilities")
	m.LLDPMgmtIP = strVal(obj, "lldp_mgmt_ip")
	m.Interface = diags.list(listVal(ctx, obj, "interface"))
}

func (r *lldpdConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan lldpdConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	body := r.body(ctx, plan, ds)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Patch(ctx, lldpdConfigPath, body)
	if err != nil {
		resp.Diagnostics.AddError("Error configuring lldpd settings", err.Error())
		return
	}
	r.read(ctx, obj, &plan, ds)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *lldpdConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state lldpdConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := r.client.GetObject(ctx, lldpdConfigPath)
	if err != nil {
		resp.Diagnostics.AddError("Error reading lldpd settings", err.Error())
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	r.read(ctx, obj, &state, ds)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *lldpdConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan lldpdConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	body := r.body(ctx, plan, ds)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Patch(ctx, lldpdConfigPath, body)
	if err != nil {
		resp.Diagnostics.AddError("Error updating lldpd settings", err.Error())
		return
	}
	r.read(ctx, obj, &plan, ds)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete is a no-op: the lldpd config singleton cannot be removed. State is
// dropped by the framework once this returns.
func (r *lldpdConfigResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func (r *lldpdConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
