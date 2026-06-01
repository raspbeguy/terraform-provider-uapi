package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const vnstatConfigPath = "/vnstat/config"

var (
	_ resource.Resource                = &vnstatConfigResource{}
	_ resource.ResourceWithConfigure   = &vnstatConfigResource{}
	_ resource.ResourceWithImportState = &vnstatConfigResource{}
)

type vnstatConfigResource struct {
	client *client.Client
}

func NewVnstatConfigResource() resource.Resource {
	return &vnstatConfigResource{}
}

type vnstatConfigModel struct {
	ID                 types.String `tfsdk:"id"`
	Managed            types.Bool   `tfsdk:"managed"`
	ETag               types.String `tfsdk:"etag"`
	DatabaseDir        types.String `tfsdk:"database_dir"`
	Interface5MinHours types.String `tfsdk:"interface_5min_hours"`
	MonthRotate        types.String `tfsdk:"month_rotate"`
}

func (r *vnstatConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vnstat_config"
}

func (r *vnstatConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *vnstatConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Global vnstat daemon settings (uci vnstat.vnstat). This is a singleton: it cannot be " +
			"created or destroyed. `terraform destroy` only removes it from state; the underlying " +
			"settings are left as-is on the router.",
		Attributes: map[string]schema.Attribute{
			"id":                   computedIDAttribute(),
			"managed":              managedAttribute(),
			"etag":                 etagAttribute(),
			"database_dir":         schema.StringAttribute{Optional: true, Description: "Directory where vnstat stores its databases."},
			"interface_5min_hours": schema.StringAttribute{Optional: true, Description: "Hours of 5-minute resolution data to keep per interface."},
			"month_rotate":         schema.StringAttribute{Optional: true, Description: "Day of the month on which monthly statistics are reset."},
		},
	}
}

func (r *vnstatConfigResource) body(_ context.Context, m vnstatConfigModel) map[string]any {
	out := map[string]any{}
	putStr(out, "DatabaseDir", m.DatabaseDir)
	putStr(out, "Interface5MinHours", m.Interface5MinHours)
	putStr(out, "MonthRotate", m.MonthRotate)
	return out
}

func (r *vnstatConfigResource) read(_ context.Context, obj map[string]any, m *vnstatConfigModel) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.DatabaseDir = strVal(obj, "DatabaseDir")
	m.Interface5MinHours = strVal(obj, "Interface5MinHours")
	m.MonthRotate = strVal(obj, "MonthRotate")
}

func (r *vnstatConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vnstatConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, err := r.client.Patch(ctx, vnstatConfigPath, r.body(ctx, plan), "")
	if err != nil {
		writeErr(&resp.Diagnostics, "configuring", "vnstat settings", err)
		return
	}
	r.read(ctx, obj, &plan)
	plan.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *vnstatConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vnstatConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := r.client.GetObject(ctx, vnstatConfigPath)
	if err != nil {
		resp.Diagnostics.AddError("Error reading vnstat settings", err.Error())
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

func (r *vnstatConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state vnstatConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, err := r.client.Patch(ctx, vnstatConfigPath, r.body(ctx, plan), state.ETag.ValueString())
	if err != nil {
		writeErr(&resp.Diagnostics, "updating", "vnstat settings", err)
		return
	}
	r.read(ctx, obj, &plan)
	plan.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete is a no-op: the vnstat config singleton cannot be removed. State is
// dropped by the framework once this returns.
func (r *vnstatConfigResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func (r *vnstatConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
