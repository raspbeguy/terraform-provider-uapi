package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &prometheusNodeExporterLuaConfigDataSource{}
	_ datasource.DataSourceWithConfigure = &prometheusNodeExporterLuaConfigDataSource{}
)

type prometheusNodeExporterLuaConfigDataSource struct{ client *client.Client }

func NewPrometheusNodeExporterLuaConfigDataSource() datasource.DataSource {
	return &prometheusNodeExporterLuaConfigDataSource{}
}

func (d *prometheusNodeExporterLuaConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prometheus_node_exporter_lua_config"
}

func (d *prometheusNodeExporterLuaConfigDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *prometheusNodeExporterLuaConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "The global prometheus-node-exporter-lua settings (uci prometheus-node-exporter-lua).",
		Attributes: map[string]dsschema.Attribute{
			"id":               dsComputedString("Stable id of the prometheus-node-exporter-lua config section."),
			"managed":          dsManagedAttribute(),
			"listen_ipv6":      dsComputedBool("Whether the exporter listens on IPv6."),
			"listen_interface": dsComputedString("Network interface the exporter is bound to."),
			"listen_port":      dsComputedString("TCP port the exporter listens on."),
			"cpu":              dsComputedBool("Whether the cpu collector is enabled."),
			"meminfo":          dsComputedBool("Whether the meminfo collector is enabled."),
			"netdev":           dsComputedBool("Whether the netdev collector is enabled."),
			"loadavg":          dsComputedBool("Whether the loadavg collector is enabled."),
			"filesystem":       dsComputedBool("Whether the filesystem collector is enabled."),
			"diskstats":        dsComputedBool("Whether the diskstats collector is enabled."),
			"uname":            dsComputedBool("Whether the uname collector is enabled."),
			"netstat":          dsComputedBool("Whether the netstat collector is enabled."),
			"stat":             dsComputedBool("Whether the stat collector is enabled."),
			"vmstat":           dsComputedBool("Whether the vmstat collector is enabled."),
			"boottime":         dsComputedBool("Whether the boottime collector is enabled."),
			"entropy":          dsComputedBool("Whether the entropy collector is enabled."),
			"time":             dsComputedBool("Whether the time collector is enabled."),
			"hwmon":            dsComputedBool("Whether the hwmon collector is enabled."),
			"textfile":         dsComputedBool("Whether the textfile collector is enabled."),
			"thermal_zone":     dsComputedBool("Whether the thermal_zone collector is enabled."),
			"edac":             dsComputedBool("Whether the edac collector is enabled."),
		},
	}
}

func (d *prometheusNodeExporterLuaConfigDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	obj, found, err := d.client.GetObject(ctx, prometheusNodeExporterLuaConfigPath)
	if err != nil {
		resp.Diagnostics.AddError("Error reading prometheus-node-exporter-lua settings", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Prometheus-node-exporter-lua settings not found", "The prometheus-node-exporter-lua config singleton is missing on the router")
		return
	}
	var m prometheusNodeExporterLuaConfigModel
	(&prometheusNodeExporterLuaConfigResource{}).read(ctx, obj, &m)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
