package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &vnstatInterfaceDataSource{}
	_ datasource.DataSourceWithConfigure = &vnstatInterfaceDataSource{}
)

type vnstatInterfaceDataSource struct{ client *client.Client }

func NewVnstatInterfaceDataSource() datasource.DataSource { return &vnstatInterfaceDataSource{} }

func (d *vnstatInterfaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vnstat_interface"
}

func (d *vnstatInterfaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *vnstatInterfaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up a vnstat monitored interface by id.",
		Attributes: map[string]dsschema.Attribute{
			"id":        dsIDAttribute(),
			"managed":   dsManagedAttribute(),
			"interface": dsComputedString("Name of the monitored network interface."),
			"enabled":   dsComputedBool("Whether monitoring of this interface is enabled."),
		},
	}
}

func (d *vnstatInterfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var m vnstatInterfaceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := d.client.GetObject(ctx, "/"+vnstatInterfaceCollection+"/"+m.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading vnstat interface", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Vnstat interface not found", "No vnstat interface with id "+m.ID.ValueString())
		return
	}
	(&vnstatInterfaceResource{}).read(ctx, obj, &m)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
