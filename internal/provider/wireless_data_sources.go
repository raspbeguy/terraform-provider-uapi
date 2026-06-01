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
	_ datasource.DataSource              = &wirelessDeviceDataSource{}
	_ datasource.DataSourceWithConfigure = &wirelessDeviceDataSource{}
)

type wirelessDeviceDataSource struct{ client *client.Client }

func NewWirelessDeviceDataSource() datasource.DataSource { return &wirelessDeviceDataSource{} }

func (d *wirelessDeviceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wireless_device"
}

func (d *wirelessDeviceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *wirelessDeviceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up a wireless radio by id.",
		Attributes: map[string]dsschema.Attribute{
			"id":       dsIDAttribute(),
			"managed":  dsManagedAttribute(),
			"etag":     dsComputedString("Opaque ETag of the resource's current state."),
			"type":     dsComputedString("Driver type: mac80211 or broadcom."),
			"band":     dsComputedString("Frequency band: 2g, 5g, 6g, or 60g."),
			"channel":  dsComputedString("Channel number or 'auto'."),
			"htmode":   dsComputedString("HT/VHT/HE mode (e.g. HT20, VHT80)."),
			"country":  dsComputedString("Regulatory country code."),
			"txpower":  dsComputedString("Transmit power in dBm."),
			"disabled": dsComputedBool("Whether the radio is disabled."),
		},
	}
}

func (d *wirelessDeviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var m wirelessDeviceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := d.client.GetObject(ctx, "/"+wirelessDeviceCollection+"/"+m.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading wireless device", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Wireless device not found", "No wireless device with id "+m.ID.ValueString())
		return
	}
	(&wirelessDeviceResource{}).read(ctx, obj, &m)
	m.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}

var (
	_ datasource.DataSource              = &wirelessInterfaceDataSource{}
	_ datasource.DataSourceWithConfigure = &wirelessInterfaceDataSource{}
)

type wirelessInterfaceDataSource struct{ client *client.Client }

func NewWirelessInterfaceDataSource() datasource.DataSource { return &wirelessInterfaceDataSource{} }

func (d *wirelessInterfaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wireless_interface"
}

func (d *wirelessInterfaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *wirelessInterfaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up a wireless interface (SSID) by id. The encryption key is never returned.",
		Attributes: map[string]dsschema.Attribute{
			"id":         dsIDAttribute(),
			"managed":    dsManagedAttribute(),
			"etag":       dsComputedString("Opaque ETag of the resource's current state."),
			"device":     dsComputedString("Wireless radio id this interface belongs to."),
			"network":    dsComputedString("Network interface this SSID is bridged to."),
			"mode":       dsComputedString("Operating mode: ap, sta, adhoc, wds, monitor, or mesh."),
			"ssid":       dsComputedString("Network name (SSID)."),
			"encryption": dsComputedString("Encryption suite (none, psk2, sae, wpa3, ...)."),
			"disabled":   dsComputedBool("Whether this SSID is disabled."),
			"hidden":     dsComputedBool("Whether the SSID is hidden."),
			"isolate":    dsComputedBool("Whether clients on this SSID are isolated."),
			"key":        dsschema.StringAttribute{Computed: true, Sensitive: true, Description: "Always null; the API never returns the key."},
			"has_key":    dsComputedBool("Whether a key is configured on the router."),
			"runtime":    wirelessInterfaceRuntimeAttribute(),
		},
	}
}

func (d *wirelessInterfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var id types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := d.client.GetObject(ctx, "/"+wirelessInterfaceCollection+"/"+id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading wireless interface", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Wireless interface not found", "No wireless interface with id "+id.ValueString())
		return
	}
	var base wirelessInterfaceModel
	(&wirelessInterfaceResource{}).read(ctx, obj, &base)
	base.ETag = types.StringValue(etag)
	m := wirelessInterfaceDS(base, obj)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
