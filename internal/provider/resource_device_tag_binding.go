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

var (
	_ resource.Resource                = &DeviceTagBindingResource{}
	_ resource.ResourceWithImportState = &DeviceTagBindingResource{}
)

func NewDeviceTagBindingResource() resource.Resource {
	return &DeviceTagBindingResource{}
}

type DeviceTagBindingResource struct {
	client *client.Client
}

type DeviceTagBindingResourceModel struct {
	ID       types.String `tfsdk:"id"`
	AgentID  types.Int64  `tfsdk:"agent_id"`
	DeviceID types.Int64  `tfsdk:"device_id"`
	TagID    types.Int64  `tfsdk:"tag_id"`
}

func (r *DeviceTagBindingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_tag_binding"
}

func (r *DeviceTagBindingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the binding between a device and a tag in Domotz.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Binding ID (format: agent_id:device_id:tag_id)",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"agent_id": schema.Int64Attribute{
				Description: "ID of the collector managing the device",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"device_id": schema.Int64Attribute{
				Description: "ID of the device",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"tag_id": schema.Int64Attribute{
				Description: "ID of the tag",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *DeviceTagBindingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *DeviceTagBindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DeviceTagBindingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	agentID := int32(plan.AgentID.ValueInt64())
	deviceID := int32(plan.DeviceID.ValueInt64())
	tagID := int32(plan.TagID.ValueInt64())

	err := r.client.BindTagToDevice(agentID, deviceID, tagID)
	if err != nil {
		resp.Diagnostics.AddError("Error binding tag to device", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%d:%d:%d", agentID, deviceID, tagID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DeviceTagBindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DeviceTagBindingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	agentID := int32(state.AgentID.ValueInt64())
	deviceID := int32(state.DeviceID.ValueInt64())
	tagID := int32(state.TagID.ValueInt64())

	// Verify the binding still exists by listing device tags
	tags, err := r.client.ListDeviceTags(agentID, deviceID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading device tags", err.Error())
		return
	}

	// Check if our tag is in the list
	found := false
	for _, tag := range tags {
		if tag.ID == tagID {
			found = true
			break
		}
	}

	if !found {
		// Tag binding no longer exists, remove from state
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DeviceTagBindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Update is not supported - all changes require replacement
	resp.Diagnostics.AddError(
		"Update not supported",
		"Device tag bindings cannot be updated. All changes require replacement.",
	)
}

func (r *DeviceTagBindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DeviceTagBindingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	agentID := int32(state.AgentID.ValueInt64())
	deviceID := int32(state.DeviceID.ValueInt64())
	tagID := int32(state.TagID.ValueInt64())

	err := r.client.UnbindTagFromDevice(agentID, deviceID, tagID)
	if err != nil {
		resp.Diagnostics.AddError("Error unbinding tag from device", err.Error())
		return
	}
}

func (r *DeviceTagBindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: "agent_id:device_id:tag_id"
	parts := strings.Split(req.ID, ":")
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format 'agent_id:device_id:tag_id'",
		)
		return
	}

	agentID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid agent ID", err.Error())
		return
	}

	deviceID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid device ID", err.Error())
		return
	}

	tagID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid tag ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("agent_id"), agentID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("device_id"), deviceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tag_id"), tagID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
