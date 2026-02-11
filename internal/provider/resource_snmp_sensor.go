package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/domotz/terraform-provider-domotz/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &SNMPSensorResource{}
	_ resource.ResourceWithImportState = &SNMPSensorResource{}
)

func NewSNMPSensorResource() resource.Resource {
	return &SNMPSensorResource{}
}

type SNMPSensorResource struct {
	client *client.Client
}

type SNMPSensorResourceModel struct {
	ID        types.String `tfsdk:"id"`
	AgentID   types.Int64  `tfsdk:"agent_id"`
	DeviceID  types.Int64  `tfsdk:"device_id"`
	Name      types.String `tfsdk:"name"`
	OID       types.String `tfsdk:"oid"`
	Category  types.String `tfsdk:"category"`
	ValueType types.String `tfsdk:"value_type"`
}

func (r *SNMPSensorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snmp_sensor"
}

func (r *SNMPSensorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an SNMP sensor in Domotz.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Sensor ID",
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
			"name": schema.StringAttribute{
				Description: "Sensor name",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"oid": schema.StringAttribute{
				Description: "SNMP OID to monitor",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"category": schema.StringAttribute{
				Description: "Sensor category",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("OTHER", "NETWORKING", "COMPUTING", "STORAGE"),
				},
			},
			"value_type": schema.StringAttribute{
				Description: "Value type (STRING, NUMERIC)",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("STRING", "NUMERIC"),
				},
			},
		},
	}
}

func (r *SNMPSensorResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SNMPSensorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SNMPSensorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := client.CreateSNMPSensorRequest{
		Name:      plan.Name.ValueString(),
		OID:       plan.OID.ValueString(),
		Category:  plan.Category.ValueString(),
		ValueType: plan.ValueType.ValueString(),
	}

	sensor, err := r.client.CreateSNMPSensor(
		ctx,
		int32(plan.AgentID.ValueInt64()),
		int32(plan.DeviceID.ValueInt64()),
		createReq,
	)
	if err != nil {
		resp.Diagnostics.AddError("Error creating SNMP sensor", err.Error())
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(int(sensor.ID)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SNMPSensorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SNMPSensorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sensorID, err := strconv.ParseInt(state.ID.ValueString(), 10, 32)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing sensor ID", err.Error())
		return
	}

	sensor, err := r.client.GetSNMPSensor(
		ctx,
		int32(state.AgentID.ValueInt64()),
		int32(state.DeviceID.ValueInt64()),
		int32(sensorID),
	)
	if err != nil {
		var notFound *client.NotFoundError
		if errors.As(err, &notFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading SNMP sensor", err.Error())
		return
	}

	state.Name = types.StringValue(sensor.Name)
	state.OID = types.StringValue(sensor.OID)
	state.Category = types.StringValue(sensor.Category)
	state.ValueType = types.StringValue(sensor.ValueType)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SNMPSensorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// SNMP sensors cannot be updated - all changes require replacement
	resp.Diagnostics.AddError(
		"Update not supported",
		"SNMP sensors cannot be updated. All changes require replacement.",
	)
}

func (r *SNMPSensorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SNMPSensorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sensorID, err := strconv.ParseInt(state.ID.ValueString(), 10, 32)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing sensor ID", err.Error())
		return
	}

	err = r.client.DeleteSNMPSensor(
		ctx,
		int32(state.AgentID.ValueInt64()),
		int32(state.DeviceID.ValueInt64()),
		int32(sensorID),
	)
	if err != nil {
		var notFound *client.NotFoundError
		if errors.As(err, &notFound) {
			return
		}
		resp.Diagnostics.AddError("Error deleting SNMP sensor", err.Error())
		return
	}
}

func (r *SNMPSensorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: "agent_id:device_id:sensor_id"
	parts := strings.Split(req.ID, ":")
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format 'agent_id:device_id:sensor_id'",
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

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("agent_id"), agentID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("device_id"), deviceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[2])...)
}
