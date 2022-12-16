package purrl

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"io"
	"net/http"
)

type Purrl struct{}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

var _ = (infer.CustomResource[PurrlInputs, PurrlOutputs])((*Purrl)(nil))
var _ = (infer.CustomUpdate[PurrlInputs, PurrlOutputs])((*Purrl)(nil))
var _ = (infer.CustomDelete[PurrlOutputs])((*Purrl)(nil))
var _ = (infer.ExplicitDependencies[PurrlInputs, PurrlOutputs])((*Purrl)(nil))

func (c *Purrl) Create(ctx p.Context, name string, input PurrlInputs, preview bool) (string, PurrlOutputs, error) {
	state := PurrlOutputs{PurrlInputs: input, Response: strPtr("")}
	var id string
	id, err := resource.NewUniqueHex(name, 8, 0)
	if err != nil {
		return id, state, err
	}
	if preview {
		return id, state, nil
	}
	endpoint, err := callAPIEndpoint(input.Method, input.Url, input.Body, input.ResponseCodes, input.Headers)
	if err != nil {
		return id, state, err
	}
	state.Response = endpoint

	return id, state, err
}

func strPtr(s string) *string {
	return &s
}

func callAPIEndpoint(method, url, body *string, responseCode *[]string, headers *map[string]string) (*string, error) {
	if method == nil || url == nil || responseCode == nil {
		return nil, errors.New("method, url and responseCode are required")
	}

	if body == nil {
		body = strPtr("")
	}

	request, err := http.NewRequestWithContext(context.TODO(), *method, *url, bytes.NewBuffer([]byte(*body)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	if headers != nil {
		headersMap := make(map[string]string)
		for k, v := range *headers {
			strKey := fmt.Sprintf("%v", k)
			strValue := fmt.Sprintf("%v", v)
			headersMap[strKey] = strValue
		}

		for k, v := range headersMap {
			request.Header.Set(k, v)
		}
	}

	resp, err := Client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error calling API endpoint: %v", err)
	}
	defer resp.Body.Close()
	respBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("error reading response body: %v", readErr)
	}

	code := fmt.Sprintf("%v", resp.StatusCode)

	stringConversionList := make([]string, len(*responseCode))
	for i, v := range *responseCode {
		stringConversionList[i] = fmt.Sprint(v)
	}

	if !responseCodeChecker(stringConversionList, code) {
		return nil, errors.New("response code not in list of expected response codes")
	}

	return strPtr(string(respBody)), nil
}

func responseCodeChecker(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func (c *Purrl) Update(ctx p.Context, id string, olds PurrlOutputs, news PurrlInputs, preview bool) (PurrlOutputs, error) {
	state := PurrlOutputs{PurrlInputs: news}
	// If in preview, don't run the command.
	if preview {
		return state, nil
	}
	endpoint, err := callAPIEndpoint(news.Method, news.Url, news.Body, news.ResponseCodes, news.Headers)
	if err != nil {
		return state, err
	}
	state.Response = endpoint
	return state, nil
}

func (c *Purrl) Delete(ctx p.Context, id string, props PurrlOutputs) error {

	// if delete props are not set, we do nothing
	if props.DeleteMethod == nil || props.DeleteUrl == nil || props.DeleteResponseCodes == nil {
		return nil
	}
	deleteResponse, err := callAPIEndpoint(props.DeleteMethod, props.DeleteUrl, props.DeleteBody, props.DeleteResponseCodes, props.DeleteHeaders)
	if err != nil {
		return err
	}
	ctx.Logf(diag.Debug, "delete response: %s", *deleteResponse)

	return nil
}

func (c *Purrl) WireDependencies(f infer.FieldSelector, args *PurrlInputs, state *PurrlOutputs) {
	nameInput := f.InputField(&args.Name)
	urlInput := f.InputField(&args.Url)
	methodInput := f.InputField(&args.Method)
	bodyInput := f.InputField(&args.Body)
	headersInput := f.InputField(&args.Headers)
	responseCodeInput := f.InputField(&args.ResponseCodes)

	deleteUrlInput := f.InputField(&args.DeleteUrl)
	deleteMethodInput := f.InputField(&args.DeleteMethod)
	deleteBodyInput := f.InputField(&args.DeleteBody)
	deleteHeadersInput := f.InputField(&args.DeleteHeaders)
	deleteResponseCodeInput := f.InputField(&args.DeleteResponseCodes)

	f.OutputField(&state.Name).DependsOn(nameInput)
	f.OutputField(&state.Url).DependsOn(urlInput)
	f.OutputField(&state.Method).DependsOn(methodInput)
	f.OutputField(&state.Body).DependsOn(bodyInput)
	f.OutputField(&state.Headers).DependsOn(headersInput)
	f.OutputField(&state.ResponseCodes).DependsOn(responseCodeInput)
	f.OutputField(&state.DeleteUrl).DependsOn(deleteUrlInput)
	f.OutputField(&state.DeleteMethod).DependsOn(deleteMethodInput)
	f.OutputField(&state.DeleteBody).DependsOn(deleteBodyInput)
	f.OutputField(&state.DeleteHeaders).DependsOn(deleteHeadersInput)
	f.OutputField(&state.DeleteResponseCodes).DependsOn(deleteResponseCodeInput)
}

var _ = (infer.Annotated)((*Purrl)(nil))

func (c *Purrl) Annotate(a infer.Annotator) {
	a.Describe(&c, "A Pulumi provider for making API calls")
}

type PurrlInputs struct {
	// The field tags are used to provide metadata on the schema representation.
	// pulumi:"optional" specifies that a field is optional. This must be a pointer.
	// provider:"replaceOnChanges" specifies that the resource will be replaced if the field changes.
	Name          *string            `pulumi:"name"`
	Url           *string            `pulumi:"url"`
	Method        *string            `pulumi:"method"`
	Body          *string            `pulumi:"body,optional"`
	Headers       *map[string]string `pulumi:"headers,optional"`
	ResponseCodes *[]string          `pulumi:"responseCodes"`

	DeleteUrl           *string            `pulumi:"deleteUrl,optional"`
	DeleteMethod        *string            `pulumi:"deleteMethod,optional"`
	DeleteBody          *string            `pulumi:"deleteBody,optional"`
	DeleteHeaders       *map[string]string `pulumi:"deleteHeaders,optional"`
	DeleteResponseCodes *[]string          `pulumi:"deleteResponseCodes,optional"`
}

func (c *PurrlInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Name, "The name for this API call.")
	a.Describe(&c.Url, "The API endpoint to call.")
	a.Describe(&c.Method, "The HTTP method to use.")
	a.Describe(&c.Body, "The body of the request.")
	a.Describe(&c.Headers, "The headers to send with the request.")
	a.Describe(&c.ResponseCodes, "The expected response code.")
	a.Describe(&c.DeleteUrl, "The API endpoint to call.")
	a.Describe(&c.DeleteMethod, "The HTTP method to use.")
	a.Describe(&c.DeleteBody, "The body of the request.")
	a.Describe(&c.DeleteHeaders, "The headers to send with the request.")
	a.Describe(&c.DeleteResponseCodes, "The expected response code.")
}

type PurrlOutputs struct {
	PurrlInputs
	Response       *string `pulumi:"response"`
	DeleteResponse *string `pulumi:"deleteResponse,optional"`
}

func (c *PurrlOutputs) Annotate(a infer.Annotator) {
	c.PurrlInputs.Annotate(a)
	a.Describe(&c.Response, "The response from the API call.")
	a.Describe(&c.DeleteResponse, "The response from the API call.")
}
