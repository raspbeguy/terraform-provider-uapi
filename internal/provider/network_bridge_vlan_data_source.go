package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &networkBridgeVlanDataSource{}
	_ datasource.DataSourceWithConfigure = &networkBridgeVlanDataSource{}
)

type networkBridgeVlanDataSource struct{ client *client.Client }

func NewNetworkBridgeVlanDataSource() datasource.DataSource { return &networkBridgeVlanDataSource{} }

func (d *networkBridgeVlanDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_bridge_vlan"
}

func (d *networkBridgeVlanDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *networkBridgeVlanDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up a bridge VLAN by id.",
		Attributes: map[string]dsschema.Attribute{
			"id":      dsIDAttribute(),
			"managed": dsManagedAttribute(),
			"etag":    dsComputedString("Opaque ETag of the resource's current state."),
			"device":  dsComputedString("Bridge device name this VLAN belongs to."),
			"vlan":    dsComputedString("VLAN id (1-4094)."),
			"ports":   dsComputedStringList("Member ports, each as <name>[:t|:u|:*]."),
		},
	}
}

func (d *networkBridgeVlanDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var m networkBridgeVlanModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := d.client.GetObject(ctx, "/"+networkBridgeVlanCollection+"/"+m.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading network bridge VLAN", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Network bridge VLAN not found", "No network bridge VLAN with id "+m.ID.ValueString())
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	(&networkBridgeVlanResource{}).read(ctx, obj, &m, ds)
	m.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
