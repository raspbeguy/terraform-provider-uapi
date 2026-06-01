package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &networkInterfaceDataSource{}
	_ datasource.DataSourceWithConfigure = &networkInterfaceDataSource{}
)

type networkInterfaceDataSource struct{ client *client.Client }

func NewNetworkInterfaceDataSource() datasource.DataSource { return &networkInterfaceDataSource{} }

func (d *networkInterfaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_interface"
}

func (d *networkInterfaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *networkInterfaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up a network interface by id.",
		Attributes: map[string]dsschema.Attribute{
			"id":              dsIDAttribute(),
			"managed":         dsManagedAttribute(),
			"etag":            dsComputedString("Opaque ETag of the resource's current state."),
			"device":          dsComputedString("Underlying device this interface binds to."),
			"proto":           dsComputedString("Protocol: static, dhcp, dhcpv6, pppoe, none, ppp, or wwan."),
			"ipaddr":          dsComputedString("IPv4 address."),
			"ipaddrs":         dsComputedStringList("Full IPv4 address list for a static interface."),
			"netmask":         dsComputedString("IPv4 netmask."),
			"gateway":         dsComputedString("Default gateway."),
			"dns":             dsComputedStringList("DNS servers."),
			"ip6assign":       dsComputedString("IPv6 prefix assignment length."),
			"mtu":             dsComputedString("Interface MTU."),
			"auto":            dsComputedBool("Whether the interface is brought up automatically."),
			"private_key":     dsschema.StringAttribute{Computed: true, Sensitive: true, Description: "Always null; the WireGuard private key is never returned."},
			"has_private_key": dsComputedBool("Whether a WireGuard private key is configured."),
			"listen_port":     dsComputedString("WireGuard UDP listen port."),
			"addresses":       dsComputedStringList("WireGuard interface addresses (CIDRs)."),
			"nohostroute":     dsComputedBool("WireGuard: whether host routes for peers are skipped."),
			"ip4table":        dsComputedString("WireGuard IPv4 routing table."),
			"ip6table":        dsComputedString("WireGuard IPv6 routing table."),
			"peerdns":         dsComputedBool("Whether DNS servers advertised by the upstream are accepted (dhcp/dhcpv6)."),
			"defaultroute":    dsComputedBool("Whether the default route received over DHCP is installed (dhcp)."),
			"metric":          dsComputedString("Default-route metric (dhcp)."),
			"hostname":        dsComputedString("Client hostname sent in DHCPDISCOVER (dhcp)."),
			"clientid":        dsComputedString("DHCP client identifier (dhcp)."),
			"reqprefix":       dsComputedString("DHCPv6 prefix-delegation request: auto, no, or a numeric prefix size (dhcpv6)."),
			"reqaddress":      dsComputedString("DHCPv6 IA_NA request mode: try, force, or none (dhcpv6)."),
			"ip6hint":         dsComputedString("Preferred IPv6 prefix hint for prefix delegation (dhcpv6)."),
			"ip6ifaceid":      dsComputedString("Static IPv6 interface id for IA_NA (dhcpv6)."),
			"delegate":        dsComputedBool("Whether prefix delegation downstream is accepted (dhcpv6)."),
			"runtime":         networkInterfaceRuntimeAttribute(),
		},
	}
}

func (d *networkInterfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var id types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := d.client.GetObject(ctx, "/"+networkInterfaceCollection+"/"+id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading network interface", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Network interface not found", "No network interface with id "+id.ValueString())
		return
	}
	var base networkInterfaceModel
	ds := newDiagsink(&resp.Diagnostics)
	(&networkInterfaceResource{}).read(ctx, obj, &base, ds)
	base.ETag = types.StringValue(etag)
	m := networkInterfaceDS(base, obj)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}

var (
	_ datasource.DataSource              = &networkDeviceDataSource{}
	_ datasource.DataSourceWithConfigure = &networkDeviceDataSource{}
)

type networkDeviceDataSource struct{ client *client.Client }

func NewNetworkDeviceDataSource() datasource.DataSource { return &networkDeviceDataSource{} }

func (d *networkDeviceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_device"
}

func (d *networkDeviceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *networkDeviceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up a network device by id.",
		Attributes: map[string]dsschema.Attribute{
			"id":      dsIDAttribute(),
			"managed": dsManagedAttribute(),
			"etag":    dsComputedString("Opaque ETag of the resource's current state."),
			"name":    dsComputedString("Device name (e.g. br-lan)."),
			"type":    dsComputedString("Device type: bridge, 8021q, 8021ad, macvlan, veth, tun, or tap."),
			"ports":   dsComputedStringList("Member interfaces."),
			"vid":     dsComputedString("VLAN id."),
			"ifname":  dsComputedString("Base interface name for VLAN/macvlan devices."),
			"mtu":     dsComputedString("Device MTU."),
			"macaddr": dsComputedString("Override MAC address."),
			"ipv6":    dsComputedBool("Whether IPv6 is enabled on the device."),
		},
	}
}

func (d *networkDeviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var m networkDeviceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := d.client.GetObject(ctx, "/"+networkDeviceCollection+"/"+m.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading network device", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Network device not found", "No network device with id "+m.ID.ValueString())
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	(&networkDeviceResource{}).read(ctx, obj, &m, ds)
	m.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
