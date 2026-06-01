package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &networkRouteDataSource{}
	_ datasource.DataSourceWithConfigure = &networkRouteDataSource{}
)

type networkRouteDataSource struct{ client *client.Client }

func NewNetworkRouteDataSource() datasource.DataSource { return &networkRouteDataSource{} }

func (d *networkRouteDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_route"
}

func (d *networkRouteDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *networkRouteDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up a static network route by id.",
		Attributes: map[string]dsschema.Attribute{
			"id":        dsIDAttribute(),
			"managed":   dsManagedAttribute(),
			"etag":      dsComputedString("Opaque ETag of the resource's current state."),
			"interface": dsComputedString("Logical network interface the route is bound to."),
			"target":    dsComputedString("Destination IPv4 address or CIDR."),
			"netmask":   dsComputedString("Destination netmask, when target is a bare address."),
			"gateway":   dsComputedString("Next-hop gateway IPv4 address."),
			"table":     dsComputedString("Routing table the route is installed into."),
			"metric":    dsComputedString("Route metric (priority)."),
			"mtu":       dsComputedString("Path MTU for the route."),
			"source":    dsComputedString("Preferred source IPv4 address or CIDR."),
			"type":      dsComputedString("Route type: unicast, blackhole, unreachable, prohibit, throw, anycast, multicast, local, or broadcast."),
		},
	}
}

func (d *networkRouteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var m networkRouteModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := d.client.GetObject(ctx, "/"+networkRouteCollection+"/"+m.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading network route", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Network route not found", "No network route with id "+m.ID.ValueString())
		return
	}
	(&networkRouteResource{}).read(ctx, obj, &m)
	m.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
