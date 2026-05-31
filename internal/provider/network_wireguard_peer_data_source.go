package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &networkWireguardPeerDataSource{}
	_ datasource.DataSourceWithConfigure = &networkWireguardPeerDataSource{}
)

type networkWireguardPeerDataSource struct{ client *client.Client }

func NewNetworkWireguardPeerDataSource() datasource.DataSource {
	return &networkWireguardPeerDataSource{}
}

func (d *networkWireguardPeerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_wireguard_peer"
}

func (d *networkWireguardPeerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *networkWireguardPeerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up a WireGuard peer by id. The preshared key is never returned.",
		Attributes: map[string]dsschema.Attribute{
			"id":                   dsIDAttribute(),
			"managed":              dsManagedAttribute(),
			"interface":            dsComputedString("Parent WireGuard interface name."),
			"description":          dsComputedString("Human-readable peer description."),
			"public_key":           dsComputedString("Peer public key (44-char base64)."),
			"preshared_key":        dsschema.StringAttribute{Computed: true, Sensitive: true, Description: "Always null; the API never returns the preshared key."},
			"has_preshared_key":    dsComputedBool("Whether a preshared key is configured on the router."),
			"allowed_ips":          dsComputedStringList("IPv4 CIDRs routed to this peer."),
			"endpoint_host":        dsComputedString("Peer endpoint hostname or IP address."),
			"endpoint_port":        dsComputedString("Peer endpoint UDP port (1-65535)."),
			"persistent_keepalive": dsComputedString("Keepalive interval in seconds (0-65535)."),
			"route_allowed_ips":    dsComputedBool("Whether routes for allowed_ips are created automatically."),
			"disabled":             dsComputedBool("Whether this peer is disabled."),
		},
	}
}

func (d *networkWireguardPeerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var m networkWireguardPeerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := d.client.GetObject(ctx, "/"+networkWireguardPeerCollection+"/"+m.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading network WireGuard peer", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Network WireGuard peer not found", "No network WireGuard peer with id "+m.ID.ValueString())
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	(&networkWireguardPeerResource{}).read(ctx, obj, &m, ds)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
