package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raspbeguy/terraform-provider-uapi/internal/client"
)

var (
	_ datasource.DataSource              = &sqmQueueDataSource{}
	_ datasource.DataSourceWithConfigure = &sqmQueueDataSource{}
)

type sqmQueueDataSource struct{ client *client.Client }

func NewSqmQueueDataSource() datasource.DataSource { return &sqmQueueDataSource{} }

func (d *sqmQueueDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sqm_queue"
}

func (d *sqmQueueDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromDataSourceConfigure(req, resp)
}

func (d *sqmQueueDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Look up an SQM queue by id.",
		Attributes: map[string]dsschema.Attribute{
			"id":        dsIDAttribute(),
			"managed":   dsManagedAttribute(),
			"etag":      dsComputedString("Opaque ETag of the resource's current state."),
			"enabled":   dsComputedBool("Whether the queue is enabled."),
			"interface": dsComputedString("Network interface the queue is attached to."),
			"download":  dsComputedString("Download rate limit in kbit/s."),
			"upload":    dsComputedString("Upload rate limit in kbit/s."),
			"qdisc":     dsComputedString("Queueing discipline: cake, fq_codel, pie, or htb."),
			"script":    dsComputedString("SQM script: piece_of_cake.qos, simple.qos, simplest.qos, or layer_cake.qos."),
			"linklayer": dsComputedString("Link layer adaptation: none, ethernet, or atm."),
			"overhead":  dsComputedString("Per-packet overhead in bytes for link layer adaptation."),
		},
	}
}

func (d *sqmQueueDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var m sqmQueueModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, etag, found, err := d.client.GetObject(ctx, "/"+sqmQueueCollection+"/"+m.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading SQM queue", err.Error())
		return
	}
	if !found {
		resp.Diagnostics.AddError("SQM queue not found", "No SQM queue with id "+m.ID.ValueString())
		return
	}
	(&sqmQueueResource{}).read(ctx, obj, &m)
	m.ETag = types.StringValue(etag)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
