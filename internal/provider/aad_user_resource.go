package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type externalUserResourceType struct{}

type externalUserResourceData struct {
	Id        types.String `tfsdk:"id"`
	LoginName types.String `tfsdk:"name"`
	Database  types.String `tfsdk:"database"`
}

type externalUserResource struct {
	provider provider
}

func (t externalUserResourceType) GetSchema(context context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Database user for AzureAD user/group/application, needed to configure access at database level",

		Attributes: map[string]tfsdk.Attribute{
			// needed to keep testing framework happy, just set to name
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				MarkdownDescription: "Display name of AzureAD user/group/application",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"database": {
				MarkdownDescription: "Database to create user in",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
		},
	}, nil
}

func (t externalUserResourceType) NewResource(context context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diagnostics := convertProviderType(in)

	return externalUserResource{
		provider: provider,
	}, diagnostics
}

func (resource externalUserResource) Create(context context.Context, request tfsdk.CreateResourceRequest, response *tfsdk.CreateResourceResponse) {
	var data externalUserResourceData

	diagnostics := request.Config.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	err := resource.provider.manager.CreateExternalUser(context, data.LoginName.Value, data.Database.Value)
	if err != nil {
		response.Diagnostics.AddError("Failed to create user", err.Error())
		return
	}

	data.Id = types.String{Value: data.LoginName.Value}

	diagnostics = response.State.Set(context, &data)
	response.Diagnostics.Append(diagnostics...)
}

func (resource externalUserResource) Read(context context.Context, request tfsdk.ReadResourceRequest, response *tfsdk.ReadResourceResponse) {
	var data externalUserResourceData

	diagnostics := request.State.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	user, err := resource.provider.manager.GetLoginUser(context, data.LoginName.Value, data.Database.Value)
	if err != nil {
		response.Diagnostics.AddError("Failed to read user", err.Error())
		return
	}

	data.LoginName = types.String{Value: user.LoginName}
	data.Database = types.String{Value: user.Database}

	diagnostics = response.State.Set(context, &data)
	response.Diagnostics.Append(diagnostics...)
}

func (resource externalUserResource) Update(context context.Context, request tfsdk.UpdateResourceRequest, response *tfsdk.UpdateResourceResponse) {
	var data externalUserResourceData

	diagnostics := request.Plan.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.AddError("Not implemented", "update should not be allowed on this resource, should always be a replace operation")
}

func (resource externalUserResource) Delete(context context.Context, request tfsdk.DeleteResourceRequest, response *tfsdk.DeleteResourceResponse) {
	var data externalUserResourceData

	diagnostics := request.State.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	err := resource.provider.manager.DeleteLoginUser(context, data.LoginName.Value, data.Database.Value)
	if err != nil {
		response.Diagnostics.AddError("Failed to delete user", err.Error())
		return
	}

	response.State.RemoveResource(context)
}

// import not supported
func (r externalUserResource) ImportState(context context.Context, request tfsdk.ImportResourceStateRequest, response *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStateNotImplemented(context, "", response)
}
