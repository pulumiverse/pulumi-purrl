package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumiverse/pulumi-purrl/sdk/go/purrl"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		purrl, err := purrl.NewPurrl(ctx, "purrl", &purrl.PurrlArgs{
			Url:  pulumi.String("https://httpbin.org/get"),
			Name: pulumi.String("httpbin"),
			ResponseCodes: pulumi.StringArray{
				pulumi.String("200"),
			},
			Method: pulumi.String("GET"),
			Headers: pulumi.StringMap{
				"test": pulumi.String("test"),
			},
			DeleteMethod: pulumi.String("DELETE"),
			DeleteUrl:    pulumi.String("https://httpbin.org/delete"),
			DeleteResponseCodes: pulumi.StringArray{
				pulumi.String("200"),
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("response", purrl.Response)
		return nil
	})
}
