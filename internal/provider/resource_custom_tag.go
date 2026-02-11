package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/domotz/terraform-provider-domotz/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &CustomTagResource{}
	_ resource.ResourceWithImportState = &CustomTagResource{}
)

func NewCustomTagResource() resource.Resource {
	return &CustomTagResource{}
}

type CustomTagResource struct {
	client *client.Client
}

type CustomTagResourceModel struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Colour types.String `tfsdk:"colour"`
}

func (r *CustomTagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_tag"
}

func (r *CustomTagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a custom tag in Domotz.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Tag ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Tag name",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"colour": schema.StringAttribute{
				Description: "Tag color in hex format (e.g., #FF5733)",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`),
						"must be a valid hex color (e.g., #FF5733)",
					),
				},
			},
		},
	}
}

func (r *CustomTagResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CustomTagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CustomTagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := client.CreateTagRequest{
		Name:   plan.Name.ValueString(),
		Colour: plan.Colour.ValueString(),
	}

	tag, err := r.client.CreateTag(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating tag", err.Error())
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(int(tag.ID)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CustomTagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CustomTagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagID, err := strconv.ParseInt(state.ID.ValueString(), 10, 32)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing tag ID", err.Error())
		return
	}

	tag, err := r.client.GetTag(ctx, int32(tagID))
	if err != nil {
		var notFound *client.NotFoundError
		if errors.As(err, &notFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading tag", err.Error())
		return
	}

	state.Name = types.StringValue(tag.Name)
	state.Colour = types.StringValue(tag.Colour)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CustomTagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CustomTagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagID, err := strconv.ParseInt(plan.ID.ValueString(), 10, 32)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing tag ID", err.Error())
		return
	}

	name := plan.Name.ValueString()
	colour := plan.Colour.ValueString()
	updateReq := client.UpdateTagRequest{
		Name:   &name,
		Colour: &colour,
	}

	tag, err := r.client.UpdateTag(ctx, int32(tagID), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tag", err.Error())
		return
	}

	plan.Name = types.StringValue(tag.Name)
	plan.Colour = types.StringValue(tag.Colour)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CustomTagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CustomTagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagID, err := strconv.ParseInt(state.ID.ValueString(), 10, 32)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing tag ID", err.Error())
		return
	}

	err = r.client.DeleteTag(ctx, int32(tagID))
	if err != nil {
		var notFound *client.NotFoundError
		if errors.As(err, &notFound) {
			return
		}
		resp.Diagnostics.AddError("Error deleting tag", err.Error())
		return
	}
}

func (r *CustomTagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
