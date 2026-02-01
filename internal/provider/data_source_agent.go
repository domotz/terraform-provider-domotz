package provider

import (
	"context"
	"fmt"

	"github.com/domotz/terraform-provider-domotz/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &AgentDataSource{}

func NewAgentDataSource() datasource.DataSource {
	return &AgentDataSource{}
}

type AgentDataSource struct {
	client *client.Client
}

type AgentDataSourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Status      types.String `tfsdk:"status"`
	TeamID      types.Int64  `tfsdk:"team_id"`
	TeamName    types.String `tfsdk:"team_name"`
}

func (d *AgentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent"
}

func (d *AgentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves details of a specific Domotz agent (collector).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Collector ID",
				Required:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "Collector display name",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Collector status (ONLINE, OFFLINE)",
				Computed:    true,
			},
			"team_id": schema.Int64Attribute{
				Description: "Team/area ID",
				Computed:    true,
			},
			"team_name": schema.StringAttribute{
				Description: "Team/area name",
				Computed:    true,
			},
		},
	}
}

func (d *AgentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = c
}

func (d *AgentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config AgentDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	agent, err := d.client.GetAgent(int32(config.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading agent", err.Error())
		return
	}

	config.DisplayName = types.StringValue(agent.DisplayName)
	config.Status = types.StringValue(agent.Status.Value)
	config.TeamID = types.Int64Value(int64(agent.Team.ID))
	config.TeamName = types.StringValue(agent.Team.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
