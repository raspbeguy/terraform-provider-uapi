package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const networkWireguardPeerCollection = "network/wireguard_peers"

var (
	_ resource.Resource                = &networkWireguardPeerResource{}
	_ resource.ResourceWithConfigure   = &networkWireguardPeerResource{}
	_ resource.ResourceWithImportState = &networkWireguardPeerResource{}
)

type networkWireguardPeerResource struct {
	client *client.Client
}

func NewNetworkWireguardPeerResource() resource.Resource {
	return &networkWireguardPeerResource{}
}

type networkWireguardPeerModel struct {
	ID                  types.String `tfsdk:"id"`
	Managed             types.Bool   `tfsdk:"managed"`
	ETag                types.String `tfsdk:"etag"`
	Interface           types.String `tfsdk:"interface"`
	Description         types.String `tfsdk:"description"`
	PublicKey           types.String `tfsdk:"public_key"`
	PresharedKey        types.String `tfsdk:"preshared_key"`
	HasPresharedKey     types.Bool   `tfsdk:"has_preshared_key"`
	AllowedIPs          types.List   `tfsdk:"allowed_ips"`
	EndpointHost        types.String `tfsdk:"endpoint_host"`
	EndpointPort        types.String `tfsdk:"endpoint_port"`
	PersistentKeepalive types.String `tfsdk:"persistent_keepalive"`
	RouteAllowedIPs     types.Bool   `tfsdk:"route_allowed_ips"`
	Disabled            types.Bool   `tfsdk:"disabled"`
}

func (r *networkWireguardPeerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_wireguard_peer"
}

func (r *networkWireguardPeerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *networkWireguardPeerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A WireGuard peer attached to a WireGuard network interface (uci network.wireguard_<interface>).",
		Attributes: map[string]schema.Attribute{
			"id":      computedIDAttribute(),
			"managed": managedAttribute(),
			"etag":    etagAttribute(),
			"interface": schema.StringAttribute{
				Required:    true,
				Description: "Parent WireGuard interface name. Must reference an existing interface with proto=wireguard.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Human-readable peer description.",
			},
			"public_key": schema.StringAttribute{
				Required:    true,
				Description: "Peer public key (44-char base64).",
			},
			"preshared_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Optional preshared key (44-char base64). Write-only: the API never returns it, so it is not refreshed from the router; rely on has_preshared_key.",
			},
			"has_preshared_key": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether a preshared key is configured on the router (the value is never returned).",
			},
			"allowed_ips": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "Non-empty list of IPv4 CIDRs routed to this peer.",
			},
			"endpoint_host": schema.StringAttribute{
				Optional:    true,
				Description: "Peer endpoint hostname or IP address.",
			},
			"endpoint_port": schema.StringAttribute{
				Optional:    true,
				Description: "Peer endpoint UDP port (1-65535).",
			},
			"persistent_keepalive": schema.StringAttribute{
				Optional:    true,
				Description: "Keepalive interval in seconds (0-65535).",
			},
			"route_allowed_ips": optionalComputedBool("Automatically create routes for allowed_ips. Defaults to false."),
			"disabled":          optionalComputedBool("Disable this peer. Defaults to false."),
		},
	}
}

func (r *networkWireguardPeerResource) body(ctx context.Context, m networkWireguardPeerModel, diags *diagsink) map[string]any {
	out := map[string]any{}
	putStr(out, "interface", m.Interface)
	putStr(out, "description", m.Description)
	putStr(out, "public_key", m.PublicKey)
	putStr(out, "preshared_key", m.PresharedKey)
	putList(ctx, out, "allowed_ips", m.AllowedIPs, diags.d)
	putStr(out, "endpoint_host", m.EndpointHost)
	putStr(out, "endpoint_port", m.EndpointPort)
	putStr(out, "persistent_keepalive", m.PersistentKeepalive)
	putBool(out, "route_allowed_ips", m.RouteAllowedIPs)
	putBool(out, "disabled", m.Disabled)
	return out
}

// read refreshes everything except preshared_key, which the API never returns
// and which must therefore be preserved from the prior plan/state value.
func (r *networkWireguardPeerResource) read(ctx context.Context, obj map[string]any, m *networkWireguardPeerModel, diags *diagsink) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.Interface = strVal(obj, "interface")
	m.Description = strVal(obj, "description")
	m.PublicKey = strVal(obj, "public_key")
	hasPSK := boolVal(obj, "has_preshared_key")
	if hasPSK.IsNull() {
		hasPSK = types.BoolValue(false)
	}
	m.HasPresharedKey = hasPSK
	m.AllowedIPs = diags.list(listVal(ctx, obj, "allowed_ips"))
	m.EndpointHost = strVal(obj, "endpoint_host")
	m.EndpointPort = strVal(obj, "endpoint_port")
	m.PersistentKeepalive = strVal(obj, "persistent_keepalive")
	m.RouteAllowedIPs = boolVal(obj, "route_allowed_ips")
	m.Disabled = boolVal(obj, "disabled")
}

func (r *networkWireguardPeerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan networkWireguardPeerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	body := r.body(ctx, plan, ds)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, err := r.client.Post(ctx, "/"+networkWireguardPeerCollection, body, "")
	if err != nil {
		writeErr(&resp.Diagnostics, "creating", "network WireGuard peer", err)
		return
	}
	r.read(ctx, obj, &plan, ds)
	plan.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkWireguardPeerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state networkWireguardPeerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := r.client.GetObject(ctx, "/"+networkWireguardPeerCollection+"/"+state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading network WireGuard peer", err.Error())
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	r.read(ctx, obj, &state, ds)
	state.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *networkWireguardPeerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state networkWireguardPeerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	body := r.body(ctx, plan, ds)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, err := r.client.Put(ctx, "/"+networkWireguardPeerCollection+"/"+plan.ID.ValueString(), body, state.ETag.ValueString())
	if err != nil {
		writeErr(&resp.Diagnostics, "updating", "network WireGuard peer", err)
		return
	}
	r.read(ctx, obj, &plan, ds)
	plan.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkWireguardPeerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state networkWireguardPeerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(ctx, "/"+networkWireguardPeerCollection+"/"+state.ID.ValueString(), state.ETag.ValueString()); err != nil {
		writeErr(&resp.Diagnostics, "deleting", "network WireGuard peer", err)
	}
}

func (r *networkWireguardPeerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByID(ctx, r.client, networkWireguardPeerCollection, "network WireGuard peer", req, resp)
}
