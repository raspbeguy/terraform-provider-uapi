package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const dhcpServerCollection = "dhcp/servers"

var (
	_ resource.Resource                = &dhcpServerResource{}
	_ resource.ResourceWithConfigure   = &dhcpServerResource{}
	_ resource.ResourceWithImportState = &dhcpServerResource{}
)

type dhcpServerResource struct {
	client *client.Client
}

func NewDhcpServerResource() resource.Resource {
	return &dhcpServerResource{}
}

type dhcpServerModel struct {
	ID          types.String `tfsdk:"id"`
	Managed     types.Bool   `tfsdk:"managed"`
	Interface   types.String `tfsdk:"interface"`
	Start       types.String `tfsdk:"start"`
	Limit       types.String `tfsdk:"limit"`
	Leasetime   types.String `tfsdk:"leasetime"`
	Ignore      types.Bool   `tfsdk:"ignore"`
	Force       types.Bool   `tfsdk:"force"`
	Dynamicdhcp types.Bool   `tfsdk:"dynamicdhcp"`
	RA          types.String `tfsdk:"ra"`
	DHCPv6      types.String `tfsdk:"dhcpv6"`
	RADefault   types.String `tfsdk:"ra_default"`
	Domain      types.String `tfsdk:"domain"`
	DHCPOption  types.List   `tfsdk:"dhcp_option"`
}

func (r *dhcpServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcp_server"
}

func (r *dhcpServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *dhcpServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A per-interface DHCP pool (uci dhcp.dhcp).",
		Attributes: map[string]schema.Attribute{
			"id":      computedIDAttribute(),
			"managed": managedAttribute(),
			"interface": schema.StringAttribute{
				Required:    true,
				Description: "Network interface this pool serves.",
			},
			"start": schema.StringAttribute{
				Optional:    true,
				Description: "Pool start offset within the /24 (0-254).",
			},
			"limit": schema.StringAttribute{
				Optional:    true,
				Description: "Pool size within the /24 (0-254).",
			},
			"leasetime": schema.StringAttribute{
				Optional:    true,
				Description: "Lease time, e.g. 12h, 30m, 1d, or plain seconds.",
			},
			"ignore":      optionalComputedBool("Ignore this interface (disable DHCP). Defaults to false."),
			"force":       optionalComputedBool("Serve DHCP even if another server is detected. Defaults to false."),
			"dynamicdhcp": optionalComputedBool("Hand out dynamic leases. Defaults to true."),
			"ra": schema.StringAttribute{
				Optional:    true,
				Description: "Router advertisement mode: disabled, server, relay, or hybrid.",
			},
			"dhcpv6": schema.StringAttribute{
				Optional:    true,
				Description: "DHCPv6 mode: disabled, server, relay, or hybrid.",
			},
			"ra_default": schema.StringAttribute{
				Optional:    true,
				Description: "Default router lifetime behavior for router advertisements.",
			},
			"domain": schema.StringAttribute{
				Optional:    true,
				Description: "DNS domain announced to clients on this interface.",
			},
			"dhcp_option": optionalComputedStringList("Raw dnsmasq DHCP options for this pool."),
		},
	}
}

func (r *dhcpServerResource) body(ctx context.Context, m dhcpServerModel, diags *diagsink) map[string]any {
	out := map[string]any{}
	putStr(out, "interface", m.Interface)
	putStr(out, "start", m.Start)
	putStr(out, "limit", m.Limit)
	putStr(out, "leasetime", m.Leasetime)
	putBool(out, "ignore", m.Ignore)
	putBool(out, "force", m.Force)
	putBool(out, "dynamicdhcp", m.Dynamicdhcp)
	putStr(out, "ra", m.RA)
	putStr(out, "dhcpv6", m.DHCPv6)
	putStr(out, "ra_default", m.RADefault)
	putStr(out, "domain", m.Domain)
	putList(ctx, out, "dhcp_option", m.DHCPOption, diags.d)
	return out
}

func (r *dhcpServerResource) read(ctx context.Context, obj map[string]any, m *dhcpServerModel, diags *diagsink) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.Interface = strVal(obj, "interface")
	m.Start = strVal(obj, "start")
	m.Limit = strVal(obj, "limit")
	m.Leasetime = strVal(obj, "leasetime")
	m.Ignore = boolVal(obj, "ignore")
	m.Force = boolVal(obj, "force")
	m.Dynamicdhcp = boolVal(obj, "dynamicdhcp")
	m.RA = strVal(obj, "ra")
	m.DHCPv6 = strVal(obj, "dhcpv6")
	m.RADefault = strVal(obj, "ra_default")
	m.Domain = strVal(obj, "domain")
	m.DHCPOption = diags.list(listVal(ctx, obj, "dhcp_option"))
}

func (r *dhcpServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dhcpServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	body := r.body(ctx, plan, ds)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Post(ctx, "/"+dhcpServerCollection, body)
	if err != nil {
		resp.Diagnostics.AddError("Error creating dhcp server", err.Error())
		return
	}
	r.read(ctx, obj, &plan, ds)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dhcpServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dhcpServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := r.client.GetObject(ctx, "/"+dhcpServerCollection+"/"+state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading dhcp server", err.Error())
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

func (r *dhcpServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan dhcpServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	body := r.body(ctx, plan, ds)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Put(ctx, "/"+dhcpServerCollection+"/"+plan.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Error updating dhcp server", err.Error())
		return
	}
	r.read(ctx, obj, &plan, ds)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dhcpServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dhcpServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(ctx, "/"+dhcpServerCollection+"/"+state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting dhcp server", err.Error())
	}
}

func (r *dhcpServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByID(ctx, r.client, dhcpServerCollection, "dhcp server", req, resp)
}
