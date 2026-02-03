package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/domotz/terraform-provider-domotz/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &DeviceResource{}
	_ resource.ResourceWithImportState = &DeviceResource{}
)

// NewDeviceResource is a helper function to simplify the provider implementation
func NewDeviceResource() resource.Resource {
	return &DeviceResource{}
}

// DeviceResource defines the resource implementation
type DeviceResource struct {
	client *client.Client
}

// DeviceResourceModel describes the resource data model
type DeviceResourceModel struct {
	ID          types.String   `tfsdk:"id"`
	AgentID     types.Int64    `tfsdk:"agent_id"`
	DisplayName types.String   `tfsdk:"display_name"`
	IPAddresses types.List     `tfsdk:"ip_addresses"`
	UserData    *UserDataModel `tfsdk:"user_data"`
	Importance  types.String   `tfsdk:"importance"`
}

// UserDataModel describes the user_data nested object
type UserDataModel struct {
	Name   types.String `tfsdk:"name"`
	Model  types.String `tfsdk:"model"`
	Vendor types.String `tfsdk:"vendor"`
	Type   types.String `tfsdk:"type"`
}

// Metadata returns the resource type name
func (r *DeviceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

// Schema defines the schema for the resource
func (r *DeviceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an external IP device in Domotz.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Device ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"agent_id": schema.Int64Attribute{
				Description: "ID of the agent (collector) managing this device",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "Display name for the device",
				Required:    true,
			},
			"ip_addresses": schema.ListAttribute{
				Description: "List of IP addresses for the device",
				Required:    true,
				ElementType: types.StringType,
			},
			"importance": schema.StringAttribute{
				Description: "Device importance level (VITAL, FLOATING)",
				Optional:    true,
				Computed:    true,
			},
			"user_data": schema.SingleNestedAttribute{
				Description: "Custom metadata for the device",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "Custom name",
						Optional:    true,
					},
					"model": schema.StringAttribute{
						Description: "Device model",
						Optional:    true,
					},
					"vendor": schema.StringAttribute{
						Description: "Device vendor",
						Optional:    true,
					},
					"type": schema.StringAttribute{
						Description: "Device type",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *DeviceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = c
}

// Create creates the resource and sets the initial Terraform state
func (r *DeviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DeviceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract IP addresses from list
	var ipAddresses []string
	resp.Diagnostics.Append(plan.IPAddresses.ElementsAs(ctx, &ipAddresses, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build create request
	createReq := client.CreateDeviceRequest{
		DisplayName: plan.DisplayName.ValueString(),
		IPAddresses: ipAddresses,
	}

	if !plan.Importance.IsNull() {
		createReq.Importance = plan.Importance.ValueString()
	}

	if plan.UserData != nil {
		typeVal, _ := strconv.ParseInt(plan.UserData.Type.ValueString(), 10, 32)
		createReq.UserData = client.DeviceUserData{
			Name:   plan.UserData.Name.ValueString(),
			Model:  plan.UserData.Model.ValueString(),
			Vendor: plan.UserData.Vendor.ValueString(),
			Type:   int32(typeVal),
		}
	}

	// Create device
	device, err := r.client.CreateDevice(int32(plan.AgentID.ValueInt64()), createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating device",
			"Could not create device: "+err.Error(),
		)
		return
	}

	// Update state with created device
	plan.ID = types.StringValue(strconv.Itoa(int(device.ID)))
	if device.Importance != "" {
		plan.Importance = types.StringValue(device.Importance)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data
func (r *DeviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DeviceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID, err := strconv.ParseInt(state.ID.ValueString(), 10, 32)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing device ID",
			"Could not parse device ID: "+err.Error(),
		)
		return
	}

	// Read device from API
	device, err := r.client.GetDevice(int32(state.AgentID.ValueInt64()), int32(deviceID))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading device",
			"Could not read device: "+err.Error(),
		)
		return
	}

	// Update state
	state.DisplayName = types.StringValue(device.DisplayName)
	state.Importance = types.StringValue(device.Importance)

	// Update IP addresses
	ipAddressesList, diags := types.ListValueFrom(ctx, types.StringType, device.IPAddresses)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.IPAddresses = ipAddressesList

	// Update user data if present
	if device.UserData.Name != "" || device.UserData.Model != "" || device.UserData.Vendor != "" || device.UserData.Type != 0 {
		state.UserData = &UserDataModel{
			Name:   types.StringValue(device.UserData.Name),
			Model:  types.StringValue(device.UserData.Model),
			Vendor: types.StringValue(device.UserData.Vendor),
			Type:   types.StringValue(fmt.Sprintf("%d", device.UserData.Type)),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success
func (r *DeviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DeviceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID, err := strconv.ParseInt(plan.ID.ValueString(), 10, 32)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing device ID",
			"Could not parse device ID: "+err.Error(),
		)
		return
	}

	// Build update request
	updateReq := client.UpdateDeviceRequest{}

	displayName := plan.DisplayName.ValueString()
	updateReq.DisplayName = &displayName

	if !plan.Importance.IsNull() {
		importance := plan.Importance.ValueString()
		updateReq.Importance = &importance
	}

	if plan.UserData != nil {
		typeVal, _ := strconv.ParseInt(plan.UserData.Type.ValueString(), 10, 32)
		userData := client.DeviceUserData{
			Name:   plan.UserData.Name.ValueString(),
			Model:  plan.UserData.Model.ValueString(),
			Vendor: plan.UserData.Vendor.ValueString(),
			Type:   int32(typeVal),
		}
		updateReq.UserData = &userData
	}

	// Update device
	device, err := r.client.UpdateDevice(int32(plan.AgentID.ValueInt64()), int32(deviceID), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating device",
			"Could not update device: "+err.Error(),
		)
		return
	}

	// Update state
	plan.DisplayName = types.StringValue(device.DisplayName)
	if device.Importance != "" {
		plan.Importance = types.StringValue(device.Importance)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success
func (r *DeviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DeviceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID, err := strconv.ParseInt(state.ID.ValueString(), 10, 32)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing device ID",
			"Could not parse device ID: "+err.Error(),
		)
		return
	}

	// Delete device
	err = r.client.DeleteDevice(int32(state.AgentID.ValueInt64()), int32(deviceID))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting device",
			"Could not delete device: "+err.Error(),
		)
		return
	}
}

// ImportState imports the resource into Terraform state
func (r *DeviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: "agent_id:device_id"
	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format 'agent_id:device_id'",
		)
		return
	}

	agentID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid agent ID",
			"Could not parse agent ID: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("agent_id"), agentID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
