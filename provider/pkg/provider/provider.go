package provider

import (
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumiverse/pulumi-purrl/provider/pkg/provider/purrl"
	"strings"

	"github.com/blang/semver"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-go-provider/integration"
	"github.com/pulumi/pulumi-go-provider/middleware/schema"
)

// NewProvider This provider uses the `pulumi-go-provider` library to produce a code-first provider definition.
func NewProvider() p.Provider {
	return infer.Provider(infer.Options{
		// This is the metadata for the provider
		Metadata: schema.Metadata{
			DisplayName: "purrl",
			Description: "A Pulumi native provider for making API calls",
			Keywords: []string{
				"pulumi",
				"command",
				"category/utility",
				"kind/native",
			},
			Homepage:          "https://pulumi.com",
			License:           "Apache-2.0",
			Repository:        "https://github.com/pulumiverse/pulumi-purrl",
			PluginDownloadURL: "github://api.github.com/pulumiverse",
			Publisher:         "Pulumiverse",
			LogoURL:           "",
			// This contains language specific details for generating the provider's SDKs
			LanguageMap: map[string]any{
				"csharp": map[string]any{
					"packageReferences": map[string]string{
						"Pulumi": "3.*",
					},
					"rootNamespace": "Pulumiverse",
				},
				"go": map[string]any{
					"generateResourceContainerTypes": true,
					"importBasePath":                 "github.com/pulumiverse/pulumi-purrl/sdk/go/purrl",
				},
				"nodejs": map[string]any{
					"dependencies": map[string]string{
						"@pulumi/pulumi": "^3.0.0",
					},
					"packageName": "@pulumiverse/purrl",
				},
				"python": map[string]any{
					"requires": map[string]string{
						"pulumi": ">=3.0.0,<4.0.0",
					},
					"packageName": "pulumiverse_purrl",
				},
				"java": map[string]any{
					"buildFiles":                      "gradle",
					"gradleNexusPublishPluginVersion": "1.1.0",
					"dependencies": map[string]any{
						"com.pulumi:pulumi":               "0.6.0",
						"com.google.code.gson:gson":       "2.8.9",
						"com.google.code.findbugs:jsr305": "3.0.2",
					},
					"basePackage": "com.pulumiverse",
				},
			},
		},
		// A list of `infer.Resource` that are provided by the provider.
		Resources: []infer.InferredResource{
			// The Command resource implementation is commented extensively for new pulumi-go-provider developers.
			infer.Resource[
				// 1. This type is an interface that implements the logic for the Resource
				//    these methods include `Create`, `Update`, `Delete`, and `WireDependencies`.
				//    `WireDependencies` should be implemented to preserve the secretness of an input
				*purrl.Purrl,
				// 2. The type of the Inputs/Arguments to supply to the Resource.
				purrl.PurrlInputs,
				// 3. The type of the Output/Properties/Fields of a created Resource.
				purrl.PurrlOutputs,
			](),
		},
		ModuleMap: map[tokens.ModuleName]tokens.ModuleName{
			"purrl": "index",
		},
	})
}

func Schema(version string) (string, error) {
	version = strings.TrimPrefix(version, "v")
	s, err := integration.NewServer("purrl", semver.MustParse(version), NewProvider()).
		GetSchema(p.GetSchemaRequest{})
	return s.Schema, err
}
