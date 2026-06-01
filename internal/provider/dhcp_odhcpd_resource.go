package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const dhcpOdhcpdPath = "/dhcp/odhcpd"

var (
	_ resource.Resource                = &dhcpOdhcpdResource{}
	_ resource.ResourceWithConfigure   = &dhcpOdhcpdResource{}
	_ resource.ResourceWithImportState = &dhcpOdhcpdResource{}
)

type dhcpOdhcpdResource struct {
	client *client.Client
}

func NewDhcpOdhcpdResource() resource.Resource {
	return &dhcpOdhcpdResource{}
}

type dhcpOdhcpdModel struct {
	ID           types.String `tfsdk:"id"`
	Managed      types.Bool   `tfsdk:"managed"`
	ETag         types.String `tfsdk:"etag"`
	Maindhcp     types.Bool   `tfsdk:"maindhcp"`
	Leasefile    types.String `tfsdk:"leasefile"`
	Leasetrigger types.String `tfsdk:"leasetrigger"`
	Loglevel     types.String `tfsdk:"loglevel"`
}

func (r *dhcpOdhcpdResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcp_odhcpd"
}

func (r *dhcpOdhcpdResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *dhcpOdhcpdResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Global odhcpd settings (uci dhcp.odhcpd). This is a singleton: it cannot be " +
			"created or destroyed. `terraform destroy` only removes it from state; the underlying " +
			"settings are left as-is on the router.",
		Attributes: map[string]schema.Attribute{
			"id":       computedIDAttribute(),
			"managed":  managedAttribute(),
			"etag":     etagAttribute(),
			"maindhcp": optionalComputedBool("Let odhcpd serve IPv4 DHCP instead of dnsmasq. Defaults to false."),
			"leasefile": schema.StringAttribute{
				Optional:    true,
				Description: "Path to the odhcpd lease file.",
			},
			"leasetrigger": schema.StringAttribute{
				Optional:    true,
				Description: "Script run when leases change.",
			},
			"loglevel": schema.StringAttribute{
				Optional:    true,
				Description: "Syslog log level (0-7).",
			},
		},
	}
}

func (r *dhcpOdhcpdResource) body(_ context.Context, m dhcpOdhcpdModel) map[string]any {
	out := map[string]any{}
	putBool(out, "maindhcp", m.Maindhcp)
	putStr(out, "leasefile", m.Leasefile)
	putStr(out, "leasetrigger", m.Leasetrigger)
	putStr(out, "loglevel", m.Loglevel)
	return out
}

func (r *dhcpOdhcpdResource) read(_ context.Context, obj map[string]any, m *dhcpOdhcpdModel) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.Maindhcp = boolVal(obj, "maindhcp")
	m.Leasefile = strVal(obj, "leasefile")
	m.Leasetrigger = strVal(obj, "leasetrigger")
	m.Loglevel = strVal(obj, "loglevel")
}

func (r *dhcpOdhcpdResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dhcpOdhcpdModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, err := r.client.Patch(ctx, dhcpOdhcpdPath, r.body(ctx, plan), "")
	if err != nil {
		writeErr(&resp.Diagnostics, "configuring", "odhcpd settings", err)
		return
	}
	r.read(ctx, obj, &plan)
	plan.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dhcpOdhcpdResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dhcpOdhcpdModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := r.client.GetObject(ctx, dhcpOdhcpdPath)
	if err != nil {
		resp.Diagnostics.AddError("Error reading odhcpd settings", err.Error())
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

func (r *dhcpOdhcpdResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state dhcpOdhcpdModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, err := r.client.Patch(ctx, dhcpOdhcpdPath, r.body(ctx, plan), state.ETag.ValueString())
	if err != nil {
		writeErr(&resp.Diagnostics, "updating", "odhcpd settings", err)
		return
	}
	r.read(ctx, obj, &plan)
	plan.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete is a no-op: the odhcpd singleton cannot be removed. State is dropped
// by the framework once this returns.
func (r *dhcpOdhcpdResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func (r *dhcpOdhcpdResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
