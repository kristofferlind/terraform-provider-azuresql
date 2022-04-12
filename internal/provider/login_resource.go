package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type loginResourceType struct{}

type loginResourceData struct {
	Id        types.String `tfsdk:"id"`
	LoginName types.String `tfsdk:"name"`
	Password  types.String `tfsdk:"password"`
}

type loginResource struct {
	provider provider
}

func (t loginResourceType) GetSchema(context context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "MSSQL Server login",

		Attributes: map[string]tfsdk.Attribute{
			// needed to keep testing framework happy, just set to login_name
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				MarkdownDescription: "Name part of login credential",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"password": {
				MarkdownDescription: "Password part of login credential",
				Required:            true,
				Type:                types.StringType,
				Sensitive:           true,
			},
		},
	}, nil
}

func (t loginResourceType) NewResource(context context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diagnostics := convertProviderType(in)

	return loginResource{
		provider: provider,
	}, diagnostics
}

func (resource loginResource) Create(context context.Context, request tfsdk.CreateResourceRequest, response *tfsdk.CreateResourceResponse) {
	var data loginResourceData

	diagnostics := request.Config.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	// create login
	err := resource.provider.manager.CreateLogin(context, data.LoginName.Value, data.Password.Value)
	if err != nil {
		response.Diagnostics.AddError("Failed to create login", err.Error())
		return
	}

	data.Id = types.String{Value: data.LoginName.Value}

	tflog.Trace(context, "created a resource")

	diagnostics = response.State.Set(context, &data)
	response.Diagnostics.Append(diagnostics...)
}

func (resource loginResource) Read(context context.Context, request tfsdk.ReadResourceRequest, response *tfsdk.ReadResourceResponse) {
	var data loginResourceData

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

func (resource loginResource) Update(context context.Context, request tfsdk.UpdateResourceRequest, response *tfsdk.UpdateResourceResponse) {
	var data loginResourceData

	diagnostics := request.Plan.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	// update login
	err := resource.provider.manager.UpdateLogin(context, data.LoginName.Value, data.Password.Value)
	if err != nil {
		response.Diagnostics.AddError("Failed to update login", err.Error())
		return
	}

	data.Id = types.String{Value: data.LoginName.Value}

	diagnostics = response.State.Set(context, &data)
	response.Diagnostics.Append(diagnostics...)
}

func (resource loginResource) Delete(context context.Context, request tfsdk.DeleteResourceRequest, response *tfsdk.DeleteResourceResponse) {
	var data loginResourceData

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
func (r loginResource) ImportState(context context.Context, request tfsdk.ImportResourceStateRequest, response *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStateNotImplemented(context, "", response)
}
