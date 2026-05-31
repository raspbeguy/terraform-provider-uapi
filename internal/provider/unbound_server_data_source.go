package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &unboundServerDataSource{}
	_ datasource.DataSourceWithConfigure = &unboundServerDataSource{}
)

type unboundServerDataSource struct{ client *client.Client }

func NewUnboundServerDataSource() datasource.DataSource { return &unboundServerDataSource{} }

func (d *unboundServerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_unbound_server"
}

func (d *unboundServerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *unboundServerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "The global unbound resolver settings (uci unbound.unbound).",
		Attributes: map[string]dsschema.Attribute{
			"id":             dsComputedString("Stable id of the unbound server section."),
			"managed":        dsManagedAttribute(),
			"enabled":        dsComputedBool("Whether the unbound resolver is enabled."),
			"listen_port":    dsComputedString("Port unbound listens on for DNS queries."),
			"dhcp_link":      dsComputedString("DHCP integration source: none, odhcpd, or dnsmasq."),
			"add_local_fqdn": dsComputedString("How aggressively LAN host FQDN records are added."),
			"add_wan_fqdn":   dsComputedString("How aggressively WAN host FQDN records are added."),
			"dnssec_enabled": dsComputedBool("Whether DNSSEC validation is enabled."),
			"recursion":      dsComputedString("Recursion tuning preset: default, passive, or aggressive."),
			"resource":       dsComputedString("Memory/cache sizing preset: tiny, small, medium, large, big, or huge."),
			"protocol":       dsComputedString("IP protocol mode: auto, ip4_only, ip6_only, or mixed."),
			"query_minimize": dsComputedBool("Whether QNAME minimization is enabled."),
			"prefetch":       dsComputedBool("Whether cache prefetching is enabled."),
		},
	}
}

func (d *unboundServerDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	obj, found, err := d.client.GetObject(ctx, unboundServerPath)
	if err != nil {
		resp.Diagnostics.AddError("Error reading unbound settings", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Unbound settings not found", "The unbound server singleton is missing on the router")
		return
	}
	var m unboundServerModel
	(&unboundServerResource{}).read(ctx, obj, &m)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
