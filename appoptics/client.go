package appoptics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Operations is a constant used to describe certain metric units.
// i.e 5 "operations" per second, etc.
const Operations = "operations"

// OperationsShort is a shortened alternative to Operations.
const OperationsShort = "ops"

// PartialFailureHeader is used to search for partial failure HTTP headers.
// If present in a request, then an error has occured.
const PartialFailureHeader = "x-partial-failure"

// Client holds the API Token and Measurements url of an AppOptics client.
// It is used for interfacing with the AppOptics measurement API.
type Client struct {
	Token           string
	MeasurementsURI string
}

// NewClient creates a new client for interfacing with the AppOptics API.
func NewClient(token, measurementsURI string) *Client {
	if measurementsURI == "" {
		measurementsURI = DefaultMeasurementsURI
	}
	return &Client{token, measurementsURI}
}

// Property string constants.
const (
	// Display attributes.
	Color             = "color"
	DisplayMax        = "display_max"
	DisplayMin        = "display_min"
	DisplayUnitsLong  = "display_units_long"
	DisplayUnitsShort = "display_units_short"
	DisplayStacked    = "display_stacked"
	DisplayTransform  = "display_transform"

	// Special gauge display attributes.
	SummarizeFunction = "summarize_function"
	Aggregate         = "aggregate"

	// Metric keys.
	Name        = "name"
	Period      = "period"
	Description = "description"
	DisplayName = "display_name"
	Attributes  = "attributes"

	// Measurement keys.
	Time  = "time"
	Tags  = "tags"
	Value = "value"

	// Special Gauge keys.
	Count  = "count"
	Sum    = "sum"
	Max    = "max"
	Min    = "min"
	StdDev = "stddev"

	// Batch keys.
	Measurements = "measurements"

	DefaultMeasurementsURI = "https://api.appoptics.com/v1/measurements"
)

// Measurement is a named metric.
type Measurement map[string]interface{}

// Batch is a collection of measurements and tags.
type Batch struct {
	Measurements []Measurement     `json:"measurements,omitempty"`
	Time         int64             `json:"time"`
	Tags         map[string]string `json:"tags"`
}

var httpClient = http.DefaultClient

// SetHTTPClient modifies the HTTP client used to reach the AppOptics API.
func SetHTTPClient(c *http.Client) {
	httpClient = c
}

// PostMetrics uploads a single batch of metrics to the AppOptics API.
func (app *Client) PostMetrics(batch Batch) (err error) {
	var (
		js   []byte
		req  *http.Request
		resp *http.Response
	)

	if len(batch.Measurements) == 0 {
		return nil
	}

	if js, err = json.Marshal(batch); err != nil {
		return
	}

	if req, err = http.NewRequest("POST", app.MeasurementsURI,
		bytes.NewBuffer(js)); err != nil {
		return
	}

	req.Header.Set("content-type", "application/json")
	req.SetBasicAuth(app.Token, "")

	if resp, err = httpClient.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted || resp.Header.Get(PartialFailureHeader) != "" {
		var body []byte
		if body, err = io.ReadAll(resp.Body); err != nil {
			body = []byte(fmt.Sprintf("(could not fetch response body for error: %s)", err))
		}
		err = fmt.Errorf("unable to post to AppOptics: %d %s %s", resp.StatusCode, resp.Status, string(body))
	}
	return
}
