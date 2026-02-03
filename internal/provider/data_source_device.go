package provider

import (
	"context"
	"fmt"

	"github.com/domotz/terraform-provider-domotz/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DeviceDataSource{}

func NewDeviceDataSource() datasource.DataSource {
	return &DeviceDataSource{}
}

type DeviceDataSource struct {
	client *client.Client
}

type DeviceDataSourceModel struct {
	AgentID     types.Int64    `tfsdk:"agent_id"`
	ID          types.Int64    `tfsdk:"id"`
	DisplayName types.String   `tfsdk:"display_name"`
	Protocol    types.String   `tfsdk:"protocol"`
	IPAddresses types.List     `tfsdk:"ip_addresses"`
	Importance  types.String   `tfsdk:"importance"`
	UserData    *UserDataModel `tfsdk:"user_data"`
}

func (d *DeviceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

func (d *DeviceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves details of a specific device.",
		Attributes: map[string]schema.Attribute{
			"agent_id": schema.Int64Attribute{
				Description: "ID of the collector managing the device",
				Required:    true,
			},
			"id": schema.Int64Attribute{
				Description: "Device ID",
				Required:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "Device display name",
				Computed:    true,
			},
			"protocol": schema.StringAttribute{
				Description: "Device protocol (IP, DUMMY, etc.)",
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
			"user_data": schema.SingleNestedAttribute{
				Description: "Custom metadata for the device",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "Custom name",
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
					"type": schema.StringAttribute{
						Description: "Device type",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *DeviceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DeviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DeviceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	device, err := d.client.GetDevice(
		int32(config.AgentID.ValueInt64()),
		int32(config.ID.ValueInt64()),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error reading device", err.Error())
		return
	}

	config.DisplayName = types.StringValue(device.DisplayName)
	config.Protocol = types.StringValue(device.Protocol)
	config.Importance = types.StringValue(device.Importance)

	ipAddressesList, diags := types.ListValueFrom(ctx, types.StringType, device.IPAddresses)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.IPAddresses = ipAddressesList

	if device.UserData.Name != "" || device.UserData.Model != "" || device.UserData.Vendor != "" || device.UserData.Type != 0 {
		config.UserData = &UserDataModel{
			Name:   types.StringValue(device.UserData.Name),
			Model:  types.StringValue(device.UserData.Model),
			Vendor: types.StringValue(device.UserData.Vendor),
			Type:   types.StringValue(fmt.Sprintf("%d", device.UserData.Type)),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
