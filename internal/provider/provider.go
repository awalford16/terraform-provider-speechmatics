// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure SpeechmaticsProvider satisfies various provider interfaces.
var _ provider.Provider = &SpeechmaticsProvider{}

// SpeechmaticsProvider defines the provider implementation.
type SpeechmaticsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
	api_key string
}

// SpeechmaticsProviderModel describes the provider data model.
type SpeechmaticsProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	ApiKey   types.String `tfsdk:"api_key"`
}

func (p *SpeechmaticsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "speechmatics"
	resp.Version = p.version
}

func (p *SpeechmaticsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Speechmatics endpoint to use",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key to authenticate with Speechmatics",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *SpeechmaticsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data SpeechmaticsProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check that API key is provided in the provider
	if data.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Speechmatics API key not provided",
			"Speechmatics requires an API key to authenticate",
		)
	}

	endpoint := os.Getenv("SPEECHMATICS_URL")
	api_key := os.Getenv("SPEECHMATICS_API_TOKEN")

	if !data.Endpoint.IsNull() {
		endpoint = data.Endpoint.ValueString()
	}

	if !data.ApiKey.IsNull() {
		api_key = data.ApiKey.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a http client to interact with Speechmatics API
	client := &http.Client{}

	// Validate API token
	get_req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/jobs", endpoint), nil)
	get_req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api_key))

	get_resp, err := client.Do(get_req)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	defer get_resp.Body.Close()

	if get_resp.StatusCode != 200 {
		fmt.Printf("Error validating Speechmatics API token, please ensure you have the correct endpoint and API key set")
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *SpeechmaticsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTranscriptionResource,
	}
}

func (p *SpeechmaticsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
	}
}

func New(version string, api_key string) func() provider.Provider {
	return func() provider.Provider {
		return &SpeechmaticsProvider{
			version: version,
			api_key: api_key,
		}
	}
}
