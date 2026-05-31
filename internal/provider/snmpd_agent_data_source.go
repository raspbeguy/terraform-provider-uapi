package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &snmpdAgentDataSource{}
	_ datasource.DataSourceWithConfigure = &snmpdAgentDataSource{}
)

type snmpdAgentDataSource struct{ client *client.Client }

func NewSnmpdAgentDataSource() datasource.DataSource { return &snmpdAgentDataSource{} }

func (d *snmpdAgentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snmpd_agent"
}

func (d *snmpdAgentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *snmpdAgentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up an SNMP agent listener configuration by id.",
		Attributes: map[string]dsschema.Attribute{
			"id":           dsIDAttribute(),
			"managed":      dsManagedAttribute(),
			"agentaddress": dsComputedStringList("Addresses the SNMP agent listens on (e.g. UDP:161, udp6:161)."),
		},
	}
}

func (d *snmpdAgentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var m snmpdAgentModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := d.client.GetObject(ctx, "/"+snmpdAgentCollection+"/"+m.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading snmpd agent", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Snmpd agent not found", "No snmpd agent with id "+m.ID.ValueString())
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	(&snmpdAgentResource{}).read(ctx, obj, &m, ds)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
