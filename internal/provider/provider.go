package provider

import (
	"context"
	"fmt"

	"github.com/kristofferlind/terraform-provider-azuresql/internal/manager"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type provider struct {
	configured bool
	version    string
	manager    manager.Manager
}

type providerData struct {
	ConnectionString types.String `tfsdk:"connection_string"`
}

func (provider *provider) Configure(context context.Context, request tfsdk.ConfigureProviderRequest, response *tfsdk.ConfigureProviderResponse) {
	var data providerData
	diagnostics := request.Config.Get(context, &data)
	response.Diagnostics.Append(diagnostics...)

	if response.Diagnostics.HasError() {
		return
	}

	provider.manager = manager.GetManager(data.ConnectionString.Value)

	provider.configured = true
}

func (provider *provider) GetResources(context context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"azuresql_login":      loginResourceType{},
		"azuresql_aad_login":  externalLoginResourceType{},
		"azuresql_login_user": loginUserResourceType{},
		"azuresql_user":       userResourceType{},
		"azuresql_aad_user":   externalUserResourceType{},
		"azuresql_user_role":  userRoleResourceType{},
	}, nil
}

func (provider *provider) GetDataSources(context context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		// nothing planned
	}, nil
}

func (provider *provider) GetSchema(context context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"connection_string": {
				MarkdownDescription: "For connecting to MSSQL Server",
				Required:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
}

func convertProviderType(in tfsdk.Provider) (provider, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	p, ok := in.(*provider)

	if !ok {
		diagnostics.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return provider{}, diagnostics
	}

	if p == nil {
		diagnostics.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return provider{}, diagnostics
	}

	return *p, diagnostics
}
