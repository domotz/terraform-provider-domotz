package provider

import (
	"context"
	"fmt"

	"github.com/domotz/terraform-provider-domotz/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DeviceVariablesDataSource{}

func NewDeviceVariablesDataSource() datasource.DataSource {
	return &DeviceVariablesDataSource{}
}

type DeviceVariablesDataSource struct {
	client *client.Client
}

type DeviceVariablesDataSourceModel struct {
	AgentID   types.Int64     `tfsdk:"agent_id"`
	DeviceID  types.Int64     `tfsdk:"device_id"`
	Variables []VariableModel `tfsdk:"variables"`
}

type VariableModel struct {
	ID            types.Int64  `tfsdk:"id"`
	Label         types.String `tfsdk:"label"`
	Path          types.String `tfsdk:"path"`
	Value         types.String `tfsdk:"value"`
	Unit          types.String `tfsdk:"unit"`
	PreviousValue types.String `tfsdk:"previous_value"`
	Metric        types.String `tfsdk:"metric"`
}

func (d *DeviceVariablesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_variables"
}

func (d *DeviceVariablesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves variables (metrics) for a specific device.",
		Attributes: map[string]schema.Attribute{
			"agent_id": schema.Int64Attribute{
				Description: "ID of the collector managing the device",
				Required:    true,
			},
			"device_id": schema.Int64Attribute{
				Description: "Device ID",
				Required:    true,
			},
			"variables": schema.ListNestedAttribute{
				Description: "List of device variables/metrics",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Variable ID",
							Computed:    true,
						},
						"label": schema.StringAttribute{
							Description: "Variable label",
							Computed:    true,
						},
						"path": schema.StringAttribute{
							Description: "Variable path",
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "Current value",
							Computed:    true,
						},
						"unit": schema.StringAttribute{
							Description: "Unit of measurement",
							Computed:    true,
						},
						"previous_value": schema.StringAttribute{
							Description: "Previous value",
							Computed:    true,
						},
						"metric": schema.StringAttribute{
							Description: "Metric type",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *DeviceVariablesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DeviceVariablesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DeviceVariablesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	variables, err := d.client.ListVariables(
		int32(config.AgentID.ValueInt64()),
		int32(config.DeviceID.ValueInt64()),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error listing device variables", err.Error())
		return
	}

	config.Variables = make([]VariableModel, 0, len(variables))
	for _, v := range variables {
		config.Variables = append(config.Variables, VariableModel{
			ID:            types.Int64Value(int64(v.ID)),
			Label:         types.StringValue(v.Label),
			Path:          types.StringValue(v.Path),
			Value:         types.StringValue(v.Value),
			Unit:          types.StringValue(v.Unit),
			PreviousValue: types.StringValue(v.PreviousValue),
			Metric:        types.StringValue(v.Metric),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
