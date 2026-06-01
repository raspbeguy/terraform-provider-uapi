package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &lldpdConfigDataSource{}
	_ datasource.DataSourceWithConfigure = &lldpdConfigDataSource{}
)

type lldpdConfigDataSource struct{ client *client.Client }

func NewLldpdConfigDataSource() datasource.DataSource { return &lldpdConfigDataSource{} }

func (d *lldpdConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lldpd_config"
}

func (d *lldpdConfigDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *lldpdConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "The global lldpd daemon settings (uci lldpd.lldpd).",
		Attributes: map[string]dsschema.Attribute{
			"id":                dsComputedString("Stable id of the lldpd config section."),
			"managed":           dsManagedAttribute(),
			"etag":              dsComputedString("Opaque ETag of the resource's current state."),
			"enable_cdp":        dsComputedBool("Whether Cisco Discovery Protocol advertisement is enabled."),
			"enable_fdp":        dsComputedBool("Whether Foundry Discovery Protocol advertisement is enabled."),
			"enable_sonmp":      dsComputedBool("Whether Nortel SONMP advertisement is enabled."),
			"enable_edp":        dsComputedBool("Whether Extreme Discovery Protocol advertisement is enabled."),
			"enable_lldpmed":    dsComputedBool("Whether LLDP-MED advertisement is enabled."),
			"lldp_class":        dsComputedString("LLDP-MED device class (1-4)."),
			"lldp_description":  dsComputedBool("Whether the system description is advertised."),
			"lldp_capabilities": dsComputedBool("Whether system capabilities are advertised."),
			"lldp_mgmt_ip":      dsComputedString("Management IP address advertised."),
			"interface":         dsComputedStringList("Network interfaces lldpd listens on. Empty means all interfaces."),
		},
	}
}

func (d *lldpdConfigDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	obj, etag, found, err := d.client.GetObject(ctx, lldpdConfigPath)
	if err != nil {
		resp.Diagnostics.AddError("Error reading lldpd settings", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("Lldpd settings not found", "The lldpd config singleton is missing on the router")
		return
	}
	var m lldpdConfigModel
	ds := newDiagsink(&resp.Diagnostics)
	(&lldpdConfigResource{}).read(ctx, obj, &m, ds)
	m.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
