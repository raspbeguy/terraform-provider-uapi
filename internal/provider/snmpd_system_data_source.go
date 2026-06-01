package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &snmpdSystemDataSource{}
	_ datasource.DataSourceWithConfigure = &snmpdSystemDataSource{}
)

type snmpdSystemDataSource struct{ client *client.Client }

func NewSnmpdSystemDataSource() datasource.DataSource { return &snmpdSystemDataSource{} }

func (d *snmpdSystemDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snmpd_system"
}

func (d *snmpdSystemDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *snmpdSystemDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "The SNMP agent system identity (uci snmpd.system).",
		Attributes: map[string]dsschema.Attribute{
			"id":            dsComputedString("Stable id of the snmpd system section."),
			"managed":       dsManagedAttribute(),
			"etag":          dsComputedString("Opaque ETag of the resource's current state."),
			"sys_location":  dsComputedString("Value reported as sysLocation (the device's physical location)."),
			"sys_contact":   dsComputedString("Value reported as sysContact (the administrative contact)."),
			"sys_name":      dsComputedString("Value reported as sysName (the administratively assigned name)."),
			"sys_services":  dsComputedString("Value reported as sysServices (the OSI layer services bitmask)."),
			"sys_descr":     dsComputedString("Value reported as sysDescr (a textual description of the entity)."),
			"sys_object_id": dsComputedString("Value reported as sysObjectID (the vendor authoritative object identifier)."),
		},
	}
}

func (d *snmpdSystemDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	obj, etag, found, err := d.client.GetObject(ctx, snmpdSystemPath)
	if err != nil {
		resp.Diagnostics.AddError("Error reading snmpd system identity", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Snmpd system identity not found", "The snmpd system singleton is missing on the router")
		return
	}
	var m snmpdSystemModel
	(&snmpdSystemResource{}).read(ctx, obj, &m)
	m.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
