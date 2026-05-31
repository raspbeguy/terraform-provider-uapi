package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &firewallDefaultsDataSource{}
	_ datasource.DataSourceWithConfigure = &firewallDefaultsDataSource{}
)

type firewallDefaultsDataSource struct{ client *client.Client }

func NewFirewallDefaultsDataSource() datasource.DataSource { return &firewallDefaultsDataSource{} }

func (d *firewallDefaultsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_defaults"
}

func (d *firewallDefaultsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *firewallDefaultsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "The global firewall defaults (uci firewall.defaults).",
		Attributes: map[string]dsschema.Attribute{
			"id":                 dsComputedString("Stable id of the defaults section."),
			"managed":            dsManagedAttribute(),
			"input":              dsComputedString("Default policy for input traffic: ACCEPT, REJECT, or DROP."),
			"output":             dsComputedString("Default policy for output traffic: ACCEPT, REJECT, or DROP."),
			"forward":            dsComputedString("Default policy for forwarded traffic: ACCEPT, REJECT, or DROP."),
			"syn_flood":          dsComputedBool("Whether SYN-flood protection is enabled."),
			"drop_invalid":       dsComputedBool("Whether packets in an invalid conntrack state are dropped."),
			"synflood_burst":     dsComputedString("SYN-flood burst limit."),
			"synflood_rate":      dsComputedString("SYN-flood rate limit."),
			"tcp_syncookies":     dsComputedBool("Whether TCP SYN cookies are enabled."),
			"flow_offloading":    dsComputedBool("Whether software flow offloading is enabled."),
			"flow_offloading_hw": dsComputedBool("Whether hardware flow offloading is enabled."),
		},
	}
}

func (d *firewallDefaultsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	obj, found, err := d.client.GetObject(ctx, firewallDefaultsPath)
	if err != nil {
		resp.Diagnostics.AddError("Error reading firewall defaults", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Firewall defaults not found", "The firewall defaults singleton is missing on the router")
		return
	}
	var m firewallDefaultsModel
	(&firewallDefaultsResource{}).read(ctx, obj, &m)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
