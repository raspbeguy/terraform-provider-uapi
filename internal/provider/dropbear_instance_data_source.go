package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &dropbearInstanceDataSource{}
	_ datasource.DataSourceWithConfigure = &dropbearInstanceDataSource{}
)

type dropbearInstanceDataSource struct{ client *client.Client }

func NewDropbearInstanceDataSource() datasource.DataSource { return &dropbearInstanceDataSource{} }

func (d *dropbearInstanceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dropbear_instance"
}

func (d *dropbearInstanceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *dropbearInstanceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up a dropbear instance by id.",
		Attributes: map[string]dsschema.Attribute{
			"id":                 dsIDAttribute(),
			"managed":            dsManagedAttribute(),
			"enable":             dsComputedBool("Whether this dropbear instance is enabled."),
			"port":               dsComputedString("TCP port to listen on."),
			"password_auth":      dsComputedBool("Whether password authentication is allowed."),
			"root_password_auth": dsComputedBool("Whether root password authentication is allowed."),
			"root_login":         dsComputedBool("Whether root logins are allowed."),
			"banner_file":        dsComputedString("Path to a file displayed before authentication."),
			"interface":          dsComputedString("Network interface to bind to."),
			"gateway_ports":      dsComputedBool("Whether remote hosts may connect to forwarded ports."),
		},
	}
}

func (d *dropbearInstanceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var m dropbearInstanceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := d.client.GetObject(ctx, "/"+dropbearInstanceCollection+"/"+m.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading dropbear instance", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("dropbear instance not found", "No dropbear instance with id "+m.ID.ValueString())
		return
	}
	(&dropbearInstanceResource{}).read(ctx, obj, &m)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
