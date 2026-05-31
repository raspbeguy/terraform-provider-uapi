package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const sqmQueueCollection = "sqm/queues"

var (
	_ resource.Resource                = &sqmQueueResource{}
	_ resource.ResourceWithConfigure   = &sqmQueueResource{}
	_ resource.ResourceWithImportState = &sqmQueueResource{}
)

type sqmQueueResource struct {
	client *client.Client
}

func NewSqmQueueResource() resource.Resource {
	return &sqmQueueResource{}
}

type sqmQueueModel struct {
	ID        types.String `tfsdk:"id"`
	Managed   types.Bool   `tfsdk:"managed"`
	Enabled   types.Bool   `tfsdk:"enabled"`
	Interface types.String `tfsdk:"interface"`
	Download  types.String `tfsdk:"download"`
	Upload    types.String `tfsdk:"upload"`
	Qdisc     types.String `tfsdk:"qdisc"`
	Script    types.String `tfsdk:"script"`
	Linklayer types.String `tfsdk:"linklayer"`
	Overhead  types.String `tfsdk:"overhead"`
}

func (r *sqmQueueResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sqm_queue"
}

func (r *sqmQueueResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *sqmQueueResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An SQM (Smart Queue Management) queue (uci sqm.queue).",
		Attributes: map[string]schema.Attribute{
			"id":      computedIDAttribute(),
			"managed": managedAttribute(),
			"enabled": optionalComputedBool("Whether the queue is enabled. Defaults to true."),
			"interface": schema.StringAttribute{
				Required:    true,
				Description: "Network interface the queue is attached to.",
			},
			"download": schema.StringAttribute{
				Optional:    true,
				Description: "Download rate limit in kbit/s. 0 disables the download shaper.",
			},
			"upload": schema.StringAttribute{
				Optional:    true,
				Description: "Upload rate limit in kbit/s. 0 disables the upload shaper.",
			},
			"qdisc": schema.StringAttribute{
				Optional:    true,
				Description: "Queueing discipline: cake, fq_codel, pie, or htb.",
			},
			"script": schema.StringAttribute{
				Optional:    true,
				Description: "SQM script: piece_of_cake.qos, simple.qos, simplest.qos, or layer_cake.qos.",
			},
			"linklayer": schema.StringAttribute{
				Optional:    true,
				Description: "Link layer adaptation: none, ethernet, or atm.",
			},
			"overhead": schema.StringAttribute{
				Optional:    true,
				Description: "Per-packet overhead in bytes for link layer adaptation.",
			},
		},
	}
}

func (r *sqmQueueResource) body(_ context.Context, m sqmQueueModel) map[string]any {
	out := map[string]any{}
	putBool(out, "enabled", m.Enabled)
	putStr(out, "interface", m.Interface)
	putStr(out, "download", m.Download)
	putStr(out, "upload", m.Upload)
	putStr(out, "qdisc", m.Qdisc)
	putStr(out, "script", m.Script)
	putStr(out, "linklayer", m.Linklayer)
	putStr(out, "overhead", m.Overhead)
	return out
}

func (r *sqmQueueResource) read(_ context.Context, obj map[string]any, m *sqmQueueModel) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.Enabled = boolVal(obj, "enabled")
	m.Interface = strVal(obj, "interface")
	m.Download = strVal(obj, "download")
	m.Upload = strVal(obj, "upload")
	m.Qdisc = strVal(obj, "qdisc")
	m.Script = strVal(obj, "script")
	m.Linklayer = strVal(obj, "linklayer")
	m.Overhead = strVal(obj, "overhead")
}

func (r *sqmQueueResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sqmQueueModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Post(ctx, "/"+sqmQueueCollection, r.body(ctx, plan))
	if err != nil {
		resp.Diagnostics.AddError("Error creating SQM queue", err.Error())
		return
	}
	r.read(ctx, obj, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *sqmQueueResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sqmQueueModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := r.client.GetObject(ctx, "/"+sqmQueueCollection+"/"+state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading SQM queue", err.Error())
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	r.read(ctx, obj, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *sqmQueueResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sqmQueueModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Put(ctx, "/"+sqmQueueCollection+"/"+plan.ID.ValueString(), r.body(ctx, plan))
	if err != nil {
		resp.Diagnostics.AddError("Error updating SQM queue", err.Error())
		return
	}
	r.read(ctx, obj, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *sqmQueueResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sqmQueueModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(ctx, "/"+sqmQueueCollection+"/"+state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting SQM queue", err.Error())
	}
}

func (r *sqmQueueResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByID(ctx, r.client, sqmQueueCollection, "SQM queue", req, resp)
}
