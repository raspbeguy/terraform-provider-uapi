package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &systemTimeserverDataSource{}
	_ datasource.DataSourceWithConfigure = &systemTimeserverDataSource{}
)

type systemTimeserverDataSource struct{ client *client.Client }

func NewSystemTimeserverDataSource() datasource.DataSource { return &systemTimeserverDataSource{} }

func (d *systemTimeserverDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_timeserver"
}

func (d *systemTimeserverDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *systemTimeserverDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up an NTP timeserver configuration by id.",
		Attributes: map[string]dsschema.Attribute{
			"id":            dsIDAttribute(),
			"managed":       dsManagedAttribute(),
			"enabled":       dsComputedBool("Whether the NTP client is enabled."),
			"enable_server": dsComputedBool("Whether the router acts as an NTP server for the local network."),
			"interface":     dsComputedString("Network interface the NTP server binds to."),
			"server":        dsComputedStringList("Upstream NTP server hostnames."),
			"use_dhcp":      dsComputedBool("Whether NTP servers advertised over DHCP are used."),
		},
	}
}

func (d *systemTimeserverDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var m systemTimeserverModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := d.client.GetObject(ctx, "/"+systemTimeserverCollection+"/"+m.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading system timeserver", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("System timeserver not found", "No system timeserver with id "+m.ID.ValueString())
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	(&systemTimeserverResource{}).read(ctx, obj, &m, ds)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
