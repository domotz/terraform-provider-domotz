package provider

import (
	"context"
	"fmt"

	"github.com/domotz/terraform-provider-domotz/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DevicesDataSource{}

func NewDevicesDataSource() datasource.DataSource {
	return &DevicesDataSource{}
}

type DevicesDataSource struct {
	client *client.Client
}

type DevicesDataSourceModel struct {
	AgentID types.Int64         `tfsdk:"agent_id"`
	Devices []DeviceListModel   `tfsdk:"devices"`
}

type DeviceListModel struct {
	ID          types.Int64  `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Protocol    types.String `tfsdk:"protocol"`
	IPAddresses types.List   `tfsdk:"ip_addresses"`
	Importance  types.String `tfsdk:"importance"`
	Vendor      types.String `tfsdk:"vendor"`       // Auto-discovered vendor
	Model       types.String `tfsdk:"model"`        // Auto-discovered model
	UserData    types.Object `tfsdk:"user_data"`    // User-editable metadata
}

type DeviceUserDataModel struct {
	Name   types.String `tfsdk:"name"`
	Model  types.String `tfsdk:"model"`
	Vendor types.String `tfsdk:"vendor"`
	Type   types.Int64  `tfsdk:"type"`
}

func (d *DevicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices"
}

func (d *DevicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a list of devices for a specific collector.",
		Attributes: map[string]schema.Attribute{
			"agent_id": schema.Int64Attribute{
				Description: "ID of the collector managing the devices",
				Required:    true,
			},
			"devices": schema.ListNestedAttribute{
				Description: "List of devices",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Device ID",
							Computed:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "Device display name",
							Computed:    true,
						},
						"protocol": schema.StringAttribute{
							Description: "Device protocol",
							Computed:    true,
						},
						"ip_addresses": schema.ListAttribute{
							Description: "List of IP addresses",
							Computed:    true,
							ElementType: types.StringType,
						},
						"importance": schema.StringAttribute{
							Description: "Device importance level",
							Computed:    true,
						},
						"vendor": schema.StringAttribute{
							Description: "Auto-discovered device vendor (e.g., Ubiquiti Inc, Apple, etc.)",
							Computed:    true,
						},
						"model": schema.StringAttribute{
							Description: "Auto-discovered device model",
							Computed:    true,
						},
						"user_data": schema.SingleNestedAttribute{
							Description: "Custom metadata for the device",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "Device name",
									Computed:    true,
								},
								"model": schema.StringAttribute{
									Description: "Device model",
									Computed:    true,
								},
								"vendor": schema.StringAttribute{
									Description: "Device vendor",
									Computed:    true,
								},
								"type": schema.Int64Attribute{
									Description: "Device type ID",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *DevicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DevicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DevicesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	devices, err := d.client.ListDevices(int32(config.AgentID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error listing devices", err.Error())
		return
	}

	config.Devices = make([]DeviceListModel, 0, len(devices))
	for _, device := range devices {
		ipAddressesList, diags := types.ListValueFrom(ctx, types.StringType, device.IPAddresses)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Build user_data object
		userDataAttrs := map[string]attr.Type{
			"name":   types.StringType,
			"model":  types.StringType,
			"vendor": types.StringType,
			"type":   types.Int64Type,
		}

		userDataValues := map[string]attr.Value{
			"name":   types.StringValue(device.UserData.Name),
			"model":  types.StringValue(device.UserData.Model),
			"vendor": types.StringValue(device.UserData.Vendor),
			"type":   types.Int64Value(int64(device.UserData.Type)),
		}

		userDataObj, diags := types.ObjectValue(userDataAttrs, userDataValues)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		config.Devices = append(config.Devices, DeviceListModel{
			ID:          types.Int64Value(int64(device.ID)),
			DisplayName: types.StringValue(device.DisplayName),
			Protocol:    types.StringValue(device.Protocol),
			IPAddresses: ipAddressesList,
			Importance:  types.StringValue(device.Importance),
			Vendor:      types.StringValue(device.Vendor),
			Model:       types.StringValue(device.Model),
			UserData:    userDataObj,
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
