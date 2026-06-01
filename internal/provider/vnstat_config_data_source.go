package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &vnstatConfigDataSource{}
	_ datasource.DataSourceWithConfigure = &vnstatConfigDataSource{}
)

type vnstatConfigDataSource struct{ client *client.Client }

func NewVnstatConfigDataSource() datasource.DataSource { return &vnstatConfigDataSource{} }

func (d *vnstatConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vnstat_config"
}

func (d *vnstatConfigDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *vnstatConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "The global vnstat daemon settings (uci vnstat.vnstat).",
		Attributes: map[string]dsschema.Attribute{
			"id":                   dsComputedString("Stable id of the vnstat config section."),
			"managed":              dsManagedAttribute(),
			"etag":                 dsComputedString("Opaque ETag of the resource's current state."),
			"database_dir":         dsComputedString("Directory where vnstat stores its databases."),
			"interface_5min_hours": dsComputedString("Hours of 5-minute resolution data to keep per interface."),
			"month_rotate":         dsComputedString("Day of the month on which monthly statistics are reset."),
		},
	}
}

func (d *vnstatConfigDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	obj, etag, found, err := d.client.GetObject(ctx, vnstatConfigPath)
	if err != nil {
		resp.Diagnostics.AddError("Error reading vnstat settings", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Vnstat settings not found", "The vnstat config singleton is missing on the router")
		return
	}
	var m vnstatConfigModel
	(&vnstatConfigResource{}).read(ctx, obj, &m)
	m.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
