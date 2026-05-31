package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

type uapiProvider struct {
	version string
}

type providerModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
	Insecure types.Bool   `tfsdk:"insecure"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &uapiProvider{version: version}
	}
}

func (p *uapiProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "uapi"
	resp.Version = p.version
}

func (p *uapiProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage OpenWrt configuration through the uapi REST API. " +
			"Only curated endpoints are exposed; the /raw passthrough is intentionally not supported.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional: true,
				Description: "Base URL of the uapi API, including the version prefix, " +
					"e.g. https://router.example.com/api/v1. May also be set via the " +
					"UAPI_ENDPOINT or UAPI_BASE environment variable.",
			},
			"token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Bearer token created with `uapi-token create`. May also be set via the UAPI_TOKEN environment variable.",
			},
			"insecure": schema.BoolAttribute{
				Optional: true,
				Description: "Skip TLS certificate verification. Needed for uapi's default self-signed " +
					"certificate; do not use in production. May also be set via UAPI_INSECURE.",
			},
		},
	}
}

func (p *uapiProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var cfg providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := firstNonEmpty(cfg.Endpoint, os.Getenv("UAPI_ENDPOINT"), os.Getenv("UAPI_BASE"))
	token := firstNonEmpty(cfg.Token, os.Getenv("UAPI_TOKEN"))

	insecure := false
	if !cfg.Insecure.IsNull() && !cfg.Insecure.IsUnknown() {
		insecure = cfg.Insecure.ValueBool()
	} else if v := os.Getenv("UAPI_INSECURE"); v == "1" || v == "true" {
		insecure = true
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(path.Root("endpoint"),
			"Missing uapi endpoint",
			"Set the provider `endpoint` argument or the UAPI_ENDPOINT environment variable.")
	}
	if token == "" {
		resp.Diagnostics.AddAttributeError(path.Root("token"),
			"Missing uapi token",
			"Set the provider `token` argument or the UAPI_TOKEN environment variable.")
	}
	if resp.Diagnostics.HasError() {
		return
	}

	c := client.New(endpoint, token, insecure)
	resp.ResourceData = c
	resp.DataSourceData = c
}

func (p *uapiProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// 1.0 surface
		NewFirewallRuleResource,
		NewFirewallZoneResource,
		NewFirewallRedirectResource,
		NewNetworkInterfaceResource,
		NewNetworkDeviceResource,
		NewWirelessDeviceResource,
		NewWirelessInterfaceResource,
		NewDHCPHostResource,
		NewSystemResource,
		// 1.1 network
		NewNetworkRouteResource,
		NewNetworkRuleResource,
		NewNetworkBridgeVlanResource,
		NewNetworkWireguardPeerResource,
		// 1.1 firewall
		NewFirewallForwardingResource,
		NewFirewallDefaultsResource,
		// 1.1 dhcp
		NewDhcpServerResource,
		NewDhcpDnsmasqResource,
		NewDhcpOdhcpdResource,
		// 1.1 snmpd
		NewSnmpdAccessResource,
		NewSnmpdAgentResource,
		NewSnmpdCom2secResource,
		NewSnmpdGroupResource,
		NewSnmpdSystemResource,
		// 1.1 uhttpd / dropbear
		NewUhttpdCertResource,
		NewUhttpdInstanceResource,
		NewDropbearInstanceResource,
		// 1.1 system / sqm
		NewSystemTimeserverResource,
		NewSqmQueueResource,
		// 1.1 vnstat + misc singletons
		NewVnstatInterfaceResource,
		NewVnstatConfigResource,
		NewLldpdConfigResource,
		NewPrometheusNodeExporterLuaConfigResource,
		NewUnboundServerResource,
		// 1.1 packages
		NewPackageResource,
		NewPackageFeedResource,
	}
}

func (p *uapiProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// 1.0 surface
		NewFirewallRuleDataSource,
		NewFirewallZoneDataSource,
		NewFirewallRedirectDataSource,
		NewNetworkInterfaceDataSource,
		NewNetworkDeviceDataSource,
		NewWirelessDeviceDataSource,
		NewWirelessInterfaceDataSource,
		NewDHCPHostDataSource,
		NewSystemDataSource,
		NewDHCPLeasesDataSource,
		// 1.1 network
		NewNetworkRouteDataSource,
		NewNetworkRuleDataSource,
		NewNetworkBridgeVlanDataSource,
		NewNetworkWireguardPeerDataSource,
		// 1.1 firewall
		NewFirewallForwardingDataSource,
		NewFirewallDefaultsDataSource,
		// 1.1 dhcp
		NewDhcpServerDataSource,
		NewDhcpDnsmasqDataSource,
		NewDhcpOdhcpdDataSource,
		// 1.1 snmpd
		NewSnmpdAccessDataSource,
		NewSnmpdAgentDataSource,
		NewSnmpdCom2secDataSource,
		NewSnmpdGroupDataSource,
		NewSnmpdSystemDataSource,
		// 1.1 uhttpd / dropbear
		NewUhttpdCertDataSource,
		NewUhttpdInstanceDataSource,
		NewDropbearInstanceDataSource,
		// 1.1 system / sqm
		NewSystemTimeserverDataSource,
		NewSqmQueueDataSource,
		// 1.1 vnstat + misc singletons
		NewVnstatInterfaceDataSource,
		NewVnstatConfigDataSource,
		NewLldpdConfigDataSource,
		NewPrometheusNodeExporterLuaConfigDataSource,
		NewUnboundServerDataSource,
		// 1.1 packages
		NewPackageDataSource,
		NewPackageFeedDataSource,
	}
}

func firstNonEmpty(configured types.String, fallbacks ...string) string {
	if !configured.IsNull() && !configured.IsUnknown() && configured.ValueString() != "" {
		return configured.ValueString()
	}
	for _, f := range fallbacks {
		if f != "" {
			return f
		}
	}
	return ""
}
