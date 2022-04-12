package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type userRoleResourceType struct{}

type userRoleResourceData struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Database types.String `tfsdk:"database"`
	Role     types.String `tfsdk:"role"`
}

type userRoleResource struct {
	provider provider
}

func (t userRoleResourceType) GetSchema(context context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Set database level role membership for user",

		Attributes: map[string]tfsdk.Attribute{
			// needed to keep testing framework happy, just set to name
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				MarkdownDescription: "Name of user to set role membership for",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"database": {
				MarkdownDescription: "Which database to set the role membership in",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"role": {
				MarkdownDescription: "Role to add user as member of",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
		},
	}, nil
}

func (t userRoleResourceType) NewResource(context context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diagnostics := convertProviderType(in)

	return userRoleResource{
		provider: provider,
	}, diagnostics
}

func (resource userRoleResource) Create(context context.Context, request tfsdk.CreateResourceRequest, response *tfsdk.CreateResourceResponse) {
	var data userRoleResourceData

	diagnostics := request.Config.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	// create login
	err := resource.provider.manager.AddRole(context, data.Name.Value, data.Role.Value, data.Database.Value)
	if err != nil {
		response.Diagnostics.AddError("Failed to add user to role members", err.Error())
		return
	}

	data.Id = types.String{Value: data.Name.Value}

	diagnostics = response.State.Set(context, &data)
	response.Diagnostics.Append(diagnostics...)
}

func (resource userRoleResource) Read(context context.Context, request tfsdk.ReadResourceRequest, response *tfsdk.ReadResourceResponse) {
	var data userRoleResourceData

	diagnostics := request.State.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	user, err := resource.provider.manager.GetUserWithRole(context, data.Name.Value, data.Role.Value, data.Database.Value)
	if err != nil {
		response.Diagnostics.AddError("Failed to read user role membership", err.Error())
		return
	}

	data.Name = types.String{Value: user.Name}
	data.Database = types.String{Value: user.Database}
	data.Role = types.String{Value: user.Role}

	diagnostics = response.State.Set(context, &data)
	response.Diagnostics.Append(diagnostics...)
}

func (resource userRoleResource) Update(context context.Context, request tfsdk.UpdateResourceRequest, response *tfsdk.UpdateResourceResponse) {
	var data userRoleResourceData

	diagnostics := request.Plan.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.AddError("Not implemented", "update should not be allowed on this resource, should always be a replace operation")
}

func (resource userRoleResource) Delete(context context.Context, request tfsdk.DeleteResourceRequest, response *tfsdk.DeleteResourceResponse) {
	var data userRoleResourceData

	diagnostics := request.State.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	// delete login
	err := resource.provider.manager.RemoveRole(context, data.Name.Value, data.Role.Value, data.Database.Value)
	if err != nil {
		response.Diagnostics.AddError("Failed to remove user from role members", err.Error())
		return
	}

	response.State.RemoveResource(context)
}

// import not supported
func (r userRoleResource) ImportState(context context.Context, request tfsdk.ImportResourceStateRequest, response *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStateNotImplemented(context, "", response)
}
