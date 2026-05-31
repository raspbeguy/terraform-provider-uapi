package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

const prometheusNodeExporterLuaConfigPath = "/prometheus_node_exporter_lua/config"

var (
	_ resource.Resource                = &prometheusNodeExporterLuaConfigResource{}
	_ resource.ResourceWithConfigure   = &prometheusNodeExporterLuaConfigResource{}
	_ resource.ResourceWithImportState = &prometheusNodeExporterLuaConfigResource{}
)

type prometheusNodeExporterLuaConfigResource struct {
	client *client.Client
}

func NewPrometheusNodeExporterLuaConfigResource() resource.Resource {
	return &prometheusNodeExporterLuaConfigResource{}
}

type prometheusNodeExporterLuaConfigModel struct {
	ID              types.String `tfsdk:"id"`
	Managed         types.Bool   `tfsdk:"managed"`
	ListenIPv6      types.Bool   `tfsdk:"listen_ipv6"`
	ListenInterface types.String `tfsdk:"listen_interface"`
	ListenPort      types.String `tfsdk:"listen_port"`
	CPU             types.Bool   `tfsdk:"cpu"`
	Meminfo         types.Bool   `tfsdk:"meminfo"`
	Netdev          types.Bool   `tfsdk:"netdev"`
	Loadavg         types.Bool   `tfsdk:"loadavg"`
	Filesystem      types.Bool   `tfsdk:"filesystem"`
	Diskstats       types.Bool   `tfsdk:"diskstats"`
	Uname           types.Bool   `tfsdk:"uname"`
	Netstat         types.Bool   `tfsdk:"netstat"`
	Stat            types.Bool   `tfsdk:"stat"`
	Vmstat          types.Bool   `tfsdk:"vmstat"`
	Boottime        types.Bool   `tfsdk:"boottime"`
	Entropy         types.Bool   `tfsdk:"entropy"`
	Time            types.Bool   `tfsdk:"time"`
	Hwmon           types.Bool   `tfsdk:"hwmon"`
	Textfile        types.Bool   `tfsdk:"textfile"`
	ThermalZone     types.Bool   `tfsdk:"thermal_zone"`
	Edac            types.Bool   `tfsdk:"edac"`
}

func (r *prometheusNodeExporterLuaConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prometheus_node_exporter_lua_config"
}

func (r *prometheusNodeExporterLuaConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromResourceConfigure(req, resp)
}

func (r *prometheusNodeExporterLuaConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Global prometheus-node-exporter-lua settings (uci prometheus-node-exporter-lua). This is a " +
			"singleton: it cannot be created or destroyed. `terraform destroy` only removes it from state; the " +
			"underlying settings are left as-is on the router.",
		Attributes: map[string]schema.Attribute{
			"id":               computedIDAttribute(),
			"managed":          managedAttribute(),
			"listen_ipv6":      optionalComputedBool("Listen on IPv6 as well as IPv4. Defaults to false."),
			"listen_interface": schema.StringAttribute{Optional: true, Description: "Network interface to bind the exporter to."},
			"listen_port":      schema.StringAttribute{Optional: true, Description: "TCP port the exporter listens on."},
			"cpu":              optionalComputedBool("Enable the cpu collector. Defaults to false."),
			"meminfo":          optionalComputedBool("Enable the meminfo collector. Defaults to false."),
			"netdev":           optionalComputedBool("Enable the netdev collector. Defaults to false."),
			"loadavg":          optionalComputedBool("Enable the loadavg collector. Defaults to false."),
			"filesystem":       optionalComputedBool("Enable the filesystem collector. Defaults to false."),
			"diskstats":        optionalComputedBool("Enable the diskstats collector. Defaults to false."),
			"uname":            optionalComputedBool("Enable the uname collector. Defaults to false."),
			"netstat":          optionalComputedBool("Enable the netstat collector. Defaults to false."),
			"stat":             optionalComputedBool("Enable the stat collector. Defaults to false."),
			"vmstat":           optionalComputedBool("Enable the vmstat collector. Defaults to false."),
			"boottime":         optionalComputedBool("Enable the boottime collector. Defaults to false."),
			"entropy":          optionalComputedBool("Enable the entropy collector. Defaults to false."),
			"time":             optionalComputedBool("Enable the time collector. Defaults to false."),
			"hwmon":            optionalComputedBool("Enable the hwmon collector. Defaults to false."),
			"textfile":         optionalComputedBool("Enable the textfile collector. Defaults to false."),
			"thermal_zone":     optionalComputedBool("Enable the thermal_zone collector. Defaults to false."),
			"edac":             optionalComputedBool("Enable the edac collector. Defaults to false."),
		},
	}
}

func (r *prometheusNodeExporterLuaConfigResource) body(_ context.Context, m prometheusNodeExporterLuaConfigModel) map[string]any {
	out := map[string]any{}
	putBool(out, "listen_ipv6", m.ListenIPv6)
	putStr(out, "listen_interface", m.ListenInterface)
	putStr(out, "listen_port", m.ListenPort)
	putBool(out, "cpu", m.CPU)
	putBool(out, "meminfo", m.Meminfo)
	putBool(out, "netdev", m.Netdev)
	putBool(out, "loadavg", m.Loadavg)
	putBool(out, "filesystem", m.Filesystem)
	putBool(out, "diskstats", m.Diskstats)
	putBool(out, "uname", m.Uname)
	putBool(out, "netstat", m.Netstat)
	putBool(out, "stat", m.Stat)
	putBool(out, "vmstat", m.Vmstat)
	putBool(out, "boottime", m.Boottime)
	putBool(out, "entropy", m.Entropy)
	putBool(out, "time", m.Time)
	putBool(out, "hwmon", m.Hwmon)
	putBool(out, "textfile", m.Textfile)
	putBool(out, "thermal_zone", m.ThermalZone)
	putBool(out, "edac", m.Edac)
	return out
}

func (r *prometheusNodeExporterLuaConfigResource) read(_ context.Context, obj map[string]any, m *prometheusNodeExporterLuaConfigModel) {
	m.ID = strVal(obj, "id")
	m.Managed = boolVal(obj, "managed")
	m.ListenIPv6 = boolVal(obj, "listen_ipv6")
	m.ListenInterface = strVal(obj, "listen_interface")
	m.ListenPort = strVal(obj, "listen_port")
	m.CPU = boolVal(obj, "cpu")
	m.Meminfo = boolVal(obj, "meminfo")
	m.Netdev = boolVal(obj, "netdev")
	m.Loadavg = boolVal(obj, "loadavg")
	m.Filesystem = boolVal(obj, "filesystem")
	m.Diskstats = boolVal(obj, "diskstats")
	m.Uname = boolVal(obj, "uname")
	m.Netstat = boolVal(obj, "netstat")
	m.Stat = boolVal(obj, "stat")
	m.Vmstat = boolVal(obj, "vmstat")
	m.Boottime = boolVal(obj, "boottime")
	m.Entropy = boolVal(obj, "entropy")
	m.Time = boolVal(obj, "time")
	m.Hwmon = boolVal(obj, "hwmon")
	m.Textfile = boolVal(obj, "textfile")
	m.ThermalZone = boolVal(obj, "thermal_zone")
	m.Edac = boolVal(obj, "edac")
}

func (r *prometheusNodeExporterLuaConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan prometheusNodeExporterLuaConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Patch(ctx, prometheusNodeExporterLuaConfigPath, r.body(ctx, plan))
	if err != nil {
		resp.Diagnostics.AddError("Error configuring prometheus-node-exporter-lua settings", err.Error())
		return
	}
	r.read(ctx, obj, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *prometheusNodeExporterLuaConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state prometheusNodeExporterLuaConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, found, err := r.client.GetObject(ctx, prometheusNodeExporterLuaConfigPath)
	if err != nil {
		resp.Diagnostics.AddError("Error reading prometheus-node-exporter-lua settings", err.Error())
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	r.read(ctx, obj, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *prometheusNodeExporterLuaConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan prometheusNodeExporterLuaConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, err := r.client.Patch(ctx, prometheusNodeExporterLuaConfigPath, r.body(ctx, plan))
	if err != nil {
		resp.Diagnostics.AddError("Error updating prometheus-node-exporter-lua settings", err.Error())
		return
	}
	r.read(ctx, obj, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete is a no-op: the prometheus-node-exporter-lua config singleton cannot be
// removed. State is dropped by the framework once this returns.
func (r *prometheusNodeExporterLuaConfigResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func (r *prometheusNodeExporterLuaConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
