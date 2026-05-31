package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const dhcpDnsmasqPath = "/dhcp/dnsmasq"

var (
	_ resource.Resource                = &dhcpDnsmasqResource{}
	_ resource.ResourceWithConfigure   = &dhcpDnsmasqResource{}
	_ resource.ResourceWithImportState = &dhcpDnsmasqResource{}
)

type dhcpDnsmasqResource struct {
	client *client.Client
}

func NewDhcpDnsmasqResource() resource.Resource {
	return &dhcpDnsmasqResource{}
}

type dhcpDnsmasqModel struct {
	ID               types.String `tfsdk:"id"`
	Managed          types.Bool   `tfsdk:"managed"`
	Domain           types.String `tfsdk:"domain"`
	Local            types.String `tfsdk:"local"`
	Noresolv         types.Bool   `tfsdk:"noresolv"`
	RebindProtection types.Bool   `tfsdk:"rebind_protection"`
	Expandhosts      types.Bool   `tfsdk:"expandhosts"`
	Cachesize        types.String `tfsdk:"cachesize"`
	Port             types.String `tfsdk:"port"`
	Domainneeded     types.Bool   `tfsdk:"domainneeded"`
	Boguspriv        types.Bool   `tfsdk:"boguspriv"`
	Filterwin2k      types.Bool   `tfsdk:"filterwin2k"`
	Authoritative    types.Bool   `tfsdk:"authoritative"`
	Readethers       types.Bool   `tfsdk:"readethers"`
	Leasefile        types.String `tfsdk:"leasefile"`
	Resolvfile       types.String `tfsdk:"resolvfile"`
	Server           types.List   `tfsdk:"server"`
	Address          types.List   `tfsdk:"address"`
	Nonwildcard      types.Bool   `tfsdk:"nonwildcard"`
}

func (r *dhcpDnsmasqResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcp_dnsmasq"
}

func (r *dhcpDnsmasqResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *dhcpDnsmasqResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Global dnsmasq settings (uci dhcp.dnsmasq). This is a singleton: it cannot be " +
			"created or destroyed. `terraform destroy` only removes it from state; the underlying " +
			"settings are left as-is on the router.",
		Attributes: map[string]schema.Attribute{
			"id":      computedIDAttribute(),
			"managed": managedAttribute(),
			"domain": schema.StringAttribute{
				Optional:    true,
				Description: "Local DNS domain.",
			},
			"local": schema.StringAttribute{
				Optional:    true,
				Description: "Local domain suffix served authoritatively (e.g. /lan/).",
			},
			"noresolv":          optionalComputedBool("Do not read upstream resolvers from resolvfile. Defaults to false."),
			"rebind_protection": optionalComputedBool("Enable DNS rebind protection. Defaults to true."),
			"expandhosts":       optionalComputedBool("Add the local domain to names in /etc/hosts. Defaults to false."),
			"cachesize": schema.StringAttribute{
				Optional:    true,
				Description: "DNS cache size (0-1000000).",
			},
			"port": schema.StringAttribute{
				Optional:    true,
				Description: "DNS service port (1-65535).",
			},
			"domainneeded":  optionalComputedBool("Never forward plain names without a domain. Defaults to true."),
			"boguspriv":     optionalComputedBool("Never forward reverse lookups for private ranges. Defaults to true."),
			"filterwin2k":   optionalComputedBool("Filter Windows DNS queries that pollute upstream. Defaults to false."),
			"authoritative": optionalComputedBool("Act as the authoritative DHCP server on the segment. Defaults to true."),
			"readethers":    optionalComputedBool("Read static leases from /etc/ethers. Defaults to true."),
			"leasefile": schema.StringAttribute{
				Optional:    true,
				Description: "Path to the DHCP lease file.",
			},
			"resolvfile": schema.StringAttribute{
				Optional:    true,
				Description: "Path to the upstream resolver file.",
			},
			"server":      optionalComputedStringList("Upstream DNS servers."),
			"address":     optionalComputedStringList("Static DNS address overrides (e.g. /router.lan/192.168.1.1)."),
			"nonwildcard": optionalComputedBool("Bind only to configured interfaces instead of the wildcard address. Defaults to true."),
		},
	}
}

func (r *dhcpDnsmasqResource) body(ctx context.Context, m dhcpDnsmasqModel, diags *diagsink) map[string]any {
	out := map[string]any{}
	putStr(out, "domain", m.Domain)
	putStr(out, "local", m.Local)
	putBool(out, "noresolv", m.Noresolv)
	putBool(out, "rebind_protection", m.RebindProtection)
	putBool(out, "expandhosts", m.Expandhosts)
	putStr(out, "cachesize", m.Cachesize)
	putStr(out, "port", m.Port)
	putBool(out, "domainneeded", m.Domainneeded)
	putBool(out, "boguspriv", m.Boguspriv)
	putBool(out, "filterwin2k", m.Filterwin2k)
	putBool(out, "authoritative", m.Authoritative)
	putBool(out, "readethers", m.Readethers)
	putStr(out, "leasefile", m.Leasefile)
	putStr(out, "resolvfile", m.Resolvfile)
	putList(ctx, out, "server", m.Server, diags.d)
	putList(ctx, out, "address", m.Address, diags.d)
	putBool(out, "nonwildcard", m.Nonwildcard)
	return out
}

func (r *dhcpDnsmasqResource) read(ctx context.Context, obj map[string]any, m *dhcpDnsmasqModel, diags *diagsink) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.Domain = strVal(obj, "domain")
	m.Local = strVal(obj, "local")
	m.Noresolv = boolVal(obj, "noresolv")
	m.RebindProtection = boolVal(obj, "rebind_protection")
	m.Expandhosts = boolVal(obj, "expandhosts")
	m.Cachesize = strVal(obj, "cachesize")
	m.Port = strVal(obj, "port")
	m.Domainneeded = boolVal(obj, "domainneeded")
	m.Boguspriv = boolVal(obj, "boguspriv")
	m.Filterwin2k = boolVal(obj, "filterwin2k")
	m.Authoritative = boolVal(obj, "authoritative")
	m.Readethers = boolVal(obj, "readethers")
	m.Leasefile = strVal(obj, "leasefile")
	m.Resolvfile = strVal(obj, "resolvfile")
	m.Server = diags.list(listVal(ctx, obj, "server"))
	m.Address = diags.list(listVal(ctx, obj, "address"))
	m.Nonwildcard = boolVal(obj, "nonwildcard")
}

func (r *dhcpDnsmasqResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dhcpDnsmasqModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	body := r.body(ctx, plan, ds)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Patch(ctx, dhcpDnsmasqPath, body)
	if err != nil {
		resp.Diagnostics.AddError("Error configuring dnsmasq settings", err.Error())
		return
	}
	r.read(ctx, obj, &plan, ds)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dhcpDnsmasqResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dhcpDnsmasqModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := r.client.GetObject(ctx, dhcpDnsmasqPath)
	if err != nil {
		resp.Diagnostics.AddError("Error reading dnsmasq settings", err.Error())
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	r.read(ctx, obj, &state, ds)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *dhcpDnsmasqResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan dhcpDnsmasqModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	body := r.body(ctx, plan, ds)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Patch(ctx, dhcpDnsmasqPath, body)
	if err != nil {
		resp.Diagnostics.AddError("Error updating dnsmasq settings", err.Error())
		return
	}
	r.read(ctx, obj, &plan, ds)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete is a no-op: the dnsmasq singleton cannot be removed. State is dropped
// by the framework once this returns.
func (r *dhcpDnsmasqResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func (r *dhcpDnsmasqResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
