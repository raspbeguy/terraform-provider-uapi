package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &uhttpdInstanceDataSource{}
	_ datasource.DataSourceWithConfigure = &uhttpdInstanceDataSource{}
)

type uhttpdInstanceDataSource struct{ client *client.Client }

func NewUhttpdInstanceDataSource() datasource.DataSource { return &uhttpdInstanceDataSource{} }

func (d *uhttpdInstanceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_uhttpd_instance"
}

func (d *uhttpdInstanceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *uhttpdInstanceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up a uhttpd instance by id.",
		Attributes: map[string]dsschema.Attribute{
			"id":              dsIDAttribute(),
			"managed":         dsManagedAttribute(),
			"etag":            dsComputedString("Opaque ETag of the resource's current state."),
			"listen_http":     dsComputedStringList("Addresses listened on for HTTP."),
			"listen_https":    dsComputedStringList("Addresses listened on for HTTPS."),
			"home":            dsComputedString("Document root served by this instance."),
			"cert":            dsComputedString("Path to the TLS certificate."),
			"key":             dsComputedString("Path to the TLS private key."),
			"cgi_prefix":      dsComputedString("URL prefix mapped to CGI scripts."),
			"lua_prefix":      dsComputedStringList("Lua handler prefixes."),
			"ucode_prefix":    dsComputedStringList("ucode handler prefixes."),
			"max_requests":    dsComputedString("Maximum number of concurrent requests."),
			"max_connections": dsComputedString("Maximum number of concurrent connections."),
			"script_timeout":  dsComputedString("Maximum seconds a CGI/Lua/ucode script may run."),
			"network_timeout": dsComputedString("Maximum seconds to wait for network activity."),
			"http_keepalive":  dsComputedString("HTTP keep-alive timeout in seconds."),
			"tcp_keepalive":   dsComputedString("TCP keep-alive interval in seconds."),
			"index_page":      dsComputedStringList("Index file names tried for directory requests."),
			"error_page":      dsComputedString("Virtual URL or CGI script handling error pages."),
			"no_dirlists":     dsComputedBool("Whether directory listings are disabled."),
			"no_symlinks":     dsComputedBool("Whether symbolic links are not followed."),
			"rfc1918_filter":  dsComputedBool("Whether public-to-private (RFC1918) requests are rejected."),
		},
	}
}

func (d *uhttpdInstanceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var m uhttpdInstanceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := d.client.GetObject(ctx, "/"+uhttpdInstanceCollection+"/"+m.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading uhttpd instance", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("uhttpd instance not found", "No uhttpd instance with id "+m.ID.ValueString())
		return
	}
	ds := newDiagsink(&resp.Diagnostics)
	(&uhttpdInstanceResource{}).read(ctx, obj, &m, ds)
	m.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
