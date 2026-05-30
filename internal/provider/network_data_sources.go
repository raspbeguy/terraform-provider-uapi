package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

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
			"id":        dsIDAttribute(),
			"managed":   dsManagedAttribute(),
			"device":    dsComputedString("Underlying device this interface binds to."),
			"proto":     dsComputedString("Protocol: static, dhcp, dhcpv6, pppoe, none, ppp, or wwan."),
			"ipaddr":    dsComputedString("IPv4 address."),
			"netmask":   dsComputedString("IPv4 netmask."),
			"gateway":   dsComputedString("Default gateway."),
			"dns":       dsComputedStringList("DNS servers."),
			"ip6assign": dsComputedString("IPv6 prefix assignment length."),
			"mtu":       dsComputedString("Interface MTU."),
			"auto":      dsComputedBool("Whether the interface is brought up automatically."),
		},
	}
}

func (d *networkInterfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var m networkInterfaceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := d.client.GetObject(ctx, "/"+networkInterfaceCollection+"/"+m.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading network interface", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Network interface not found", "No network interface with id "+m.ID.ValueString())
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	(&networkInterfaceResource{}).read(ctx, obj, &m, ds)
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
	obj, found, err := d.client.GetObject(ctx, "/"+networkDeviceCollection+"/"+m.ID.ValueString())
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
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
