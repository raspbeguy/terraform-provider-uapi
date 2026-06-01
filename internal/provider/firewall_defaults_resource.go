package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const firewallDefaultsPath = "/firewall/defaults"

var (
	_ resource.Resource                = &firewallDefaultsResource{}
	_ resource.ResourceWithConfigure   = &firewallDefaultsResource{}
	_ resource.ResourceWithImportState = &firewallDefaultsResource{}
)

type firewallDefaultsResource struct {
	client *client.Client
}

func NewFirewallDefaultsResource() resource.Resource {
	return &firewallDefaultsResource{}
}

type firewallDefaultsModel struct {
	ID               types.String `tfsdk:"id"`
	Managed          types.Bool   `tfsdk:"managed"`
	ETag             types.String `tfsdk:"etag"`
	Input            types.String `tfsdk:"input"`
	Output           types.String `tfsdk:"output"`
	Forward          types.String `tfsdk:"forward"`
	SynFlood         types.Bool   `tfsdk:"syn_flood"`
	DropInvalid      types.Bool   `tfsdk:"drop_invalid"`
	SynfloodBurst    types.String `tfsdk:"synflood_burst"`
	SynfloodRate     types.String `tfsdk:"synflood_rate"`
	TCPSyncookies    types.Bool   `tfsdk:"tcp_syncookies"`
	FlowOffloading   types.Bool   `tfsdk:"flow_offloading"`
	FlowOffloadingHW types.Bool   `tfsdk:"flow_offloading_hw"`
}

func (r *firewallDefaultsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_defaults"
}

func (r *firewallDefaultsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *firewallDefaultsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Global firewall defaults (uci firewall.defaults). This is a singleton: it cannot be " +
			"created or destroyed. `terraform destroy` only removes it from state; the underlying " +
			"settings are left as-is on the router.",
		Attributes: map[string]schema.Attribute{
			"id":      computedIDAttribute(),
			"managed": managedAttribute(),
			"etag":    etagAttribute(),
			"input": schema.StringAttribute{
				Optional:    true,
				Description: "Default policy for input traffic: ACCEPT, REJECT, or DROP.",
			},
			"output": schema.StringAttribute{
				Optional:    true,
				Description: "Default policy for output traffic: ACCEPT, REJECT, or DROP.",
			},
			"forward": schema.StringAttribute{
				Optional:    true,
				Description: "Default policy for forwarded traffic: ACCEPT, REJECT, or DROP.",
			},
			"syn_flood":    optionalComputedBool("Enable SYN-flood protection. Defaults to false."),
			"drop_invalid": optionalComputedBool("Drop packets in an invalid conntrack state. Defaults to false."),
			"synflood_burst": schema.StringAttribute{
				Optional:    true,
				Description: "SYN-flood burst limit (positive integer).",
			},
			"synflood_rate": schema.StringAttribute{
				Optional:    true,
				Description: "SYN-flood rate limit (positive integer).",
			},
			"tcp_syncookies":     optionalComputedBool("Enable TCP SYN cookies. Defaults to false."),
			"flow_offloading":    optionalComputedBool("Enable software flow offloading. Defaults to false."),
			"flow_offloading_hw": optionalComputedBool("Enable hardware flow offloading. Defaults to false."),
		},
	}
}

func (r *firewallDefaultsResource) body(_ context.Context, m firewallDefaultsModel) map[string]any {
	out := map[string]any{}
	putStr(out, "input", m.Input)
	putStr(out, "output", m.Output)
	putStr(out, "forward", m.Forward)
	putBool(out, "syn_flood", m.SynFlood)
	putBool(out, "drop_invalid", m.DropInvalid)
	putStr(out, "synflood_burst", m.SynfloodBurst)
	putStr(out, "synflood_rate", m.SynfloodRate)
	putBool(out, "tcp_syncookies", m.TCPSyncookies)
	putBool(out, "flow_offloading", m.FlowOffloading)
	putBool(out, "flow_offloading_hw", m.FlowOffloadingHW)
	return out
}

func (r *firewallDefaultsResource) read(_ context.Context, obj map[string]any, m *firewallDefaultsModel) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.Input = strVal(obj, "input")
	m.Output = strVal(obj, "output")
	m.Forward = strVal(obj, "forward")
	m.SynFlood = boolVal(obj, "syn_flood")
	m.DropInvalid = boolVal(obj, "drop_invalid")
	m.SynfloodBurst = strVal(obj, "synflood_burst")
	m.SynfloodRate = strVal(obj, "synflood_rate")
	m.TCPSyncookies = boolVal(obj, "tcp_syncookies")
	m.FlowOffloading = boolVal(obj, "flow_offloading")
	m.FlowOffloadingHW = boolVal(obj, "flow_offloading_hw")
}

func (r *firewallDefaultsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan firewallDefaultsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, err := r.client.Patch(ctx, firewallDefaultsPath, r.body(ctx, plan), "")
	if err != nil {
		writeErr(&resp.Diagnostics, "configuring", "firewall defaults", err)
		return
	}
	r.read(ctx, obj, &plan)
	plan.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *firewallDefaultsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state firewallDefaultsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := r.client.GetObject(ctx, firewallDefaultsPath)
	if err != nil {
		resp.Diagnostics.AddError("Error reading firewall defaults", err.Error())
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	r.read(ctx, obj, &state)
	state.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *firewallDefaultsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state firewallDefaultsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, err := r.client.Patch(ctx, firewallDefaultsPath, r.body(ctx, plan), state.ETag.ValueString())
	if err != nil {
		writeErr(&resp.Diagnostics, "updating", "firewall defaults", err)
		return
	}
	r.read(ctx, obj, &plan)
	plan.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete is a no-op: the firewall defaults singleton cannot be removed. State is
// dropped by the framework once this returns.
func (r *firewallDefaultsResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func (r *firewallDefaultsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
