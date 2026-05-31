package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const snmpdGroupCollection = "snmpd/groups"

var (
	_ resource.Resource                = &snmpdGroupResource{}
	_ resource.ResourceWithConfigure   = &snmpdGroupResource{}
	_ resource.ResourceWithImportState = &snmpdGroupResource{}
)

type snmpdGroupResource struct {
	client *client.Client
}

func NewSnmpdGroupResource() resource.Resource {
	return &snmpdGroupResource{}
}

type snmpdGroupModel struct {
	ID      types.String `tfsdk:"id"`
	Managed types.Bool   `tfsdk:"managed"`
	Group   types.String `tfsdk:"group"`
	Version types.String `tfsdk:"version"`
	Secname types.String `tfsdk:"secname"`
}

func (r *snmpdGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snmpd_group"
}

func (r *snmpdGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *snmpdGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An SNMP VACM group definition (uci snmpd.group). Binds a security name to a group.",
		Attributes: map[string]schema.Attribute{
			"id":      computedIDAttribute(),
			"managed": managedAttribute(),
			"group": schema.StringAttribute{
				Required:    true,
				Description: "Group name referenced by snmpd access entries.",
			},
			"version": schema.StringAttribute{
				Optional:    true,
				Description: "Security model for the membership: v1, v2c, or usm.",
			},
			"secname": schema.StringAttribute{
				Optional:    true,
				Description: "Security name added to the group.",
			},
		},
	}
}

func (r *snmpdGroupResource) body(_ context.Context, m snmpdGroupModel) map[string]any {
	out := map[string]any{}
	putStr(out, "group", m.Group)
	putStr(out, "version", m.Version)
	putStr(out, "secname", m.Secname)
	return out
}

func (r *snmpdGroupResource) read(_ context.Context, obj map[string]any, m *snmpdGroupModel) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.Group = strVal(obj, "group")
	m.Version = strVal(obj, "version")
	m.Secname = strVal(obj, "secname")
}

func (r *snmpdGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan snmpdGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Post(ctx, "/"+snmpdGroupCollection, r.body(ctx, plan))
	if err != nil {
		resp.Diagnostics.AddError("Error creating snmpd group", err.Error())
		return
	}
	r.read(ctx, obj, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *snmpdGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state snmpdGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := r.client.GetObject(ctx, "/"+snmpdGroupCollection+"/"+state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading snmpd group", err.Error())
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	r.read(ctx, obj, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *snmpdGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan snmpdGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Put(ctx, "/"+snmpdGroupCollection+"/"+plan.ID.ValueString(), r.body(ctx, plan))
	if err != nil {
		resp.Diagnostics.AddError("Error updating snmpd group", err.Error())
		return
	}
	r.read(ctx, obj, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *snmpdGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state snmpdGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(ctx, "/"+snmpdGroupCollection+"/"+state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting snmpd group", err.Error())
	}
}

func (r *snmpdGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByID(ctx, r.client, snmpdGroupCollection, "snmpd group", req, resp)
}
