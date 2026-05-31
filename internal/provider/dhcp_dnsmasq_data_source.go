package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &dhcpDnsmasqDataSource{}
	_ datasource.DataSourceWithConfigure = &dhcpDnsmasqDataSource{}
)

type dhcpDnsmasqDataSource struct{ client *client.Client }

func NewDhcpDnsmasqDataSource() datasource.DataSource { return &dhcpDnsmasqDataSource{} }

func (d *dhcpDnsmasqDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcp_dnsmasq"
}

func (d *dhcpDnsmasqDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *dhcpDnsmasqDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "The global dnsmasq settings (uci dhcp.dnsmasq).",
		Attributes: map[string]dsschema.Attribute{
			"id":                dsComputedString("Stable id of the dnsmasq section."),
			"managed":           dsManagedAttribute(),
			"domain":            dsComputedString("Local DNS domain."),
			"local":             dsComputedString("Local domain suffix served authoritatively."),
			"noresolv":          dsComputedBool("Whether upstream resolvers are read from resolvfile."),
			"rebind_protection": dsComputedBool("Whether DNS rebind protection is enabled."),
			"expandhosts":       dsComputedBool("Whether the local domain is added to /etc/hosts names."),
			"cachesize":         dsComputedString("DNS cache size."),
			"port":              dsComputedString("DNS service port."),
			"domainneeded":      dsComputedBool("Whether plain names without a domain are never forwarded."),
			"boguspriv":         dsComputedBool("Whether reverse lookups for private ranges are never forwarded."),
			"filterwin2k":       dsComputedBool("Whether Windows DNS queries are filtered."),
			"authoritative":     dsComputedBool("Whether dnsmasq acts as the authoritative DHCP server."),
			"readethers":        dsComputedBool("Whether static leases are read from /etc/ethers."),
			"leasefile":         dsComputedString("Path to the DHCP lease file."),
			"resolvfile":        dsComputedString("Path to the upstream resolver file."),
			"server":            dsComputedStringList("Upstream DNS servers."),
			"address":           dsComputedStringList("Static DNS address overrides."),
			"nonwildcard":       dsComputedBool("Whether dnsmasq binds only to configured interfaces."),
		},
	}
}

func (d *dhcpDnsmasqDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	obj, found, err := d.client.GetObject(ctx, dhcpDnsmasqPath)
	if err != nil {
		resp.Diagnostics.AddError("Error reading dnsmasq settings", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Dnsmasq settings not found", "The dnsmasq singleton is missing on the router")
		return
	}
	var m dhcpDnsmasqModel
	ds := newDiagsink(&resp.Diagnostics)
	(&dhcpDnsmasqResource{}).read(ctx, obj, &m, ds)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
