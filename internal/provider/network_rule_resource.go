package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const networkRuleCollection = "network/rules"

var (
	_ resource.Resource                = &networkRuleResource{}
	_ resource.ResourceWithConfigure   = &networkRuleResource{}
	_ resource.ResourceWithImportState = &networkRuleResource{}
)

type networkRuleResource struct {
	client *client.Client
}

func NewNetworkRuleResource() resource.Resource {
	return &networkRuleResource{}
}

type networkRuleModel struct {
	ID       types.String `tfsdk:"id"`
	Managed  types.Bool   `tfsdk:"managed"`
	In       types.String `tfsdk:"in"`
	Out      types.String `tfsdk:"out"`
	Src      types.String `tfsdk:"src"`
	Dest     types.String `tfsdk:"dest"`
	Priority types.String `tfsdk:"priority"`
	Lookup   types.String `tfsdk:"lookup"`
	Goto     types.String `tfsdk:"goto"`
	Action   types.String `tfsdk:"action"`
	Invert   types.Bool   `tfsdk:"invert"`
	Mark     types.String `tfsdk:"mark"`
}

func (r *networkRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_rule"
}

func (r *networkRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *networkRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An IP routing policy rule (uci network.rule). At least one of in, out, src, or dest must be set.",
		Attributes: map[string]schema.Attribute{
			"id":      computedIDAttribute(),
			"managed": managedAttribute(),
			"in": schema.StringAttribute{
				Optional:    true,
				Description: "Incoming interface selector.",
			},
			"out": schema.StringAttribute{
				Optional:    true,
				Description: "Outgoing interface selector.",
			},
			"src": schema.StringAttribute{
				Optional:    true,
				Description: "Source IPv4 address or CIDR selector.",
			},
			"dest": schema.StringAttribute{
				Optional:    true,
				Description: "Destination IPv4 address or CIDR selector.",
			},
			"priority": schema.StringAttribute{
				Optional:    true,
				Description: "Rule priority (0-32766).",
			},
			"lookup": schema.StringAttribute{
				Optional:    true,
				Description: "Routing table to look up. Required when action is lookup.",
			},
			"goto": schema.StringAttribute{
				Optional:    true,
				Description: "Priority to jump to. Required when action is goto.",
			},
			"action": optionalComputedString("Rule action: lookup, goto, unreachable, prohibit, blackhole, or throw. Defaults to lookup."),
			"invert": optionalComputedBool("Invert the rule selectors. Defaults to false."),
			"mark": schema.StringAttribute{
				Optional:    true,
				Description: "Firewall mark to match.",
			},
		},
	}
}

func (r *networkRuleResource) body(_ context.Context, m networkRuleModel) map[string]any {
	out := map[string]any{}
	putStr(out, "in", m.In)
	putStr(out, "out", m.Out)
	putStr(out, "src", m.Src)
	putStr(out, "dest", m.Dest)
	putStr(out, "priority", m.Priority)
	putStr(out, "lookup", m.Lookup)
	putStr(out, "goto", m.Goto)
	putStr(out, "action", m.Action)
	putBool(out, "invert", m.Invert)
	putStr(out, "mark", m.Mark)
	return out
}

func (r *networkRuleResource) read(_ context.Context, obj map[string]any, m *networkRuleModel) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.In = strVal(obj, "in")
	m.Out = strVal(obj, "out")
	m.Src = strVal(obj, "src")
	m.Dest = strVal(obj, "dest")
	m.Priority = strVal(obj, "priority")
	m.Lookup = strVal(obj, "lookup")
	m.Goto = strVal(obj, "goto")
	m.Action = strVal(obj, "action")
	m.Invert = boolVal(obj, "invert")
	m.Mark = strVal(obj, "mark")
}

func (r *networkRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan networkRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Post(ctx, "/"+networkRuleCollection, r.body(ctx, plan))
	if err != nil {
		resp.Diagnostics.AddError("Error creating network rule", err.Error())
		return
	}
	r.read(ctx, obj, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state networkRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := r.client.GetObject(ctx, "/"+networkRuleCollection+"/"+state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading network rule", err.Error())
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	r.read(ctx, obj, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *networkRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan networkRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Put(ctx, "/"+networkRuleCollection+"/"+plan.ID.ValueString(), r.body(ctx, plan))
	if err != nil {
		resp.Diagnostics.AddError("Error updating network rule", err.Error())
		return
	}
	r.read(ctx, obj, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state networkRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(ctx, "/"+networkRuleCollection+"/"+state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting network rule", err.Error())
	}
}

func (r *networkRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByID(ctx, r.client, networkRuleCollection, "network rule", req, resp)
}
