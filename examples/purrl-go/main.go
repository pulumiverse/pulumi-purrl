package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumiverse/pulumi-purrl/sdk/go/purrl"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		purrlCommand, err := purrl.NewPurrl(ctx, "httpbin", &purrl.PurrlArgs{
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
		ctx.Export("actual response code", purrlCommand.ResponseCodes)
		return nil
	})
}
