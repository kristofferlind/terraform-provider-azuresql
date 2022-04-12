package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type externalLoginResourceType struct{}

type externalLoginResourceData struct {
	Id        types.String `tfsdk:"id"`
	LoginName types.String `tfsdk:"name"`
}

type externalLoginResource struct {
	provider provider
}

func (t externalLoginResourceType) GetSchema(context context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "MSSQL Server login for Azure AD user, group or application",

		Attributes: map[string]tfsdk.Attribute{
			// needed to keep testing framework happy, just set to login_name
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				MarkdownDescription: "Supplied as name of login, should match display_name of AzureAD user or access group",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
		},
	}, nil
}

func (t externalLoginResourceType) NewResource(context context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diagnostics := convertProviderType(in)

	return externalLoginResource{
		provider: provider,
	}, diagnostics
}

func (resource externalLoginResource) Create(context context.Context, request tfsdk.CreateResourceRequest, response *tfsdk.CreateResourceResponse) {
	var data externalLoginResourceData

	diagnostics := request.Config.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	// create login
	err := resource.provider.manager.CreateAADLogin(context, data.LoginName.Value)
	if err != nil {
		response.Diagnostics.AddError("Failed to create login", err.Error())
		return
	}

	data.Id = types.String{Value: data.LoginName.Value}

	tflog.Trace(context, "created a resource")

	diagnostics = response.State.Set(context, &data)
	response.Diagnostics.Append(diagnostics...)
}

func (resource externalLoginResource) Read(context context.Context, request tfsdk.ReadResourceRequest, response *tfsdk.ReadResourceResponse) {
	var data externalLoginResourceData

	diagnostics := request.State.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	// read login
	login, err := resource.provider.manager.GetLogin(context, data.LoginName.Value)
	if err != nil {
		response.Diagnostics.AddError("Failed to read login", err.Error())
		return
	}

	data.LoginName = types.String{Value: login.LoginName}

	diagnostics = response.State.Set(context, &data)
	response.Diagnostics.Append(diagnostics...)
}

func (resource externalLoginResource) Update(context context.Context, request tfsdk.UpdateResourceRequest, response *tfsdk.UpdateResourceResponse) {
	var data externalLoginResourceData

	diagnostics := request.Plan.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.AddError("Not implemented", "update should not be allowed on this resource, should always be a replace operation")
}

func (resource externalLoginResource) Delete(context context.Context, request tfsdk.DeleteResourceRequest, response *tfsdk.DeleteResourceResponse) {
	var data externalLoginResourceData

	diagnostics := request.State.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	// delete login
	err := resource.provider.manager.DeleteLogin(context, data.LoginName.Value)
	if err != nil {
		response.Diagnostics.AddError("Failed to delete login", err.Error())
		return
	}

	response.State.RemoveResource(context)
}

// import not supported
func (r externalLoginResource) ImportState(context context.Context, request tfsdk.ImportResourceStateRequest, response *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStateNotImplemented(context, "", response)
}
