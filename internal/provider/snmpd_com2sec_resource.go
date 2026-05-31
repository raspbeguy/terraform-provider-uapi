package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const snmpdCom2secCollection = "snmpd/com2secs"

var (
	_ resource.Resource                = &snmpdCom2secResource{}
	_ resource.ResourceWithConfigure   = &snmpdCom2secResource{}
	_ resource.ResourceWithImportState = &snmpdCom2secResource{}
)

type snmpdCom2secResource struct {
	client *client.Client
}

func NewSnmpdCom2secResource() resource.Resource {
	return &snmpdCom2secResource{}
}

type snmpdCom2secModel struct {
	ID        types.String `tfsdk:"id"`
	Managed   types.Bool   `tfsdk:"managed"`
	Secname   types.String `tfsdk:"secname"`
	Source    types.String `tfsdk:"source"`
	Community types.String `tfsdk:"community"`
}

func (r *snmpdCom2secResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snmpd_com2sec"
}

func (r *snmpdCom2secResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *snmpdCom2secResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An SNMP community-to-security-name mapping (uci snmpd.com2sec).",
		Attributes: map[string]schema.Attribute{
			"id":      computedIDAttribute(),
			"managed": managedAttribute(),
			"secname": schema.StringAttribute{
				Required:    true,
				Description: "Security name the community maps to.",
			},
			"source": schema.StringAttribute{
				Required:    true,
				Description: "Source network range or 'default'.",
			},
			"community": schema.StringAttribute{
				Required:    true,
				Description: "SNMP community string.",
			},
		},
	}
}

func (r *snmpdCom2secResource) body(_ context.Context, m snmpdCom2secModel) map[string]any {
	out := map[string]any{}
	putStr(out, "secname", m.Secname)
	putStr(out, "source", m.Source)
	putStr(out, "community", m.Community)
	return out
}

func (r *snmpdCom2secResource) read(_ context.Context, obj map[string]any, m *snmpdCom2secModel) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.Secname = strVal(obj, "secname")
	m.Source = strVal(obj, "source")
	m.Community = strVal(obj, "community")
}

func (r *snmpdCom2secResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan snmpdCom2secModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Post(ctx, "/"+snmpdCom2secCollection, r.body(ctx, plan))
	if err != nil {
		resp.Diagnostics.AddError("Error creating snmpd com2sec", err.Error())
		return
	}
	r.read(ctx, obj, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *snmpdCom2secResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state snmpdCom2secModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := r.client.GetObject(ctx, "/"+snmpdCom2secCollection+"/"+state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading snmpd com2sec", err.Error())
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	r.read(ctx, obj, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *snmpdCom2secResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan snmpdCom2secModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Put(ctx, "/"+snmpdCom2secCollection+"/"+plan.ID.ValueString(), r.body(ctx, plan))
	if err != nil {
		resp.Diagnostics.AddError("Error updating snmpd com2sec", err.Error())
		return
	}
	r.read(ctx, obj, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *snmpdCom2secResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state snmpdCom2secModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(ctx, "/"+snmpdCom2secCollection+"/"+state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting snmpd com2sec", err.Error())
	}
}

func (r *snmpdCom2secResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByID(ctx, r.client, snmpdCom2secCollection, "snmpd com2sec", req, resp)
}
