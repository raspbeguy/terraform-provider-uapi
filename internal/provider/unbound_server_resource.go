package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const unboundServerPath = "/unbound/server"

var (
	_ resource.Resource                = &unboundServerResource{}
	_ resource.ResourceWithConfigure   = &unboundServerResource{}
	_ resource.ResourceWithImportState = &unboundServerResource{}
)

type unboundServerResource struct {
	client *client.Client
}

func NewUnboundServerResource() resource.Resource {
	return &unboundServerResource{}
}

type unboundServerModel struct {
	ID               types.String `tfsdk:"id"`
	Managed          types.Bool   `tfsdk:"managed"`
	ETag             types.String `tfsdk:"etag"`
	Enabled          types.Bool   `tfsdk:"enabled"`
	ListenPort       types.String `tfsdk:"listen_port"`
	DHCPLink         types.String `tfsdk:"dhcp_link"`
	AddLocalFQDN     types.String `tfsdk:"add_local_fqdn"`
	AddWANFQDN       types.String `tfsdk:"add_wan_fqdn"`
	DNSSECEnabled    types.Bool   `tfsdk:"dnssec_enabled"`
	Recursion        types.String `tfsdk:"recursion"`
	Resource         types.String `tfsdk:"resource"`
	Protocol         types.String `tfsdk:"protocol"`
	QueryMinimize    types.Bool   `tfsdk:"query_minimize"`
	Prefetch         types.Bool   `tfsdk:"prefetch"`
	ManualConf       types.Bool   `tfsdk:"manual_conf"`
	ExtendedStats    types.Bool   `tfsdk:"extended_stats"`
	InterfaceAuto    types.Bool   `tfsdk:"interface_auto"`
	LocalService     types.Bool   `tfsdk:"localservice"`
	HideBinddata     types.Bool   `tfsdk:"hide_binddata"`
	RebindProtection types.String `tfsdk:"rebind_protection"`
	NumThreads       types.String `tfsdk:"num_threads"`
	TTLMin           types.String `tfsdk:"ttl_min"`
	Domain           types.String `tfsdk:"domain"`
	DomainType       types.String `tfsdk:"domain_type"`
}

func (r *unboundServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_unbound_server"
}

func (r *unboundServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *unboundServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Global unbound resolver settings (uci unbound.unbound). This is a singleton: it cannot be " +
			"created or destroyed. `terraform destroy` only removes it from state; the underlying " +
			"settings are left as-is on the router.",
		Attributes: map[string]schema.Attribute{
			"id":                computedIDAttribute(),
			"managed":           managedAttribute(),
			"etag":              etagAttribute(),
			"enabled":           optionalComputedBool("Whether the unbound resolver is enabled. Defaults to true."),
			"listen_port":       schema.StringAttribute{Optional: true, Description: "Port unbound listens on for DNS queries."},
			"dhcp_link":         schema.StringAttribute{Optional: true, Description: "DHCP integration source: none, odhcpd, or dnsmasq."},
			"add_local_fqdn":    schema.StringAttribute{Optional: true, Description: "How aggressively to add LAN host FQDN records."},
			"add_wan_fqdn":      schema.StringAttribute{Optional: true, Description: "How aggressively to add WAN host FQDN records."},
			"dnssec_enabled":    optionalComputedBool("Enable DNSSEC validation. Defaults to false."),
			"recursion":         schema.StringAttribute{Optional: true, Description: "Recursion tuning preset: default, passive, or aggressive."},
			"resource":          schema.StringAttribute{Optional: true, Description: "Memory/cache sizing preset: tiny, small, medium, large, big, or huge."},
			"protocol":          schema.StringAttribute{Optional: true, Description: "IP protocol mode: auto, ip4_only, ip6_only, or mixed."},
			"query_minimize":    optionalComputedBool("Enable QNAME minimization. Defaults to false."),
			"prefetch":          optionalComputedBool("Prefetch popular cache entries before they expire. Defaults to false."),
			"manual_conf":       optionalComputedBool("Skip uci generation and use a hand-written /etc/unbound/unbound.conf. Defaults to false."),
			"extended_stats":    optionalComputedBool("Emit extended statistics (stats-extended). Defaults to false."),
			"interface_auto":    optionalComputedBool("Bind to all interfaces (interface-automatic). Disable to bind manually. Defaults to true."),
			"localservice":      optionalComputedBool("Restrict access to clients on local subnets. Defaults to true."),
			"hide_binddata":     optionalComputedBool("Hide identity and version from binddata queries. Defaults to true."),
			"rebind_protection": schema.StringAttribute{Optional: true, Description: "DNS rebind protection: 0 (off), 1 (private nets), or 2 (all rebind attacks blocked)."},
			"num_threads":       schema.StringAttribute{Optional: true, Description: "Number of resolver threads (1-64)."},
			"ttl_min":           schema.StringAttribute{Optional: true, Description: "Minimum TTL in seconds to enforce on cached records (0-86400)."},
			"domain":            schema.StringAttribute{Optional: true, Description: "Local domain name unbound serves authoritatively."},
			"domain_type":       schema.StringAttribute{Optional: true, Description: "Local-zone type for the configured domain (e.g. deny, refuse, static, transparent, redirect)."},
		},
	}
}

func (r *unboundServerResource) body(_ context.Context, m unboundServerModel) map[string]any {
	out := map[string]any{}
	putBool(out, "enabled", m.Enabled)
	putStr(out, "listen_port", m.ListenPort)
	putStr(out, "dhcp_link", m.DHCPLink)
	putStr(out, "add_local_fqdn", m.AddLocalFQDN)
	putStr(out, "add_wan_fqdn", m.AddWANFQDN)
	putBool(out, "dnssec_enabled", m.DNSSECEnabled)
	putStr(out, "recursion", m.Recursion)
	putStr(out, "resource", m.Resource)
	putStr(out, "protocol", m.Protocol)
	putBool(out, "query_minimize", m.QueryMinimize)
	putBool(out, "prefetch", m.Prefetch)
	putBool(out, "manual_conf", m.ManualConf)
	putBool(out, "extended_stats", m.ExtendedStats)
	putBool(out, "interface_auto", m.InterfaceAuto)
	putBool(out, "localservice", m.LocalService)
	putBool(out, "hide_binddata", m.HideBinddata)
	putStr(out, "rebind_protection", m.RebindProtection)
	putStr(out, "num_threads", m.NumThreads)
	putStr(out, "ttl_min", m.TTLMin)
	putStr(out, "domain", m.Domain)
	putStr(out, "domain_type", m.DomainType)
	return out
}

func (r *unboundServerResource) read(_ context.Context, obj map[string]any, m *unboundServerModel) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.Enabled = boolVal(obj, "enabled")
	m.ListenPort = strVal(obj, "listen_port")
	m.DHCPLink = strVal(obj, "dhcp_link")
	m.AddLocalFQDN = strVal(obj, "add_local_fqdn")
	m.AddWANFQDN = strVal(obj, "add_wan_fqdn")
	m.DNSSECEnabled = boolVal(obj, "dnssec_enabled")
	m.Recursion = strVal(obj, "recursion")
	m.Resource = strVal(obj, "resource")
	m.Protocol = strVal(obj, "protocol")
	m.QueryMinimize = boolVal(obj, "query_minimize")
	m.Prefetch = boolVal(obj, "prefetch")
	m.ManualConf = boolVal(obj, "manual_conf")
	m.ExtendedStats = boolVal(obj, "extended_stats")
	m.InterfaceAuto = boolVal(obj, "interface_auto")
	m.LocalService = boolVal(obj, "localservice")
	m.HideBinddata = boolVal(obj, "hide_binddata")
	m.RebindProtection = strVal(obj, "rebind_protection")
	m.NumThreads = strVal(obj, "num_threads")
	m.TTLMin = strVal(obj, "ttl_min")
	m.Domain = strVal(obj, "domain")
	m.DomainType = strVal(obj, "domain_type")
}

func (r *unboundServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan unboundServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, err := r.client.Patch(ctx, unboundServerPath, r.body(ctx, plan), "")
	if err != nil {
		writeErr(&resp.Diagnostics, "configuring", "unbound settings", err)
		return
	}
	r.read(ctx, obj, &plan)
	plan.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *unboundServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state unboundServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := r.client.GetObject(ctx, unboundServerPath)
	if err != nil {
		resp.Diagnostics.AddError("Error reading unbound settings", err.Error())
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	r.read(ctx, obj, &state)
	state.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *unboundServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state unboundServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, err := r.client.Patch(ctx, unboundServerPath, r.body(ctx, plan), state.ETag.ValueString())
	if err != nil {
		writeErr(&resp.Diagnostics, "updating", "unbound settings", err)
		return
	}
	r.read(ctx, obj, &plan)
	plan.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete is a no-op: the unbound server singleton cannot be removed. State is
// dropped by the framework once this returns.
func (r *unboundServerResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func (r *unboundServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
