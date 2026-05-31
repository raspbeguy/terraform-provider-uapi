package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const firewallForwardingCollection = "firewall/forwardings"

var (
	_ resource.Resource                = &firewallForwardingResource{}
	_ resource.ResourceWithConfigure   = &firewallForwardingResource{}
	_ resource.ResourceWithImportState = &firewallForwardingResource{}
)

type firewallForwardingResource struct {
	client *client.Client
}

func NewFirewallForwardingResource() resource.Resource {
	return &firewallForwardingResource{}
}

type firewallForwardingModel struct {
	ID      types.String `tfsdk:"id"`
	Managed types.Bool   `tfsdk:"managed"`
	Src     types.String `tfsdk:"src"`
	Dest    types.String `tfsdk:"dest"`
	Family  types.String `tfsdk:"family"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

func (r *firewallForwardingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_forwarding"
}

func (r *firewallForwardingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *firewallForwardingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A firewall forwarding between two zones (uci firewall.forwarding).",
		Attributes: map[string]schema.Attribute{
			"id":      computedIDAttribute(),
			"managed": managedAttribute(),
			"src": schema.StringAttribute{
				Required:    true,
				Description: "Source zone name.",
			},
			"dest": schema.StringAttribute{
				Required:    true,
				Description: "Destination zone name.",
			},
			"family":  optionalComputedString("Address family: any, ipv4, or ipv6. Defaults to any."),
			"enabled": optionalComputedBool("Whether the forwarding is active. Defaults to true."),
		},
	}
}

func (r *firewallForwardingResource) body(_ context.Context, m firewallForwardingModel) map[string]any {
	out := map[string]any{}
	putStr(out, "src", m.Src)
	putStr(out, "dest", m.Dest)
	putStr(out, "family", m.Family)
	putBool(out, "enabled", m.Enabled)
	return out
}

func (r *firewallForwardingResource) read(_ context.Context, obj map[string]any, m *firewallForwardingModel) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.Src = strVal(obj, "src")
	m.Dest = strVal(obj, "dest")
	m.Family = strVal(obj, "family")
	m.Enabled = boolVal(obj, "enabled")
}

func (r *firewallForwardingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan firewallForwardingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Post(ctx, "/"+firewallForwardingCollection, r.body(ctx, plan))
	if err != nil {
		resp.Diagnostics.AddError("Error creating firewall forwarding", err.Error())
		return
	}
	r.read(ctx, obj, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *firewallForwardingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state firewallForwardingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := r.client.GetObject(ctx, "/"+firewallForwardingCollection+"/"+state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading firewall forwarding", err.Error())
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	r.read(ctx, obj, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *firewallForwardingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan firewallForwardingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Put(ctx, "/"+firewallForwardingCollection+"/"+plan.ID.ValueString(), r.body(ctx, plan))
	if err != nil {
		resp.Diagnostics.AddError("Error updating firewall forwarding", err.Error())
		return
	}
	r.read(ctx, obj, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *firewallForwardingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state firewallForwardingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(ctx, "/"+firewallForwardingCollection+"/"+state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting firewall forwarding", err.Error())
	}
}

func (r *firewallForwardingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByID(ctx, r.client, firewallForwardingCollection, "firewall forwarding", req, resp)
}
