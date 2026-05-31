package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const vnstatInterfaceCollection = "vnstat/interfaces"

var (
	_ resource.Resource                = &vnstatInterfaceResource{}
	_ resource.ResourceWithConfigure   = &vnstatInterfaceResource{}
	_ resource.ResourceWithImportState = &vnstatInterfaceResource{}
)

type vnstatInterfaceResource struct {
	client *client.Client
}

func NewVnstatInterfaceResource() resource.Resource {
	return &vnstatInterfaceResource{}
}

type vnstatInterfaceModel struct {
	ID        types.String `tfsdk:"id"`
	Managed   types.Bool   `tfsdk:"managed"`
	Interface types.String `tfsdk:"interface"`
	Enabled   types.Bool   `tfsdk:"enabled"`
}

func (r *vnstatInterfaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vnstat_interface"
}

func (r *vnstatInterfaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *vnstatInterfaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A vnstat monitored interface (uci vnstat.interface).",
		Attributes: map[string]schema.Attribute{
			"id":      computedIDAttribute(),
			"managed": managedAttribute(),
			"interface": schema.StringAttribute{
				Required:    true,
				Description: "Name of the network interface to monitor. Must reference an existing network interface.",
			},
			"enabled": optionalComputedBool("Whether monitoring of this interface is enabled. Defaults to true."),
		},
	}
}

func (r *vnstatInterfaceResource) body(_ context.Context, m vnstatInterfaceModel) map[string]any {
	out := map[string]any{}
	putStr(out, "interface", m.Interface)
	putBool(out, "enabled", m.Enabled)
	return out
}

func (r *vnstatInterfaceResource) read(_ context.Context, obj map[string]any, m *vnstatInterfaceModel) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.Interface = strVal(obj, "interface")
	m.Enabled = boolVal(obj, "enabled")
}

func (r *vnstatInterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vnstatInterfaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Post(ctx, "/"+vnstatInterfaceCollection, r.body(ctx, plan))
	if err != nil {
		resp.Diagnostics.AddError("Error creating vnstat interface", err.Error())
		return
	}
	r.read(ctx, obj, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *vnstatInterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vnstatInterfaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := r.client.GetObject(ctx, "/"+vnstatInterfaceCollection+"/"+state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading vnstat interface", err.Error())
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	r.read(ctx, obj, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *vnstatInterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan vnstatInterfaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Put(ctx, "/"+vnstatInterfaceCollection+"/"+plan.ID.ValueString(), r.body(ctx, plan))
	if err != nil {
		resp.Diagnostics.AddError("Error updating vnstat interface", err.Error())
		return
	}
	r.read(ctx, obj, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *vnstatInterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vnstatInterfaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(ctx, "/"+vnstatInterfaceCollection+"/"+state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting vnstat interface", err.Error())
	}
}

func (r *vnstatInterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByID(ctx, r.client, vnstatInterfaceCollection, "vnstat interface", req, resp)
}
