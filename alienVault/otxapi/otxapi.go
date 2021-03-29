package otxapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"

	"github.com/google/go-querystring/query"
)

const (
	get                  = "GET"
	libraryVersion       = "0.1"
	userAgent            = "go-otx-api/" + libraryVersion
	defaultBaseURL       = "https://otx.alienvault.com/"
	subscriptionsURLPath = "api/v1/pulses/subscribed"
	pulseDetailURLPath   = "api/v1/pulses/"
	userURLPath          = "api/v1/user/"
	apiVersion           = "v1"
)

// A Client manages communication with the OTX API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.  Defaults to the otx.alienvault.com:443.
	// BaseURL should always be specified with a trailing slash
	BaseURL *url.URL

	UserAgent string

	// OTX API Services
	UserDetail  *OTXUserDetailService
	PulseDetail *OTXPulseDetailService
	ThreatIntel *OTXThreatIntelFeedService
}

// Response is a otx API response.  This wraps the standard http.Response
// returned from OTX and provides convenient access to things like
// pagination links.
type Response struct {
	*http.Response
	// RawContent - raw stream
	RawContent []uint8
	// Content - additional way to access the content body of the response.
	Content map[string]interface{} `json:"results,omitempty"`
}

// ListOptions specifies the optional parameters to various List methods that
// support pagination.  our list options: ?limit=50&page_num=1
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"limit,omitempty"`
}

// addOptions adds the parameters in opt as URL query parameters to s.  opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

func (c *OTXPulseDetailService) Get(id_string string) (PulseDetail, Response, error) {
	client := &http.Client{}

	req, _ := http.NewRequest(get, fmt.Sprintf("%s/%s/%s/", defaultBaseURL, pulseDetailURLPath, id_string), nil)
	req.Header.Set("X-OTX-API-KEY", fmt.Sprintf("%s", os.Getenv("X_OTX_API_KEY")))

	response, _ := client.Do(req)
	resp := Response{Response: response}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	pulse_detail := new(PulseDetail)
	json.Unmarshal(contents, &(pulse_detail))
	json.Unmarshal(contents, &(resp.Content))

	return *pulse_detail, resp, err
}

func (c *OTXThreatIntelFeedService) List(opt *ListOptions) (ThreatIntelFeed, Response, error) {
	client := &http.Client{}
	requestpath, err := addOptions(defaultBaseURL+subscriptionsURLPath, opt)
	if err != nil {
		return ThreatIntelFeed{}, Response{}, err
	}

	req, _ := http.NewRequest(get, requestpath, nil)
	req.Header.Set("X-OTX-API-KEY", fmt.Sprintf("%s", os.Getenv("X_OTX_API_KEY")))

	response, _ := client.Do(req)
	resp := Response{Response: response}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	pulse_list := new(ThreatIntelFeed)
	err = json.Unmarshal(contents, &(pulse_list))
	json.Unmarshal(contents, &(resp.Content))
	if err != nil {
		fmt.Println("error not nil on json unmarshall")
		fmt.Println(err)
	}

	return *pulse_list, resp, err
}

func (c *OTXUserDetailService) Get() (UserDetail, *Response, error) {

	req, err := c.client.NewRequest(get, userURLPath, nil)
	if err != nil {
		return UserDetail{}, nil, err
	}
	req.Header.Set("X-OTX-API-KEY", fmt.Sprintf("%s", os.Getenv("X_OTX_API_KEY")))

	userdetail := &UserDetail{}
	resp, err := c.client.Do(req, userdetail)
	if err != nil {
		return UserDetail{}, resp, err
	}
	err = json.Unmarshal(resp.RawContent, &(userdetail))
	json.Unmarshal(resp.RawContent, &(resp.Content))

	return *userdetail, resp, err
}

// NewClient returns a new OTX API client.  If a nil httpClient is
// provided, http.DefaultClient will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}

	c.UserDetail = &OTXUserDetailService{client: c}
	return c
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := newResponse(resp)

	// check response for error
	err = CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	response.RawContent = content
	return response, err
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.  If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if c.UserAgent != "" {
		req.Header.Add("User-Agent", c.UserAgent)
	}
	return req, nil
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// CheckResponse checks the API response for errors, and returns them if
// present.  A response is considered an error if it has a status code outside
// the 200 range.  API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.  Any other
// response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

type Error struct {
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: %v",
		e.Message)
}

type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"detail"` // error message
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message)
}
