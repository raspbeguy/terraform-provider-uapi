package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const networkInterfaceCollection = "network/interfaces"

var (
	_ resource.Resource                = &networkInterfaceResource{}
	_ resource.ResourceWithConfigure   = &networkInterfaceResource{}
	_ resource.ResourceWithImportState = &networkInterfaceResource{}
)

type networkInterfaceResource struct {
	client *client.Client
}

func NewNetworkInterfaceResource() resource.Resource {
	return &networkInterfaceResource{}
}

type networkInterfaceModel struct {
	ID        types.String `tfsdk:"id"`
	Managed   types.Bool   `tfsdk:"managed"`
	ETag      types.String `tfsdk:"etag"`
	Device    types.String `tfsdk:"device"`
	Proto     types.String `tfsdk:"proto"`
	IPAddr    types.String `tfsdk:"ipaddr"`
	IPAddrs   types.List   `tfsdk:"ipaddrs"`
	Netmask   types.String `tfsdk:"netmask"`
	Gateway   types.String `tfsdk:"gateway"`
	DNS       types.List   `tfsdk:"dns"`
	IP6Assign types.String `tfsdk:"ip6assign"`
	MTU       types.String `tfsdk:"mtu"`
	Auto      types.Bool   `tfsdk:"auto"`
	// WireGuard fields (proto = "wireguard"); present in the response only then.
	PrivateKey    types.String `tfsdk:"private_key"`
	HasPrivateKey types.Bool   `tfsdk:"has_private_key"`
	ListenPort    types.String `tfsdk:"listen_port"`
	Addresses     types.List   `tfsdk:"addresses"`
	Nohostroute   types.Bool   `tfsdk:"nohostroute"`
	IP4Table      types.String `tfsdk:"ip4table"`
	IP6Table      types.String `tfsdk:"ip6table"`
	// DHCP client options (proto = "dhcp").
	PeerDNS      types.Bool   `tfsdk:"peerdns"`
	DefaultRoute types.Bool   `tfsdk:"defaultroute"`
	Metric       types.String `tfsdk:"metric"`
	Hostname     types.String `tfsdk:"hostname"`
	ClientID     types.String `tfsdk:"clientid"`
	// DHCPv6 client options (proto = "dhcpv6").
	ReqPrefix  types.String `tfsdk:"reqprefix"`
	ReqAddress types.String `tfsdk:"reqaddress"`
	IP6Hint    types.String `tfsdk:"ip6hint"`
	IP6IfaceID types.String `tfsdk:"ip6ifaceid"`
	Delegate   types.Bool   `tfsdk:"delegate"`
}

func (r *networkInterfaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_interface"
}

func (r *networkInterfaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *networkInterfaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A network interface (uci network.interface). " +
			"Caution: reloading an interface that backs your management connection can lock you out; " +
			"uapi only observes the init script exit code, not runtime convergence.",
		Attributes: map[string]schema.Attribute{
			"id":      computedIDAttribute(),
			"managed": managedAttribute(),
			"etag":    etagAttribute(),
			"device": schema.StringAttribute{
				Optional:    true,
				Description: "Underlying device this interface binds to.",
			},
			"proto": schema.StringAttribute{
				Required:    true,
				Description: "Protocol: static, dhcp, dhcpv6, pppoe, none, ppp, wwan, or wireguard.",
			},
			"ipaddr":    optionalComputedString("IPv4 address (required when proto is static, unless ipaddrs is set)."),
			"ipaddrs":   optionalComputedStringList("Full IPv4 address list for a static interface (uci `list ipaddr`). Preferred over ipaddr for multi-address interfaces."),
			"netmask":   optionalComputedString("IPv4 netmask."),
			"gateway":   optionalComputedString("Default gateway."),
			"dns":       optionalComputedStringList("DNS servers."),
			"ip6assign": optionalComputedString("IPv6 prefix assignment length."),
			"mtu":       optionalComputedString("Interface MTU."),
			"auto":      optionalComputedBool("Bring the interface up automatically. Defaults to true."),
			"private_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "WireGuard private key (proto = wireguard). Write-only: the API never returns it, so it is not refreshed from the router.",
			},
			"has_private_key": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether a WireGuard private key is configured (the key itself is never returned).",
			},
			"listen_port":  optionalComputedString("WireGuard UDP listen port (proto = wireguard)."),
			"addresses":    optionalComputedStringList("WireGuard interface addresses as CIDRs (proto = wireguard)."),
			"nohostroute":  optionalComputedBool("WireGuard: skip adding host routes for peers. Defaults to false."),
			"ip4table":     optionalComputedString("WireGuard IPv4 routing table (proto = wireguard)."),
			"ip6table":     optionalComputedString("WireGuard IPv6 routing table (proto = wireguard)."),
			"peerdns":      optionalComputedBool("Accept DNS servers advertised by the upstream (proto = dhcp or dhcpv6). Defaults to true."),
			"defaultroute": optionalComputedBool("Install the default route received over DHCP (proto = dhcp). Defaults to true."),
			"metric": schema.StringAttribute{
				Optional:    true,
				Description: "Default-route metric (proto = dhcp).",
			},
			"hostname": schema.StringAttribute{
				Optional:    true,
				Description: "Client hostname sent in DHCPDISCOVER (proto = dhcp).",
			},
			"clientid": schema.StringAttribute{
				Optional:    true,
				Description: "DHCP client identifier (proto = dhcp).",
			},
			"reqprefix": schema.StringAttribute{
				Optional:    true,
				Description: "DHCPv6 prefix-delegation request: auto, no, or a numeric prefix size (proto = dhcpv6).",
			},
			"reqaddress": schema.StringAttribute{
				Optional:    true,
				Description: "DHCPv6 IA_NA request mode: try, force, or none (proto = dhcpv6).",
			},
			"ip6hint": schema.StringAttribute{
				Optional:    true,
				Description: "Preferred IPv6 prefix hint for prefix delegation, like 2001:db8::/56 (proto = dhcpv6).",
			},
			"ip6ifaceid": schema.StringAttribute{
				Optional:    true,
				Description: "Static IPv6 interface id for IA_NA, like ::1 or an EUI-64 form (proto = dhcpv6).",
			},
			"delegate": optionalComputedBool("Accept prefix delegation downstream (proto = dhcpv6). Defaults to true."),
		},
	}
}

func (r *networkInterfaceResource) body(ctx context.Context, m networkInterfaceModel, diags *diagsink) map[string]any {
	out := map[string]any{}
	putStr(out, "device", m.Device)
	putStr(out, "proto", m.Proto)
	putStr(out, "ipaddr", m.IPAddr)
	putList(ctx, out, "ipaddrs", m.IPAddrs, diags.d)
	putStr(out, "netmask", m.Netmask)
	putStr(out, "gateway", m.Gateway)
	putList(ctx, out, "dns", m.DNS, diags.d)
	putStr(out, "ip6assign", m.IP6Assign)
	putStr(out, "mtu", m.MTU)
	putBool(out, "auto", m.Auto)
	putStr(out, "private_key", m.PrivateKey)
	putStr(out, "listen_port", m.ListenPort)
	putList(ctx, out, "addresses", m.Addresses, diags.d)
	putBool(out, "nohostroute", m.Nohostroute)
	putStr(out, "ip4table", m.IP4Table)
	putStr(out, "ip6table", m.IP6Table)
	putBool(out, "peerdns", m.PeerDNS)
	putBool(out, "defaultroute", m.DefaultRoute)
	putStr(out, "metric", m.Metric)
	putStr(out, "hostname", m.Hostname)
	putStr(out, "clientid", m.ClientID)
	putStr(out, "reqprefix", m.ReqPrefix)
	putStr(out, "reqaddress", m.ReqAddress)
	putStr(out, "ip6hint", m.IP6Hint)
	putStr(out, "ip6ifaceid", m.IP6IfaceID)
	putBool(out, "delegate", m.Delegate)
	return out
}

func (r *networkInterfaceResource) read(ctx context.Context, obj map[string]any, m *networkInterfaceModel, diags *diagsink) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.Device = strVal(obj, "device")
	m.Proto = strVal(obj, "proto")
	m.IPAddr = strVal(obj, "ipaddr")
	m.IPAddrs = diags.list(listVal(ctx, obj, "ipaddrs"))
	m.Netmask = strVal(obj, "netmask")
	m.Gateway = strVal(obj, "gateway")
	m.DNS = diags.list(listVal(ctx, obj, "dns"))
	m.IP6Assign = strVal(obj, "ip6assign")
	m.MTU = strVal(obj, "mtu")
	m.Auto = boolVal(obj, "auto")
	// private_key is write-only: leave m.PrivateKey untouched (preserve planned value).
	m.ListenPort = strVal(obj, "listen_port")
	m.Addresses = diags.list(listVal(ctx, obj, "addresses"))
	m.Nohostroute = boolVal(obj, "nohostroute")
	m.IP4Table = strVal(obj, "ip4table")
	m.IP6Table = strVal(obj, "ip6table")
	hasKey := boolVal(obj, "has_private_key")
	if hasKey.IsNull() {
		hasKey = types.BoolValue(false)
	}
	m.HasPrivateKey = hasKey
	m.PeerDNS = boolVal(obj, "peerdns")
	m.DefaultRoute = boolVal(obj, "defaultroute")
	m.Metric = strVal(obj, "metric")
	m.Hostname = strVal(obj, "hostname")
	m.ClientID = strVal(obj, "clientid")
	m.ReqPrefix = strVal(obj, "reqprefix")
	m.ReqAddress = strVal(obj, "reqaddress")
	m.IP6Hint = strVal(obj, "ip6hint")
	m.IP6IfaceID = strVal(obj, "ip6ifaceid")
	m.Delegate = boolVal(obj, "delegate")
}

func (r *networkInterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan networkInterfaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	body := r.body(ctx, plan, ds)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, err := r.client.Post(ctx, "/"+networkInterfaceCollection, body, "")
	if err != nil {
		writeErr(&resp.Diagnostics, "creating", "network interface", err)
		return
	}
	r.read(ctx, obj, &plan, ds)
	plan.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkInterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state networkInterfaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := r.client.GetObject(ctx, "/"+networkInterfaceCollection+"/"+state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading network interface", err.Error())
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

func (r *networkInterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state networkInterfaceModel
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
	obj, etag, err := r.client.Put(ctx, "/"+networkInterfaceCollection+"/"+plan.ID.ValueString(), body, state.ETag.ValueString())
	if err != nil {
		writeErr(&resp.Diagnostics, "updating", "network interface", err)
		return
	}
	r.read(ctx, obj, &plan, ds)
	plan.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkInterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state networkInterfaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(ctx, "/"+networkInterfaceCollection+"/"+state.ID.ValueString(), state.ETag.ValueString()); err != nil {
		writeErr(&resp.Diagnostics, "deleting", "network interface", err)
	}
}

func (r *networkInterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByID(ctx, r.client, networkInterfaceCollection, "network interface", req, resp)
}
