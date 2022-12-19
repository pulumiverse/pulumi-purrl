---
title: PuCurl
meta_desc: Provides an overview of the PuCurl Provider for Pulumi.
layout: overview
---

This provider is designed to be a flexible extension of your Pulumi code to make API calls to your target endpoint. PuCurl is useful when a provider does not have a resource or data source that you require, so PuCurl can be used to make substitute API calls.

## Example

{{< chooser language "typescript,python,go,csharp" >}}
{{% choosable language typescript %}}

TODO

{{% /choosable %}}
{{% choosable language python %}}

TODO

{{% /choosable %}}
{{% choosable language go %}}

```go
package main

import (
	"github.com/pulumiverse/pulumi-purrl/sdk/go/pucurl/pucurl"
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
```

{{% /choosable %}}

{{< /chooser >}}
