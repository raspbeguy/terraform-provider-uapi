package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &systemDataSource{}
	_ datasource.DataSourceWithConfigure = &systemDataSource{}
)

type systemDataSource struct{ client *client.Client }

func NewSystemDataSource() datasource.DataSource { return &systemDataSource{} }

func (d *systemDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system"
}

func (d *systemDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *systemDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "The global system settings (uci system.system).",
		Attributes: map[string]dsschema.Attribute{
			"id":           dsComputedString("Stable id of the system section."),
			"managed":      dsManagedAttribute(),
			"hostname":     dsComputedString("System hostname."),
			"description":  dsComputedString("Short device description."),
			"notes":        dsComputedString("Free-form notes."),
			"timezone":     dsComputedString("POSIX timezone string."),
			"zonename":     dsComputedString("IANA zone name (e.g. Europe/Paris)."),
			"log_size":     dsComputedString("Kernel log buffer size in KiB."),
			"log_ip":       dsComputedString("Remote syslog server IP."),
			"log_proto":    dsComputedString("Remote syslog protocol (udp or tcp)."),
			"log_remote":   dsComputedBool("Whether remote logging is enabled."),
			"urandom_seed": dsComputedBool("Whether a random seed is saved across reboots."),
		},
	}
}

func (d *systemDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	obj, found, err := d.client.GetObject(ctx, systemPath)
	if err != nil {
		resp.Diagnostics.AddError("Error reading system settings", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("System settings not found", "The system singleton is missing on the router")
		return
	}
	var m systemModel
	(&systemResource{}).read(ctx, obj, &m)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
