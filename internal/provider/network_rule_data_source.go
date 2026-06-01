package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &networkRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &networkRuleDataSource{}
)

type networkRuleDataSource struct{ client *client.Client }

func NewNetworkRuleDataSource() datasource.DataSource { return &networkRuleDataSource{} }

func (d *networkRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_rule"
}

func (d *networkRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *networkRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up an IP routing policy rule by id.",
		Attributes: map[string]dsschema.Attribute{
			"id":       dsIDAttribute(),
			"managed":  dsManagedAttribute(),
			"etag":     dsComputedString("Opaque ETag of the resource's current state."),
			"in":       dsComputedString("Incoming interface selector."),
			"out":      dsComputedString("Outgoing interface selector."),
			"src":      dsComputedString("Source IPv4 address or CIDR selector."),
			"dest":     dsComputedString("Destination IPv4 address or CIDR selector."),
			"priority": dsComputedString("Rule priority (0-32766)."),
			"lookup":   dsComputedString("Routing table to look up."),
			"goto":     dsComputedString("Priority to jump to."),
			"action":   dsComputedString("Rule action: lookup, goto, unreachable, prohibit, blackhole, or throw."),
			"invert":   dsComputedBool("Whether the rule selectors are inverted."),
			"mark":     dsComputedString("Firewall mark to match."),
		},
	}
}

func (d *networkRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var m networkRuleModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := d.client.GetObject(ctx, "/"+networkRuleCollection+"/"+m.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading network rule", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Network rule not found", "No network rule with id "+m.ID.ValueString())
		return
	}
	(&networkRuleResource{}).read(ctx, obj, &m)
	m.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
