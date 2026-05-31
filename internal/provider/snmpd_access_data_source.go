package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &snmpdAccessDataSource{}
	_ datasource.DataSourceWithConfigure = &snmpdAccessDataSource{}
)

type snmpdAccessDataSource struct{ client *client.Client }

func NewSnmpdAccessDataSource() datasource.DataSource { return &snmpdAccessDataSource{} }

func (d *snmpdAccessDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snmpd_access"
}

func (d *snmpdAccessDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *snmpdAccessDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up an SNMP VACM access entry by id.",
		Attributes: map[string]dsschema.Attribute{
			"id":      dsIDAttribute(),
			"managed": dsManagedAttribute(),
			"group":   dsComputedString("Name of the snmpd group this access entry applies to."),
			"context": dsComputedString("SNMP context the access entry matches."),
			"version": dsComputedString("Security model the entry matches: any, v1, v2c, or usm."),
			"level":   dsComputedString("Required security level: noauth, auth, or priv."),
			"prefix":  dsComputedString("Context match mode: exact or prefix."),
			"read":    dsComputedString("View name granted read access."),
			"write":   dsComputedString("View name granted write access."),
			"notify":  dsComputedString("View name granted notify access."),
		},
	}
}

func (d *snmpdAccessDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var m snmpdAccessModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := d.client.GetObject(ctx, "/"+snmpdAccessCollection+"/"+m.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading snmpd access", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Snmpd access not found", "No snmpd access with id "+m.ID.ValueString())
		return
	}
	(&snmpdAccessResource{}).read(ctx, obj, &m)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
