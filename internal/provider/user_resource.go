package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type userResourceType struct{}

type userResourceData struct {
	Id        types.String `tfsdk:"id"`
	LoginName types.String `tfsdk:"name"`
	Password  types.String `tfsdk:"password"`
	Database  types.String `tfsdk:"database"`
}

type userResource struct {
	provider provider
}

func (t userResourceType) GetSchema(context context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "MSSQL Database user for login",

		Attributes: map[string]tfsdk.Attribute{
			// needed to keep testing framework happy, just set to name
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				MarkdownDescription: "Name of user",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"password": {
				MarkdownDescription: "Password part of user credential",
				Required:            true,
				Type:                types.StringType,
				Sensitive:           true,
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

func (t userResourceType) NewResource(context context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diagnostics := convertProviderType(in)

	return userResource{
		provider: provider,
	}, diagnostics
}

func (resource userResource) Create(context context.Context, request tfsdk.CreateResourceRequest, response *tfsdk.CreateResourceResponse) {
	var data userResourceData

	diagnostics := request.Config.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	err := resource.provider.manager.CreateUser(context, data.LoginName.Value, data.Password.Value, data.Database.Value)
	if err != nil {
		response.Diagnostics.AddError("Failed to create user", err.Error())
		return
	}

	data.Id = types.String{Value: data.LoginName.Value}

	diagnostics = response.State.Set(context, &data)
	response.Diagnostics.Append(diagnostics...)
}

func (resource userResource) Read(context context.Context, request tfsdk.ReadResourceRequest, response *tfsdk.ReadResourceResponse) {
	var data userResourceData

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

func (resource userResource) Update(context context.Context, request tfsdk.UpdateResourceRequest, response *tfsdk.UpdateResourceResponse) {
	var data userResourceData

	diagnostics := request.Plan.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.AddError("Not implemented", "update should not be allowed on this resource, should always be a replace operation")
}

func (resource userResource) Delete(context context.Context, request tfsdk.DeleteResourceRequest, response *tfsdk.DeleteResourceResponse) {
	var data userResourceData

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
func (r userResource) ImportState(context context.Context, request tfsdk.ImportResourceStateRequest, response *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStateNotImplemented(context, "", response)
}
