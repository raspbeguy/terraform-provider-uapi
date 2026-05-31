package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const snmpdSystemPath = "/snmpd/system"

var (
	_ resource.Resource                = &snmpdSystemResource{}
	_ resource.ResourceWithConfigure   = &snmpdSystemResource{}
	_ resource.ResourceWithImportState = &snmpdSystemResource{}
)

type snmpdSystemResource struct {
	client *client.Client
}

func NewSnmpdSystemResource() resource.Resource {
	return &snmpdSystemResource{}
}

type snmpdSystemModel struct {
	ID          types.String `tfsdk:"id"`
	Managed     types.Bool   `tfsdk:"managed"`
	SysLocation types.String `tfsdk:"sys_location"`
	SysContact  types.String `tfsdk:"sys_contact"`
	SysName     types.String `tfsdk:"sys_name"`
	SysServices types.String `tfsdk:"sys_services"`
	SysDescr    types.String `tfsdk:"sys_descr"`
	SysObjectID types.String `tfsdk:"sys_object_id"`
}

func (r *snmpdSystemResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snmpd_system"
}

func (r *snmpdSystemResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *snmpdSystemResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "SNMP agent system identity (uci snmpd.system). This is a singleton: it cannot be " +
			"created or destroyed. `terraform destroy` only removes it from state; the underlying " +
			"settings are left as-is on the router.",
		Attributes: map[string]schema.Attribute{
			"id":      computedIDAttribute(),
			"managed": managedAttribute(),
			"sys_location": schema.StringAttribute{
				Optional:    true,
				Description: "Value reported as sysLocation (the device's physical location).",
			},
			"sys_contact": schema.StringAttribute{
				Optional:    true,
				Description: "Value reported as sysContact (the administrative contact).",
			},
			"sys_name": schema.StringAttribute{
				Optional:    true,
				Description: "Value reported as sysName (the administratively assigned name).",
			},
			"sys_services": schema.StringAttribute{
				Optional:    true,
				Description: "Value reported as sysServices (the OSI layer services bitmask).",
			},
			"sys_descr": schema.StringAttribute{
				Optional:    true,
				Description: "Value reported as sysDescr (a textual description of the entity).",
			},
			"sys_object_id": schema.StringAttribute{
				Optional:    true,
				Description: "Value reported as sysObjectID (the vendor authoritative object identifier).",
			},
		},
	}
}

func (r *snmpdSystemResource) body(_ context.Context, m snmpdSystemModel) map[string]any {
	out := map[string]any{}
	putStr(out, "sysLocation", m.SysLocation)
	putStr(out, "sysContact", m.SysContact)
	putStr(out, "sysName", m.SysName)
	putStr(out, "sysServices", m.SysServices)
	putStr(out, "sysDescr", m.SysDescr)
	putStr(out, "sysObjectID", m.SysObjectID)
	return out
}

func (r *snmpdSystemResource) read(_ context.Context, obj map[string]any, m *snmpdSystemModel) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.SysLocation = strVal(obj, "sysLocation")
	m.SysContact = strVal(obj, "sysContact")
	m.SysName = strVal(obj, "sysName")
	m.SysServices = strVal(obj, "sysServices")
	m.SysDescr = strVal(obj, "sysDescr")
	m.SysObjectID = strVal(obj, "sysObjectID")
}

func (r *snmpdSystemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan snmpdSystemModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Patch(ctx, snmpdSystemPath, r.body(ctx, plan))
	if err != nil {
		resp.Diagnostics.AddError("Error configuring snmpd system identity", err.Error())
		return
	}
	r.read(ctx, obj, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *snmpdSystemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state snmpdSystemModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := r.client.GetObject(ctx, snmpdSystemPath)
	if err != nil {
		resp.Diagnostics.AddError("Error reading snmpd system identity", err.Error())
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	r.read(ctx, obj, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *snmpdSystemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan snmpdSystemModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Patch(ctx, snmpdSystemPath, r.body(ctx, plan))
	if err != nil {
		resp.Diagnostics.AddError("Error updating snmpd system identity", err.Error())
		return
	}
	r.read(ctx, obj, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete is a no-op: the snmpd system singleton cannot be removed. State is
// dropped by the framework once this returns.
func (r *snmpdSystemResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func (r *snmpdSystemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
