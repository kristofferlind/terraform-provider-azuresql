package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type loginRoleResourceType struct{}

type loginRoleResourceData struct {
	Id        types.String `tfsdk:"id"`
	LoginName types.String `tfsdk:"name"`
	Role      types.String `tfsdk:"role"`
}

type loginRoleResource struct {
	provider provider
}

func (t loginRoleResourceType) GetSchema(context context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Grant server role to login",

		Attributes: map[string]tfsdk.Attribute{
			// needed to keep testing framework happy, just set to login_name
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				MarkdownDescription: "Name of login to grant role",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"role": {
				MarkdownDescription: "Server role to grant login",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
		},
	}, nil
}

func (t loginRoleResourceType) NewResource(context context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diagnostics := convertProviderType(in)

	return loginRoleResource{
		provider: provider,
	}, diagnostics
}

func (resource loginRoleResource) Create(context context.Context, request tfsdk.CreateResourceRequest, response *tfsdk.CreateResourceResponse) {
	var data loginRoleResourceData

	diagnostics := request.Config.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	// create login
	err := resource.provider.manager.GrantServerRole(context, data.LoginName.Value, data.Role.Value)
	if err != nil {
		response.Diagnostics.AddError("Failed to grant server role", err.Error())
		return
	}

	data.Id = types.String{Value: data.LoginName.Value}

	diagnostics = response.State.Set(context, &data)
	response.Diagnostics.Append(diagnostics...)
}

func (resource loginRoleResource) Read(context context.Context, request tfsdk.ReadResourceRequest, response *tfsdk.ReadResourceResponse) {
	var data loginRoleResourceData

	diagnostics := request.State.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	// read login
	login, err := resource.provider.manager.GetServerRole(context, data.LoginName.Value, data.Role.Value)
	if err != nil {
		response.Diagnostics.AddError("Failed to read server roles", err.Error())
		return
	}

	data.LoginName = types.String{Value: login.Name}
	data.Role = types.String{Value: login.Role}

	diagnostics = response.State.Set(context, &data)
	response.Diagnostics.Append(diagnostics...)
}

func (resource loginRoleResource) Update(context context.Context, request tfsdk.UpdateResourceRequest, response *tfsdk.UpdateResourceResponse) {
	var data loginRoleResourceData

	diagnostics := request.Plan.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.AddError("Not implemented", "update should not be allowed on this resource, should always be a replace operation")
}

func (resource loginRoleResource) Delete(context context.Context, request tfsdk.DeleteResourceRequest, response *tfsdk.DeleteResourceResponse) {
	var data loginRoleResourceData

	diagnostics := request.State.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	// delete login
	err := resource.provider.manager.RevokeServerRole(context, data.LoginName.Value, data.Role.Value)
	if err != nil {
		response.Diagnostics.AddError("Failed to revoke server role", err.Error())
		return
	}

	response.State.RemoveResource(context)
}

// import not supported
func (r loginRoleResource) ImportState(context context.Context, request tfsdk.ImportResourceStateRequest, response *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStateNotImplemented(context, "", response)
}
