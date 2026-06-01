package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &dhcpOdhcpdDataSource{}
	_ datasource.DataSourceWithConfigure = &dhcpOdhcpdDataSource{}
)

type dhcpOdhcpdDataSource struct{ client *client.Client }

func NewDhcpOdhcpdDataSource() datasource.DataSource { return &dhcpOdhcpdDataSource{} }

func (d *dhcpOdhcpdDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcp_odhcpd"
}

func (d *dhcpOdhcpdDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *dhcpOdhcpdDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "The global odhcpd settings (uci dhcp.odhcpd).",
		Attributes: map[string]dsschema.Attribute{
			"id":           dsComputedString("Stable id of the odhcpd section."),
			"managed":      dsManagedAttribute(),
			"etag":         dsComputedString("Opaque ETag of the resource's current state."),
			"maindhcp":     dsComputedBool("Whether odhcpd serves IPv4 DHCP instead of dnsmasq."),
			"leasefile":    dsComputedString("Path to the odhcpd lease file."),
			"leasetrigger": dsComputedString("Script run when leases change."),
			"loglevel":     dsComputedString("Syslog log level (0-7)."),
		},
	}
}

func (d *dhcpOdhcpdDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	obj, etag, found, err := d.client.GetObject(ctx, dhcpOdhcpdPath)
	if err != nil {
		resp.Diagnostics.AddError("Error reading odhcpd settings", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Odhcpd settings not found", "The odhcpd singleton is missing on the router")
		return
	}
	var m dhcpOdhcpdModel
	(&dhcpOdhcpdResource{}).read(ctx, obj, &m)
	m.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
