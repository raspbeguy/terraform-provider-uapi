package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &firewallRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &firewallRuleDataSource{}
)

type firewallRuleDataSource struct{ client *client.Client }

func NewFirewallRuleDataSource() datasource.DataSource { return &firewallRuleDataSource{} }

func (d *firewallRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_rule"
}

func (d *firewallRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *firewallRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up a firewall rule by id.",
		Attributes: map[string]dsschema.Attribute{
			"id":      dsIDAttribute(),
			"managed": dsManagedAttribute(),
			"etag":    dsComputedString("Opaque ETag of the resource's current state."),
			"name":    dsComputedString("Rule name."),
			"target":  dsComputedString("Rule action: ACCEPT, REJECT, DROP, NOTRACK, or MARK."),
			"enabled": dsComputedBool("Whether the rule is active."),
			"match": dsschema.SingleNestedAttribute{
				Computed:    true,
				Description: "Match conditions for the rule.",
				Attributes: map[string]dsschema.Attribute{
					"src_zone":  dsComputedString("Source firewall zone name."),
					"dest_zone": dsComputedString("Destination firewall zone name."),
					"src_ip":    dsComputedStringList("Source IP addresses or CIDRs."),
					"dest_ip":   dsComputedStringList("Destination IP addresses or CIDRs."),
					"src_port":  dsComputedStringList("Source ports."),
					"dest_port": dsComputedStringList("Destination ports."),
					"proto":     dsComputedStringList("Protocols: tcp, udp, icmp, icmpv6, esp, ah, any, all."),
					"family":    dsComputedString("Address family: any, ipv4, or ipv6."),
				},
			},
		},
	}
}

func (d *firewallRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var m firewallRuleModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := d.client.GetObject(ctx, "/"+firewallRuleCollection+"/"+m.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading firewall rule", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Firewall rule not found", "No firewall rule with id "+m.ID.ValueString())
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	(&firewallRuleResource{}).read(ctx, obj, &m, ds)
	m.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}

var (
	_ datasource.DataSource              = &firewallZoneDataSource{}
	_ datasource.DataSourceWithConfigure = &firewallZoneDataSource{}
)

type firewallZoneDataSource struct{ client *client.Client }

func NewFirewallZoneDataSource() datasource.DataSource { return &firewallZoneDataSource{} }

func (d *firewallZoneDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_zone"
}

func (d *firewallZoneDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *firewallZoneDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up a firewall zone by id.",
		Attributes: map[string]dsschema.Attribute{
			"id":      dsIDAttribute(),
			"managed": dsManagedAttribute(),
			"etag":    dsComputedString("Opaque ETag of the resource's current state."),
			"name":    dsComputedString("Zone name."),
			"input":   dsComputedString("Default policy for input traffic: ACCEPT, REJECT, or DROP."),
			"output":  dsComputedString("Default policy for output traffic."),
			"forward": dsComputedString("Default policy for forwarded traffic."),
			"network": dsComputedStringList("Network interfaces covered by this zone."),
			"masq":    dsComputedBool("Whether masquerading (NAT) is enabled."),
			"mtu_fix": dsComputedBool("Whether MSS clamping is enabled."),
			"family":  dsComputedString("Address family: any, ipv4, or ipv6."),
		},
	}
}

func (d *firewallZoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var m firewallZoneModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := d.client.GetObject(ctx, "/"+firewallZoneCollection+"/"+m.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading firewall zone", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Firewall zone not found", "No firewall zone with id "+m.ID.ValueString())
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	(&firewallZoneResource{}).read(ctx, obj, &m, ds)
	m.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}

var (
	_ datasource.DataSource              = &firewallRedirectDataSource{}
	_ datasource.DataSourceWithConfigure = &firewallRedirectDataSource{}
)

type firewallRedirectDataSource struct{ client *client.Client }

func NewFirewallRedirectDataSource() datasource.DataSource { return &firewallRedirectDataSource{} }

func (d *firewallRedirectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_redirect"
}

func (d *firewallRedirectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *firewallRedirectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up a firewall redirect by id.",
		Attributes: map[string]dsschema.Attribute{
			"id":      dsIDAttribute(),
			"managed": dsManagedAttribute(),
			"etag":    dsComputedString("Opaque ETag of the resource's current state."),
			"name":    dsComputedString("Redirect name."),
			"target":  dsComputedString("Redirect type: DNAT or SNAT."),
			"enabled": dsComputedBool("Whether the redirect is active."),
			"match": dsschema.SingleNestedAttribute{
				Computed:    true,
				Description: "Match conditions for the redirect.",
				Attributes: map[string]dsschema.Attribute{
					"src_zone":  dsComputedString("Source firewall zone name."),
					"dest_zone": dsComputedString("Destination firewall zone name."),
					"src_ip":    dsComputedStringList("Source IP addresses or CIDRs."),
					"src_port":  dsComputedStringList("Source ports."),
					"src_dport": dsComputedStringList("Incoming (destination) ports to redirect."),
					"dest_ip":   dsComputedStringList("Internal destination IP addresses."),
					"dest_port": dsComputedStringList("Internal destination ports."),
					"proto":     dsComputedStringList("Protocols: tcp, udp, icmp, icmpv6, esp, ah, any, all."),
					"family":    dsComputedString("Address family: any, ipv4, or ipv6."),
				},
			},
			"reflection":      dsComputedBool("Whether NAT loopback / hairpinning is enabled for this redirect."),
			"reflection_src":  dsComputedString("Source address used for hairpinned packets: internal or external."),
			"reflection_zone": dsComputedStringList("Firewall zones in which NAT reflection is applied."),
		},
	}
}

func (d *firewallRedirectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var m firewallRedirectModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := d.client.GetObject(ctx, "/"+firewallRedirectCollection+"/"+m.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading firewall redirect", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Firewall redirect not found", "No firewall redirect with id "+m.ID.ValueString())
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	(&firewallRedirectResource{}).read(ctx, obj, &m, ds)
	m.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}

func dsIDAttribute() dsschema.StringAttribute {
	return dsschema.StringAttribute{Required: true, Description: "Resource id to look up."}
}

func dsComputedString(desc string) dsschema.StringAttribute {
	return dsschema.StringAttribute{Computed: true, Description: desc}
}

func dsComputedBool(desc string) dsschema.BoolAttribute {
	return dsschema.BoolAttribute{Computed: true, Description: desc}
}

func dsComputedStringList(desc string) dsschema.ListAttribute {
	return dsschema.ListAttribute{ElementType: types.StringType, Computed: true, Description: desc}
}

func dsManagedAttribute() dsschema.BoolAttribute {
	return dsComputedBool("Whether the underlying uci section is uapi-managed.")
}
