package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &dhcpServerDataSource{}
	_ datasource.DataSourceWithConfigure = &dhcpServerDataSource{}
)

type dhcpServerDataSource struct{ client *client.Client }

func NewDhcpServerDataSource() datasource.DataSource { return &dhcpServerDataSource{} }

func (d *dhcpServerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcp_server"
}

func (d *dhcpServerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *dhcpServerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up a DHCP server pool by id.",
		Attributes: map[string]dsschema.Attribute{
			"id":          dsIDAttribute(),
			"managed":     dsManagedAttribute(),
			"etag":        dsComputedString("Opaque ETag of the resource's current state."),
			"interface":   dsComputedString("Network interface this pool serves."),
			"start":       dsComputedString("Pool start offset within the /24."),
			"limit":       dsComputedString("Pool size within the /24."),
			"leasetime":   dsComputedString("Lease time, e.g. 12h, 30m, 1d, or plain seconds."),
			"ignore":      dsComputedBool("Whether DHCP is disabled on this interface."),
			"force":       dsComputedBool("Whether DHCP is served even if another server is detected."),
			"dynamicdhcp": dsComputedBool("Whether dynamic leases are handed out."),
			"ra":          dsComputedString("Router advertisement mode: disabled, server, relay, or hybrid."),
			"dhcpv6":      dsComputedString("DHCPv6 mode: disabled, server, relay, or hybrid."),
			"ra_default":  dsComputedString("Default router lifetime behavior for router advertisements."),
			"domain":      dsComputedString("DNS domain announced to clients on this interface."),
			"dhcp_option": dsComputedStringList("Raw dnsmasq DHCP options for this pool."),
		},
	}
}

func (d *dhcpServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var m dhcpServerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := d.client.GetObject(ctx, "/"+dhcpServerCollection+"/"+m.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading dhcp server", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("DHCP server not found", "No dhcp server with id "+m.ID.ValueString())
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	(&dhcpServerResource{}).read(ctx, obj, &m, ds)
	m.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
