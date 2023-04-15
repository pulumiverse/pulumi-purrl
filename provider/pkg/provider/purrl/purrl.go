package purrl

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
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

var _ = (infer.CustomResource[PurrlInputs, PurrlOutputs])((*Purrl)(nil))
var _ = (infer.CustomUpdate[PurrlInputs, PurrlOutputs])((*Purrl)(nil))
var _ = (infer.CustomDelete[PurrlOutputs])((*Purrl)(nil))
var _ = (infer.ExplicitDependencies[PurrlInputs, PurrlOutputs])((*Purrl)(nil))

func (c *Purrl) Check(ctx p.Context, name string, oldInputs, newInputs resource.PropertyMap) (
	PurrlInputs, []p.CheckFailure, error) {
	_, hasExpectedResponseCodes := newInputs["expectedResponseCodes"]
	responseCodes, hasResponseCodes := newInputs["responseCodes"]
	if hasExpectedResponseCodes && hasResponseCodes {
		failure := p.CheckFailure{Property: "expectedResponseCodes", Reason: "Only one of responseCodes or expectedResponseCodes must be provided"}
		return PurrlInputs{}, []p.CheckFailure{failure}, errors.New("only one of responseCodes or expectedResponseCodes must be provided")
	} else if !hasExpectedResponseCodes && !hasResponseCodes {
		failure := p.CheckFailure{Property: "expectedResponseCodes", Reason: "At least one of responseCodes or expectedResponseCodes must be provided"}
		return PurrlInputs{}, []p.CheckFailure{failure}, errors.New("at least one of responseCodes or expectedResponseCodes must be provided")
	} else if hasResponseCodes {
		newInputs["expectedResponseCodes"] = responseCodes
	}

	_, hasExpectedDeleteResponseCodes := newInputs["expectedDeletedResponseCodes"]
	deleteResponseCodes, hasDeleteResponseCodes := newInputs["deleteResponseCodes"]
	if hasExpectedDeleteResponseCodes && hasDeleteResponseCodes {
		failure := p.CheckFailure{Property: "expectedDeleteResponseCodes", Reason: "Only one of deleteResponseCodes or expectedDeleteResponseCodes must be provided"}
		return PurrlInputs{}, []p.CheckFailure{failure}, errors.New("only one of deleteResponseCodes or expectedDeleteResponseCodes must be provided")
	} else if !hasExpectedDeleteResponseCodes && !hasDeleteResponseCodes {
		failure := p.CheckFailure{Property: "expectedDeleteResponseCodes", Reason: "At least one of deleteResponseCodes or expectedDeleteResponseCodes must be provided"}
		return PurrlInputs{}, []p.CheckFailure{failure}, errors.New("at least one of deleteResponseCodes or expectedDeleteResponseCodes must be provided")
	} else if hasDeleteResponseCodes {
		newInputs["expectedDeleteResponseCodes"] = deleteResponseCodes
	}

	return infer.DefaultCheck[PurrlInputs](newInputs)
}

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
	// At this point, after `Check` has been called, we should be safe in only using input.ExpectedResponseCodes
	code, endpoint, err := callAPIEndpoint(input.Method, input.URL, input.Body, input.ExpectedResponseCodes, input.Headers,
		input.InsecureSkipTLSVerify, input.CaCert, input.Cert, input.Key)
	if err != nil {
		return id, state, err
	}
	state.ResponseCode = code
	state.Response = endpoint

	return id, state, err
}

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}

func callAPIEndpoint(method, url, body *string, expectedResponseCodes *[]string, headers *map[string]string,
	insecureSkipVerify *bool, caCert, cert, key *string) (*int, *string, error) {
	if method == nil || url == nil || expectedResponseCodes == nil {
		return nil, nil, errors.New("method, url and expectedResponseCodes are required")
	}

	if insecureSkipVerify == nil {
		insecureSkipVerify = boolPtr(false)
	}

	if body == nil {
		body = strPtr("")
	}

	request, err := http.NewRequestWithContext(context.TODO(), *method, *url, bytes.NewBuffer([]byte(*body)))
	if err != nil {
		return nil, nil, fmt.Errorf("error creating request: %v", err)
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

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: *insecureSkipVerify,
			},
		},
	}

	if cert != nil && key != nil {
		certificate, err := tls.X509KeyPair([]byte(*cert), []byte(*key))
		if err != nil {
			return nil, nil, fmt.Errorf("error creating certificate: %v", err)
		}
		caCertPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, nil, fmt.Errorf("error creating SystemCertPool: %v", err)
		}
		if caCert == nil {
			caCertPool = x509.NewCertPool()
		} else {
			caCertPool.AppendCertsFromPEM([]byte(*caCert))
		}
		caCertPool.AppendCertsFromPEM([]byte(*cert))
		client.Transport.(*http.Transport).TLSClientConfig.RootCAs = caCertPool
		client.Transport.(*http.Transport).TLSClientConfig.Certificates = []tls.Certificate{certificate}
	}

	resp, err := client.Do(request)
	if err != nil {
		return intPtr(resp.StatusCode), nil, fmt.Errorf("error calling API endpoint: %v", err)
	}
	defer resp.Body.Close()
	respBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return intPtr(resp.StatusCode), nil, fmt.Errorf("error reading response body: %v", readErr)
	}

	code := fmt.Sprintf("%v", resp.StatusCode)

	stringConversionList := make([]string, len(*expectedResponseCodes))
	for i, v := range *expectedResponseCodes {
		stringConversionList[i] = fmt.Sprint(v)
	}

	if !responseCodeChecker(stringConversionList, code) {
		return intPtr(resp.StatusCode), nil, errors.New("response code not in list of expected response codes")
	}

	return intPtr(resp.StatusCode), strPtr(string(respBody)), nil
}

func responseCodeChecker(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func (c *Purrl) Update(ctx p.Context, id string, olds PurrlOutputs,
	news PurrlInputs, preview bool) (PurrlOutputs, error) {
	state := PurrlOutputs{PurrlInputs: news}
	// If in preview, don't run the command.
	if preview {
		return state, nil
	}
	code, endpoint, err := callAPIEndpoint(news.Method, news.URL, news.Body, news.ExpectedResponseCodes, news.Headers,
		news.InsecureSkipTLSVerify, news.CaCert, news.Cert, news.Key)
	if err != nil {
		return state, err
	}
	state.ResponseCode = code
	state.Response = endpoint
	return state, nil
}

func (c *Purrl) Delete(ctx p.Context, id string, props PurrlOutputs) error {

	// if delete props are not set, we do nothing
	if props.DeleteMethod == nil || props.DeleteURL == nil || props.ExpectedDeleteResponseCodes == nil {
		return nil
	}
	code, deleteResponse, err := callAPIEndpoint(props.DeleteMethod, props.DeleteURL, props.DeleteBody,
		props.ExpectedDeleteResponseCodes, props.DeleteHeaders, props.DeleteInsecureSkipTLSVerify,
		props.DeleteCaCert, props.DeleteCert, props.DeleteKey)
	if err != nil {
		return err
	}
	ctx.Logf(diag.Debug, "delete response: %s %s", *code, *deleteResponse)

	return nil
}

func (c *Purrl) WireDependencies(f infer.FieldSelector, args *PurrlInputs, state *PurrlOutputs) {
	nameInput := f.InputField(&args.Name)
	urlInput := f.InputField(&args.URL)
	methodInput := f.InputField(&args.Method)
	bodyInput := f.InputField(&args.Body)
	headersInput := f.InputField(&args.Headers)
	responseCodeInput := f.InputField(&args.ExpectedResponseCodes)
	insecureSkipTLSVerifyInput := f.InputField(&args.InsecureSkipTLSVerify)
	caCertInput := f.InputField(&args.CaCert)
	certInput := f.InputField(&args.Cert)
	keyInput := f.InputField(&args.Key)

	deleteURLInput := f.InputField(&args.DeleteURL)
	deleteMethodInput := f.InputField(&args.DeleteMethod)
	deleteBodyInput := f.InputField(&args.DeleteBody)
	deleteHeadersInput := f.InputField(&args.DeleteHeaders)
	deleteResponseCodeInput := f.InputField(&args.ExpectedDeleteResponseCodes)
	deleteInsecureSkipTLSVerifyInput := f.InputField(&args.DeleteInsecureSkipTLSVerify)
	deleteCaCertInput := f.InputField(&args.DeleteCaCert)
	deleteCertInput := f.InputField(&args.DeleteCert)
	deleteKeyInput := f.InputField(&args.DeleteKey)

	f.OutputField(&state.Name).DependsOn(nameInput)
	f.OutputField(&state.URL).DependsOn(urlInput)
	f.OutputField(&state.Method).DependsOn(methodInput)
	f.OutputField(&state.Body).DependsOn(bodyInput)
	f.OutputField(&state.Headers).DependsOn(headersInput)
	f.OutputField(&state.ExpectedResponseCodes).DependsOn(responseCodeInput)
	f.OutputField(&state.InsecureSkipTLSVerify).DependsOn(insecureSkipTLSVerifyInput)
	f.OutputField(&state.CaCert).DependsOn(caCertInput)
	f.OutputField(&state.Cert).DependsOn(certInput)
	f.OutputField(&state.Key).DependsOn(keyInput)
	f.OutputField(&state.DeleteURL).DependsOn(deleteURLInput)
	f.OutputField(&state.DeleteMethod).DependsOn(deleteMethodInput)
	f.OutputField(&state.DeleteBody).DependsOn(deleteBodyInput)
	f.OutputField(&state.DeleteHeaders).DependsOn(deleteHeadersInput)
	f.OutputField(&state.ExpectedDeleteResponseCodes).DependsOn(deleteResponseCodeInput)
	f.OutputField(&state.DeleteInsecureSkipTLSVerify).DependsOn(deleteInsecureSkipTLSVerifyInput)
	f.OutputField(&state.DeleteCaCert).DependsOn(deleteCaCertInput)
	f.OutputField(&state.DeleteCert).DependsOn(deleteCertInput)
	f.OutputField(&state.DeleteKey).DependsOn(deleteKeyInput)
}

var _ = (infer.Annotated)((*Purrl)(nil))

func (c *Purrl) Annotate(a infer.Annotator) {
	a.Describe(&c, "A Pulumi provider for making API calls")
}

type PurrlInputs struct {
	// The field tags are used to provide metadata on the schema representation.
	// pulumi:"optional" specifies that a field is optional. This must be a pointer.
	// provider:"replaceOnChanges" specifies that the resource will be replaced if the field changes.
	Name                  *string            `pulumi:"name"`
	URL                   *string            `pulumi:"url"`
	Method                *string            `pulumi:"method"`
	Body                  *string            `pulumi:"body,optional"`
	Headers               *map[string]string `pulumi:"headers,optional"`
	ResponseCodes         *[]string          `pulumi:"responseCodes,optional"`
	ExpectedResponseCodes *[]string          `pulumi:"expectedResponseCodes,optional"`
	InsecureSkipTLSVerify *bool              `pulumi:"insecureSkipTLSVerify,optional"`
	CaCert                *string            `pulumi:"caCert,optional"`
	Cert                  *string            `pulumi:"cert,optional"`
	Key                   *string            `pulumi:"key,optional"`

	DeleteURL                   *string            `pulumi:"deleteUrl,optional"`
	DeleteMethod                *string            `pulumi:"deleteMethod,optional"`
	DeleteBody                  *string            `pulumi:"deleteBody,optional"`
	DeleteHeaders               *map[string]string `pulumi:"deleteHeaders,optional"`
	DeleteResponseCodes         *[]string          `pulumi:"deleteResponseCodes,optional"`
	ExpectedDeleteResponseCodes *[]string          `pulumi:"expectedDeleteResponseCodes,optional"`
	DeleteInsecureSkipTLSVerify *bool              `pulumi:"deleteInsecureSkipTLSVerify,optional"`
	DeleteCaCert                *string            `pulumi:"deleteCaCert,optional"`
	DeleteCert                  *string            `pulumi:"deleteCert,optional"`
	DeleteKey                   *string            `pulumi:"deleteKey,optional"`
}

func (c *PurrlInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Name, "The name for this API call.")
	a.Describe(&c.URL, "The API endpoint to call.")
	a.Describe(&c.Method, "The HTTP method to use.")
	a.Describe(&c.Body, "The body of the request.")
	a.Describe(&c.Headers, "The headers to send with the request.")
	a.Describe(&c.ResponseCodes, "The expected response code(s). Deprecated -- use `expectedResponseCodes` instead.")
	a.Describe(&c.ExpectedResponseCodes, "The expected response code(s).")
	a.Describe(&c.InsecureSkipTLSVerify, "Skip TLS verification.")
	a.Describe(&c.CaCert, "The CA certificate if server cert is not signed by a trusted CA.")
	a.Describe(&c.Cert, "The client certificate to use for TLS verification.")
	a.Describe(&c.Key, "The client key to use for TLS verification.")
	a.Describe(&c.DeleteURL, "The API endpoint to call.")
	a.Describe(&c.DeleteMethod, "The HTTP method to use.")
	a.Describe(&c.DeleteBody, "The body of the request.")
	a.Describe(&c.DeleteHeaders, "The headers to send with the request.")
	a.Describe(&c.DeleteResponseCodes, "The expected response code(s) for deletion. Deprecated -- use `expectedDeleteResponseCodes` instead.")
	a.Describe(&c.ExpectedDeleteResponseCodes, "The expected response code(s) for deletion.")
	a.Describe(&c.DeleteInsecureSkipTLSVerify, "Skip TLS verification.")
	a.Describe(&c.DeleteCaCert, "The CA certificate if server cert is not signed by a trusted CA.")
	a.Describe(&c.DeleteCert, "The client certificate to use for TLS verification.")
	a.Describe(&c.DeleteKey, "The client key to use for TLS verification.")
}

type PurrlOutputs struct {
	PurrlInputs
	Response       *string `pulumi:"response"`
	ResponseCode   *int    `pulumi:"responseCode"`
	DeleteResponse *string `pulumi:"deleteResponse,optional"`
}

func (c *PurrlOutputs) Annotate(a infer.Annotator) {
	c.PurrlInputs.Annotate(a)
	a.Describe(&c.Response, "The response from the API call.")
	a.Describe(&c.DeleteResponse, "The response from the (delete) API call.")
}
