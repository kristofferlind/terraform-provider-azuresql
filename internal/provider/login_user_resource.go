package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type loginUserResourceType struct{}

type loginUserResourceData struct {
	Id        types.String `tfsdk:"id"`
	LoginName types.String `tfsdk:"name"`
	Database  types.String `tfsdk:"database"`
}

type loginUserResource struct {
	provider provider
}

func (t loginUserResourceType) GetSchema(context context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "MSSQL Database user for login",

		Attributes: map[string]tfsdk.Attribute{
			// needed to keep testing framework happy, just set to login_name
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				MarkdownDescription: "Name of login to create user for",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"database": {
				MarkdownDescription: "Database to create user for the login in",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
		},
	}, nil
}

func (t loginUserResourceType) NewResource(context context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diagnostics := convertProviderType(in)

	return loginUserResource{
		provider: provider,
	}, diagnostics
}

func (resource loginUserResource) Create(context context.Context, request tfsdk.CreateResourceRequest, response *tfsdk.CreateResourceResponse) {
	var data loginUserResourceData

	diagnostics := request.Config.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	// create login
	err := resource.provider.manager.CreateLoginUser(context, data.LoginName.Value, data.Database.Value)
	if err != nil {
		response.Diagnostics.AddError("Failed to create login user", err.Error())
		return
	}

	data.Id = types.String{Value: data.LoginName.Value}

	diagnostics = response.State.Set(context, &data)
	response.Diagnostics.Append(diagnostics...)
}

func (resource loginUserResource) Read(context context.Context, request tfsdk.ReadResourceRequest, response *tfsdk.ReadResourceResponse) {
	var data loginUserResourceData

	diagnostics := request.State.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	user, err := resource.provider.manager.GetLoginUser(context, data.LoginName.Value, data.Database.Value)
	if err != nil {
		response.Diagnostics.AddError("Failed to read login user", err.Error())
		return
	}

	data.LoginName = types.String{Value: user.LoginName}
	data.Database = types.String{Value: user.Database}

	diagnostics = response.State.Set(context, &data)
	response.Diagnostics.Append(diagnostics...)
}

func (resource loginUserResource) Update(context context.Context, request tfsdk.UpdateResourceRequest, response *tfsdk.UpdateResourceResponse) {
	var data loginUserResourceData

	diagnostics := request.Plan.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.AddError("Not implemented", "update should not be allowed on this resource, should always be a replace operation")
}

func (resource loginUserResource) Delete(context context.Context, request tfsdk.DeleteResourceRequest, response *tfsdk.DeleteResourceResponse) {
	var data loginUserResourceData

	diagnostics := request.State.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	// delete login
	err := resource.provider.manager.DeleteLoginUser(context, data.LoginName.Value, data.Database.Value)
	if err != nil {
		response.Diagnostics.AddError("Failed to delete login user", err.Error())
		return
	}

	response.State.RemoveResource(context)
}

// import not supported
func (r loginUserResource) ImportState(context context.Context, request tfsdk.ImportResourceStateRequest, response *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStateNotImplemented(context, "", response)
}
