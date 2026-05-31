package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &uhttpdCertDataSource{}
	_ datasource.DataSourceWithConfigure = &uhttpdCertDataSource{}
)

type uhttpdCertDataSource struct{ client *client.Client }

func NewUhttpdCertDataSource() datasource.DataSource { return &uhttpdCertDataSource{} }

func (d *uhttpdCertDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_uhttpd_cert"
}

func (d *uhttpdCertDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *uhttpdCertDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up a uhttpd certificate by id.",
		Attributes: map[string]dsschema.Attribute{
			"id":           dsIDAttribute(),
			"managed":      dsManagedAttribute(),
			"commonname":   dsComputedString("Certificate common name (CN)."),
			"days":         dsComputedString("Certificate validity in days."),
			"bits":         dsComputedString("RSA key size in bits."),
			"organization": dsComputedString("Certificate organization (O)."),
			"location":     dsComputedString("Certificate locality (L)."),
			"state":        dsComputedString("Certificate state or province (ST)."),
			"country":      dsComputedString("Two-letter country code (C)."),
		},
	}
}

func (d *uhttpdCertDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var m uhttpdCertModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := d.client.GetObject(ctx, "/"+uhttpdCertCollection+"/"+m.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading uhttpd certificate", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("uhttpd certificate not found", "No uhttpd certificate with id "+m.ID.ValueString())
		return
	}
	(&uhttpdCertResource{}).read(ctx, obj, &m)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
