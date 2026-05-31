package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &snmpdGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &snmpdGroupDataSource{}
)

type snmpdGroupDataSource struct{ client *client.Client }

func NewSnmpdGroupDataSource() datasource.DataSource { return &snmpdGroupDataSource{} }

func (d *snmpdGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snmpd_group"
}

func (d *snmpdGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *snmpdGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up an SNMP VACM group definition by id.",
		Attributes: map[string]dsschema.Attribute{
			"id":      dsIDAttribute(),
			"managed": dsManagedAttribute(),
			"group":   dsComputedString("Group name referenced by snmpd access entries."),
			"version": dsComputedString("Security model for the membership: v1, v2c, or usm."),
			"secname": dsComputedString("Security name added to the group."),
		},
	}
}

func (d *snmpdGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var m snmpdGroupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := d.client.GetObject(ctx, "/"+snmpdGroupCollection+"/"+m.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading snmpd group", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Snmpd group not found", "No snmpd group with id "+m.ID.ValueString())
		return
	}
	(&snmpdGroupResource{}).read(ctx, obj, &m)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
