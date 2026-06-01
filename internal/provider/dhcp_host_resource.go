package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const dhcpHostCollection = "dhcp/hosts"

var (
	_ resource.Resource                = &dhcpHostResource{}
	_ resource.ResourceWithConfigure   = &dhcpHostResource{}
	_ resource.ResourceWithImportState = &dhcpHostResource{}
)

type dhcpHostResource struct {
	client *client.Client
}

func NewDHCPHostResource() resource.Resource {
	return &dhcpHostResource{}
}

type dhcpHostModel struct {
	ID         types.String `tfsdk:"id"`
	Managed    types.Bool   `tfsdk:"managed"`
	ETag       types.String `tfsdk:"etag"`
	Name       types.String `tfsdk:"name"`
	MAC        types.String `tfsdk:"mac"`
	MACAliases types.List   `tfsdk:"mac_aliases"`
	DUID       types.String `tfsdk:"duid"`
	HostID     types.String `tfsdk:"hostid"`
	IP         types.String `tfsdk:"ip"`
	Leasetime  types.String `tfsdk:"leasetime"`
	Tag        types.String `tfsdk:"tag"`
	DNS        types.Bool   `tfsdk:"dns"`
	Broadcast  types.Bool   `tfsdk:"broadcast"`
	Instance   types.String `tfsdk:"instance"`
}

func (r *dhcpHostResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcp_host"
}

func (r *dhcpHostResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *dhcpHostResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A static DHCP lease (uci dhcp.host).",
		Attributes: map[string]schema.Attribute{
			"id":      computedIDAttribute(),
			"managed": managedAttribute(),
			"etag":    etagAttribute(),
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Hostname for the static lease.",
			},
			"mac": schema.StringAttribute{
				Optional:    true,
				Description: "Primary client MAC address for an IPv4 reservation (aa:bb:cc:dd:ee:ff). Either mac or duid is required.",
			},
			"mac_aliases": optionalComputedStringList("Additional MAC addresses sharing the same reservation."),
			"duid": schema.StringAttribute{
				Optional:    true,
				Description: "Client DUID for a DHCPv6 reservation. Either mac or duid is required.",
			},
			"hostid": schema.StringAttribute{
				Optional:    true,
				Description: "Static IPv6 host id hint (suffix), like '::42'.",
			},
			"ip": schema.StringAttribute{
				Required:    true,
				Description: "IPv4 or IPv6 address to assign.",
			},
			"leasetime": optionalComputedString("Lease duration like '12h', '30m', '1d', or seconds."),
			"tag":       optionalComputedString("dnsmasq tag to apply to the host."),
			"dns":       optionalComputedBool("Add a DNS entry for the host. Defaults to false."),
			"broadcast": optionalComputedBool("Force broadcast replies for clients that need it. Defaults to false."),
			"instance": schema.StringAttribute{
				Optional:    true,
				Description: "Pin this reservation to a specific dhcp/dnsmasq instance (section name).",
			},
		},
	}
}

func (r *dhcpHostResource) body(ctx context.Context, m dhcpHostModel, diags *diagsink) map[string]any {
	out := map[string]any{}
	putStr(out, "name", m.Name)
	putStr(out, "mac", m.MAC)
	putList(ctx, out, "mac_aliases", m.MACAliases, diags.d)
	putStr(out, "duid", m.DUID)
	putStr(out, "hostid", m.HostID)
	putStr(out, "ip", m.IP)
	putStr(out, "leasetime", m.Leasetime)
	putStr(out, "tag", m.Tag)
	putBool(out, "dns", m.DNS)
	putBool(out, "broadcast", m.Broadcast)
	putStr(out, "instance", m.Instance)
	return out
}

func (r *dhcpHostResource) read(ctx context.Context, obj map[string]any, m *dhcpHostModel, diags *diagsink) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.Name = strVal(obj, "name")
	m.MAC = strVal(obj, "mac")
	m.MACAliases = diags.list(listVal(ctx, obj, "mac_aliases"))
	m.DUID = strVal(obj, "duid")
	m.HostID = strVal(obj, "hostid")
	m.IP = strVal(obj, "ip")
	m.Leasetime = strVal(obj, "leasetime")
	m.Tag = strVal(obj, "tag")
	m.DNS = boolVal(obj, "dns")
	m.Broadcast = boolVal(obj, "broadcast")
	m.Instance = strVal(obj, "instance")
}

func (r *dhcpHostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dhcpHostModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	body := r.body(ctx, plan, ds)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, err := r.client.Post(ctx, "/"+dhcpHostCollection, body, "")
	if err != nil {
		writeErr(&resp.Diagnostics, "creating", "dhcp host", err)
		return
	}
	r.read(ctx, obj, &plan, ds)
	plan.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dhcpHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dhcpHostModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := r.client.GetObject(ctx, "/"+dhcpHostCollection+"/"+state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading dhcp host", err.Error())
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	r.read(ctx, obj, &state, ds)
	state.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *dhcpHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state dhcpHostModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	body := r.body(ctx, plan, ds)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, err := r.client.Put(ctx, "/"+dhcpHostCollection+"/"+plan.ID.ValueString(), body, state.ETag.ValueString())
	if err != nil {
		writeErr(&resp.Diagnostics, "updating", "dhcp host", err)
		return
	}
	r.read(ctx, obj, &plan, ds)
	plan.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dhcpHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dhcpHostModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(ctx, "/"+dhcpHostCollection+"/"+state.ID.ValueString(), state.ETag.ValueString()); err != nil {
		writeErr(&resp.Diagnostics, "deleting", "dhcp host", err)
	}
}

func (r *dhcpHostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByID(ctx, r.client, dhcpHostCollection, "dhcp host", req, resp)
}
