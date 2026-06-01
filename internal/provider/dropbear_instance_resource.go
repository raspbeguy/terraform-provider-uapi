package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const dropbearInstanceCollection = "dropbear/instances"

var (
	_ resource.Resource                = &dropbearInstanceResource{}
	_ resource.ResourceWithConfigure   = &dropbearInstanceResource{}
	_ resource.ResourceWithImportState = &dropbearInstanceResource{}
)

type dropbearInstanceResource struct {
	client *client.Client
}

func NewDropbearInstanceResource() resource.Resource {
	return &dropbearInstanceResource{}
}

type dropbearInstanceModel struct {
	ID               types.String `tfsdk:"id"`
	Managed          types.Bool   `tfsdk:"managed"`
	ETag             types.String `tfsdk:"etag"`
	Enable           types.Bool   `tfsdk:"enable"`
	Port             types.String `tfsdk:"port"`
	PasswordAuth     types.Bool   `tfsdk:"password_auth"`
	RootPasswordAuth types.Bool   `tfsdk:"root_password_auth"`
	RootLogin        types.Bool   `tfsdk:"root_login"`
	BannerFile       types.String `tfsdk:"banner_file"`
	Interface        types.String `tfsdk:"interface"`
	GatewayPorts     types.Bool   `tfsdk:"gateway_ports"`
}

func (r *dropbearInstanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dropbear_instance"
}

func (r *dropbearInstanceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *dropbearInstanceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A dropbear SSH server instance (uci dropbear.dropbear).",
		Attributes: map[string]schema.Attribute{
			"id":                 computedIDAttribute(),
			"managed":            managedAttribute(),
			"etag":               etagAttribute(),
			"enable":             optionalComputedBool("Whether this dropbear instance is enabled. Defaults to true."),
			"port":               schema.StringAttribute{Optional: true, Description: "TCP port to listen on (1-65535)."},
			"password_auth":      optionalComputedBool("Allow password authentication. Defaults to true."),
			"root_password_auth": optionalComputedBool("Allow root password authentication. Defaults to true."),
			"root_login":         optionalComputedBool("Allow root logins. Defaults to true."),
			"banner_file":        schema.StringAttribute{Optional: true, Description: "Path to a file displayed before authentication."},
			"interface":          schema.StringAttribute{Optional: true, Description: "Network interface to bind to."},
			"gateway_ports":      optionalComputedBool("Allow remote hosts to connect to forwarded ports. Defaults to false."),
		},
	}
}

func (r *dropbearInstanceResource) body(m dropbearInstanceModel) map[string]any {
	out := map[string]any{}
	putBool(out, "enable", m.Enable)
	putStr(out, "Port", m.Port)
	putBool(out, "PasswordAuth", m.PasswordAuth)
	putBool(out, "RootPasswordAuth", m.RootPasswordAuth)
	putBool(out, "RootLogin", m.RootLogin)
	putStr(out, "BannerFile", m.BannerFile)
	putStr(out, "Interface", m.Interface)
	putBool(out, "GatewayPorts", m.GatewayPorts)
	return out
}

func (r *dropbearInstanceResource) read(_ context.Context, obj map[string]any, m *dropbearInstanceModel) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.Enable = boolVal(obj, "enable")
	m.Port = strVal(obj, "Port")
	m.PasswordAuth = boolVal(obj, "PasswordAuth")
	m.RootPasswordAuth = boolVal(obj, "RootPasswordAuth")
	m.RootLogin = boolVal(obj, "RootLogin")
	m.BannerFile = strVal(obj, "BannerFile")
	m.Interface = strVal(obj, "Interface")
	m.GatewayPorts = boolVal(obj, "GatewayPorts")
}

func (r *dropbearInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dropbearInstanceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, err := r.client.Post(ctx, "/"+dropbearInstanceCollection, r.body(plan), "")
	if err != nil {
		writeErr(&resp.Diagnostics, "creating", "dropbear instance", err)
		return
	}
	r.read(ctx, obj, &plan)
	plan.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dropbearInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dropbearInstanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := r.client.GetObject(ctx, "/"+dropbearInstanceCollection+"/"+state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading dropbear instance", err.Error())
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

func (r *dropbearInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state dropbearInstanceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, err := r.client.Put(ctx, "/"+dropbearInstanceCollection+"/"+plan.ID.ValueString(), r.body(plan), state.ETag.ValueString())
	if err != nil {
		writeErr(&resp.Diagnostics, "updating", "dropbear instance", err)
		return
	}
	r.read(ctx, obj, &plan)
	plan.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dropbearInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dropbearInstanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(ctx, "/"+dropbearInstanceCollection+"/"+state.ID.ValueString(), state.ETag.ValueString()); err != nil {
		writeErr(&resp.Diagnostics, "deleting", "dropbear instance", err)
	}
}

func (r *dropbearInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByID(ctx, r.client, dropbearInstanceCollection, "dropbear instance", req, resp)
}
