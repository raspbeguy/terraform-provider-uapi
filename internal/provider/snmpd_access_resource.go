package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const snmpdAccessCollection = "snmpd/accesses"

var (
	_ resource.Resource                = &snmpdAccessResource{}
	_ resource.ResourceWithConfigure   = &snmpdAccessResource{}
	_ resource.ResourceWithImportState = &snmpdAccessResource{}
)

type snmpdAccessResource struct {
	client *client.Client
}

func NewSnmpdAccessResource() resource.Resource {
	return &snmpdAccessResource{}
}

type snmpdAccessModel struct {
	ID      types.String `tfsdk:"id"`
	Managed types.Bool   `tfsdk:"managed"`
	Group   types.String `tfsdk:"group"`
	Context types.String `tfsdk:"context"`
	Version types.String `tfsdk:"version"`
	Level   types.String `tfsdk:"level"`
	Prefix  types.String `tfsdk:"prefix"`
	Read    types.String `tfsdk:"read"`
	Write   types.String `tfsdk:"write"`
	Notify  types.String `tfsdk:"notify"`
}

func (r *snmpdAccessResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snmpd_access"
}

func (r *snmpdAccessResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *snmpdAccessResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An SNMP VACM access entry (uci snmpd.access). Grants a group access to MIB views.",
		Attributes: map[string]schema.Attribute{
			"id":      computedIDAttribute(),
			"managed": managedAttribute(),
			"group": schema.StringAttribute{
				Required:    true,
				Description: "Name of the snmpd group this access entry applies to.",
			},
			"context": schema.StringAttribute{
				Optional:    true,
				Description: "SNMP context the access entry matches.",
			},
			"version": schema.StringAttribute{
				Optional:    true,
				Description: "Security model the entry matches: any, v1, v2c, or usm.",
			},
			"level": schema.StringAttribute{
				Optional:    true,
				Description: "Required security level: noauth, auth, or priv.",
			},
			"prefix": schema.StringAttribute{
				Optional:    true,
				Description: "Context match mode: exact or prefix.",
			},
			"read": schema.StringAttribute{
				Optional:    true,
				Description: "View name granted read access.",
			},
			"write": schema.StringAttribute{
				Optional:    true,
				Description: "View name granted write access.",
			},
			"notify": schema.StringAttribute{
				Optional:    true,
				Description: "View name granted notify access.",
			},
		},
	}
}

func (r *snmpdAccessResource) body(_ context.Context, m snmpdAccessModel) map[string]any {
	out := map[string]any{}
	putStr(out, "group", m.Group)
	putStr(out, "context", m.Context)
	putStr(out, "version", m.Version)
	putStr(out, "level", m.Level)
	putStr(out, "prefix", m.Prefix)
	putStr(out, "read", m.Read)
	putStr(out, "write", m.Write)
	putStr(out, "notify", m.Notify)
	return out
}

func (r *snmpdAccessResource) read(_ context.Context, obj map[string]any, m *snmpdAccessModel) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.Group = strVal(obj, "group")
	m.Context = strVal(obj, "context")
	m.Version = strVal(obj, "version")
	m.Level = strVal(obj, "level")
	m.Prefix = strVal(obj, "prefix")
	m.Read = strVal(obj, "read")
	m.Write = strVal(obj, "write")
	m.Notify = strVal(obj, "notify")
}

func (r *snmpdAccessResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan snmpdAccessModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Post(ctx, "/"+snmpdAccessCollection, r.body(ctx, plan))
	if err != nil {
		resp.Diagnostics.AddError("Error creating snmpd access", err.Error())
		return
	}
	r.read(ctx, obj, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *snmpdAccessResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state snmpdAccessModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := r.client.GetObject(ctx, "/"+snmpdAccessCollection+"/"+state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading snmpd access", err.Error())
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	r.read(ctx, obj, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *snmpdAccessResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan snmpdAccessModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Put(ctx, "/"+snmpdAccessCollection+"/"+plan.ID.ValueString(), r.body(ctx, plan))
	if err != nil {
		resp.Diagnostics.AddError("Error updating snmpd access", err.Error())
		return
	}
	r.read(ctx, obj, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *snmpdAccessResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state snmpdAccessModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(ctx, "/"+snmpdAccessCollection+"/"+state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting snmpd access", err.Error())
	}
}

func (r *snmpdAccessResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByID(ctx, r.client, snmpdAccessCollection, "snmpd access", req, resp)
}
