package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const uhttpdCertCollection = "uhttpd/certs"

var (
	_ resource.Resource                = &uhttpdCertResource{}
	_ resource.ResourceWithConfigure   = &uhttpdCertResource{}
	_ resource.ResourceWithImportState = &uhttpdCertResource{}
)

type uhttpdCertResource struct {
	client *client.Client
}

func NewUhttpdCertResource() resource.Resource {
	return &uhttpdCertResource{}
}

type uhttpdCertModel struct {
	ID           types.String `tfsdk:"id"`
	Managed      types.Bool   `tfsdk:"managed"`
	Days         types.String `tfsdk:"days"`
	Bits         types.String `tfsdk:"bits"`
	CommonName   types.String `tfsdk:"commonname"`
	Organization types.String `tfsdk:"organization"`
	Location     types.String `tfsdk:"location"`
	State        types.String `tfsdk:"state"`
	Country      types.String `tfsdk:"country"`
}

func (r *uhttpdCertResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_uhttpd_cert"
}

func (r *uhttpdCertResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *uhttpdCertResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A self-signed certificate generator for uhttpd (uci uhttpd.cert).",
		Attributes: map[string]schema.Attribute{
			"id":      computedIDAttribute(),
			"managed": managedAttribute(),
			"commonname": schema.StringAttribute{
				Required:    true,
				Description: "Certificate common name (CN).",
			},
			"days": schema.StringAttribute{
				Optional:    true,
				Description: "Certificate validity in days (1-36500).",
			},
			"bits": schema.StringAttribute{
				Optional:    true,
				Description: "RSA key size in bits (>= 1024).",
			},
			"organization": schema.StringAttribute{
				Optional:    true,
				Description: "Certificate organization (O).",
			},
			"location": schema.StringAttribute{
				Optional:    true,
				Description: "Certificate locality (L).",
			},
			"state": schema.StringAttribute{
				Optional:    true,
				Description: "Certificate state or province (ST).",
			},
			"country": schema.StringAttribute{
				Optional:    true,
				Description: "Two-letter country code (C).",
			},
		},
	}
}

func (r *uhttpdCertResource) body(m uhttpdCertModel) map[string]any {
	out := map[string]any{}
	putStr(out, "commonname", m.CommonName)
	putStr(out, "days", m.Days)
	putStr(out, "bits", m.Bits)
	putStr(out, "organization", m.Organization)
	putStr(out, "location", m.Location)
	putStr(out, "state", m.State)
	putStr(out, "country", m.Country)
	return out
}

func (r *uhttpdCertResource) read(_ context.Context, obj map[string]any, m *uhttpdCertModel) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.CommonName = strVal(obj, "commonname")
	m.Days = strVal(obj, "days")
	m.Bits = strVal(obj, "bits")
	m.Organization = strVal(obj, "organization")
	m.Location = strVal(obj, "location")
	m.State = strVal(obj, "state")
	m.Country = strVal(obj, "country")
}

func (r *uhttpdCertResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan uhttpdCertModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Post(ctx, "/"+uhttpdCertCollection, r.body(plan))
	if err != nil {
		resp.Diagnostics.AddError("Error creating uhttpd certificate", err.Error())
		return
	}
	r.read(ctx, obj, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *uhttpdCertResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state uhttpdCertModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := r.client.GetObject(ctx, "/"+uhttpdCertCollection+"/"+state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading uhttpd certificate", err.Error())
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	r.read(ctx, obj, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *uhttpdCertResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan uhttpdCertModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Put(ctx, "/"+uhttpdCertCollection+"/"+plan.ID.ValueString(), r.body(plan))
	if err != nil {
		resp.Diagnostics.AddError("Error updating uhttpd certificate", err.Error())
		return
	}
	r.read(ctx, obj, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *uhttpdCertResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state uhttpdCertModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(ctx, "/"+uhttpdCertCollection+"/"+state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting uhttpd certificate", err.Error())
	}
}

func (r *uhttpdCertResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByID(ctx, r.client, uhttpdCertCollection, "uhttpd certificate", req, resp)
}
