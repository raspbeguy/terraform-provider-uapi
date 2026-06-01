package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const networkRouteCollection = "network/routes"

var (
	_ resource.Resource                = &networkRouteResource{}
	_ resource.ResourceWithConfigure   = &networkRouteResource{}
	_ resource.ResourceWithImportState = &networkRouteResource{}
)

type networkRouteResource struct {
	client *client.Client
}

func NewNetworkRouteResource() resource.Resource {
	return &networkRouteResource{}
}

type networkRouteModel struct {
	ID        types.String `tfsdk:"id"`
	Managed   types.Bool   `tfsdk:"managed"`
	ETag      types.String `tfsdk:"etag"`
	Interface types.String `tfsdk:"interface"`
	Target    types.String `tfsdk:"target"`
	Netmask   types.String `tfsdk:"netmask"`
	Gateway   types.String `tfsdk:"gateway"`
	Table     types.String `tfsdk:"table"`
	Metric    types.String `tfsdk:"metric"`
	MTU       types.String `tfsdk:"mtu"`
	Source    types.String `tfsdk:"source"`
	Type      types.String `tfsdk:"type"`
}

func (r *networkRouteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_route"
}

func (r *networkRouteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *networkRouteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A static network route (uci network.route).",
		Attributes: map[string]schema.Attribute{
			"id":      computedIDAttribute(),
			"managed": managedAttribute(),
			"etag":    etagAttribute(),
			"interface": schema.StringAttribute{
				Optional:    true,
				Description: "Logical network interface the route is bound to. Must reference an existing interface for unicast routes.",
			},
			"target": schema.StringAttribute{
				Required:    true,
				Description: "Destination IPv4 address or CIDR.",
			},
			"netmask": schema.StringAttribute{
				Optional:    true,
				Description: "Destination netmask, when target is given as a bare address.",
			},
			"gateway": schema.StringAttribute{
				Optional:    true,
				Description: "Next-hop gateway IPv4 address.",
			},
			"table": schema.StringAttribute{
				Optional:    true,
				Description: "Routing table to install the route into.",
			},
			"metric": schema.StringAttribute{
				Optional:    true,
				Description: "Route metric (priority); non-negative integer.",
			},
			"mtu": schema.StringAttribute{
				Optional:    true,
				Description: "Path MTU for the route; non-negative integer.",
			},
			"source": schema.StringAttribute{
				Optional:    true,
				Description: "Preferred source IPv4 address or CIDR.",
			},
			"type": optionalComputedString("Route type: unicast, blackhole, unreachable, prohibit, throw, anycast, multicast, local, or broadcast. Defaults to unicast."),
		},
	}
}

func (r *networkRouteResource) body(_ context.Context, m networkRouteModel) map[string]any {
	out := map[string]any{}
	putStr(out, "interface", m.Interface)
	putStr(out, "target", m.Target)
	putStr(out, "netmask", m.Netmask)
	putStr(out, "gateway", m.Gateway)
	putStr(out, "table", m.Table)
	putStr(out, "metric", m.Metric)
	putStr(out, "mtu", m.MTU)
	putStr(out, "source", m.Source)
	putStr(out, "type", m.Type)
	return out
}

func (r *networkRouteResource) read(_ context.Context, obj map[string]any, m *networkRouteModel) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.Interface = strVal(obj, "interface")
	m.Target = strVal(obj, "target")
	m.Netmask = strVal(obj, "netmask")
	m.Gateway = strVal(obj, "gateway")
	m.Table = strVal(obj, "table")
	m.Metric = strVal(obj, "metric")
	m.MTU = strVal(obj, "mtu")
	m.Source = strVal(obj, "source")
	m.Type = strVal(obj, "type")
}

func (r *networkRouteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan networkRouteModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, err := r.client.Post(ctx, "/"+networkRouteCollection, r.body(ctx, plan), "")
	if err != nil {
		writeErr(&resp.Diagnostics, "creating", "network route", err)
		return
	}
	r.read(ctx, obj, &plan)
	plan.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkRouteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state networkRouteModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := r.client.GetObject(ctx, "/"+networkRouteCollection+"/"+state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading network route", err.Error())
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

func (r *networkRouteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state networkRouteModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, err := r.client.Put(ctx, "/"+networkRouteCollection+"/"+plan.ID.ValueString(), r.body(ctx, plan), state.ETag.ValueString())
	if err != nil {
		writeErr(&resp.Diagnostics, "updating", "network route", err)
		return
	}
	r.read(ctx, obj, &plan)
	plan.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkRouteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state networkRouteModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(ctx, "/"+networkRouteCollection+"/"+state.ID.ValueString(), state.ETag.ValueString()); err != nil {
		writeErr(&resp.Diagnostics, "deleting", "network route", err)
	}
}

func (r *networkRouteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByID(ctx, r.client, networkRouteCollection, "network route", req, resp)
}
