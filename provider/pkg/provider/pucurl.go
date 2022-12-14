package provider

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

type PuCurl struct{}


func (c *PuCurl) Create(ctx p.Context, name string, input PuCurlInputs, preview bool) (id string, output PuCurlOutputs, err error) {
	//TODO implement me
	panic("implement me")
}

func (c *PuCurl) Update(ctx p.Context, id string, olds PuCurlOutputs, news PuCurlInputs, preview bool) (PuCurlOutputs, error) {
	//TODO implement me
	panic("implement me")
}

func (c *PuCurl) Delete(ctx p.Context, id string, props PuCurlOutputs) error {
	//TODO implement me
	panic("implement me")
}

func (c *PuCurl) WireDependencies(f infer.FieldSelector, args *PuCurlInputs, state *PuCurlOutputs) {

}


var _ = (infer.Annotated)((*PuCurl)(nil))

func (c *PuCurl) Annotate(a infer.Annotator) {
	a.Describe(&c, "A local command to be executed.\n"+
		"This command can be inserted into the life cycles of other resources using the\n"+
		"`dependsOn` or `parent` resource options. A command is considered to have\n"+
		"failed when it finished with a non-zero exit code. This will fail the CRUD step\n"+
		"of the `Command` resource.")
}

type PuCurlInputs struct {
	// The field tags are used to provide metadata on the schema representation.
	// pulumi:"optional" specifies that a field is optional. This must be a pointer.
	// provider:"replaceOnChanges" specifies that the resource will be replaced if the field changes.
	Name                *string            `pulumi:"name"`
	Url                 *string            `pulumi:"url"`
	Method              *string            `pulumi:"method,optional"`
	Body                *string            `pulumi:"body,optional"`
	Headers             *map[string]string `pulumi:"headers,optional"`
	ResponseCode        *int               `pulumi:"responseCode"`
	DeleteUrl           *bool              `pulumi:"deleteUrl,optional"`
	DeleteMethod        *bool              `pulumi:"deleteMethod,optional"`
	DeleteBody          *bool              `pulumi:"deleteBody,optional"`
	DestroyHeaders      *map[string]string `pulumi:"destroyHeaders,optional"`
	DestroyResponseCode *int               `pulumi:"destroyResponseCode"`
}

func (c *PuCurlInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Name, "The name for this API call.")
	a.Describe(&c.Url, "The API endpoint to call.")
	a.Describe(&c.Method, "The HTTP method to use.")
	a.Describe(&c.Body, "The body of the request.")
	a.Describe(&c.Headers, "The headers to send with the request.")
	a.Describe(&c.ResponseCode, "The expected response code.")
	a.Describe(&c.DeleteUrl, "The API endpoint to call.")
	a.Describe(&c.DeleteMethod, "The HTTP method to use.")
	a.Describe(&c.DeleteBody, "The body of the request.")
	a.Describe(&c.DestroyHeaders, "The headers to send with the request.")
	a.Describe(&c.DestroyResponseCode, "The expected response code.")
}

type PuCurlOutputs struct {
	PuCurlInputs
	Response       *string `pulumi:"response"`
	DeleteResponse *string `pulumi:"deleteResponse"`
}

func (c *PuCurlOutputs) Annotate(a infer.Annotator) {
	c.PuCurlInputs.Annotate(a)
	a.Describe(&c.Response, "The response from the API call.")
	a.Describe(&c.DeleteResponse, "The response from the API call.")
}
