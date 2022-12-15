package main

import (
	"github.com/dirien/pulumi-pucurl/sdk/go/pucurl/pucurl"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		puCurl, err := pucurl.NewPuCurl(ctx, "pucurl", &pucurl.PuCurlArgs{
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
		ctx.Export("puCurl", puCurl.Response)
		return nil
	})
}
