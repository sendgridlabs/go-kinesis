package kinesis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
  "strconv"
  "sort"
	"time"
)

const (
	ACTION_KEY = "Action"
)

var (
	timeNow = time.Now
)

type Kinesis struct {
	client  *Client
	Region  string
	Version string
}

// New

func New(access_key, secret_key string) *Kinesis {
	keys := &Auth{
		AccessKey: access_key,
		SecretKey: secret_key,
	}
	return &Kinesis{client: NewClient(keys), Version: "20131104", Region: "us-east-1"}
}

// filter

type Filter struct {
	m map[string][]string
}

// NewFilter creates a new Filter.
func NewFilter() *Filter {
	return &Filter{make(map[string][]string)}
}

// Add appends a filtering parameter with the given name and value(s).
func (f *Filter) Add(name string, value ...string) {
	f.m[name] = append(f.m[name], value...)
}

func (f *Filter) addParams(params map[string]string) {
	if f != nil {
		a := make([]string, len(f.m))
		i := 0
		for k := range f.m {
			a[i] = k
			i++
		}
		sort.StringSlice(a).Sort()
		for i, k := range a {
			prefix := "Filter." + strconv.Itoa(i+1)
			params[prefix+".Name"] = k
			for j, v := range f.m[k] {
				params[prefix+".Value."+strconv.Itoa(j+1)] = v
			}
		}
	}
}

// params
func makeParams(action string) map[string]string {
	params := make(map[string]string)
	params[ACTION_KEY] = action
	return params
}

// errors

type Error struct {
  // HTTP status code (200, 403, ...)
  StatusCode int
  // EC2 error code ("UnsupportedOperation", ...)
  Code string
  // The human-oriented error message
  Message   string
  RequestId string `xml:"RequestID"`
}

func (err *Error) Error() string {
  if err.Code == "" {
    return err.Message
  }
  return fmt.Sprintf("%s (%s)", err.Message, err.Code)
}

type jsonErrors struct {
  RequestId string
  Errors    []Error
}

func buildError(r *http.Response) error {
  errors := jsonErrors{}
  json.NewDecoder(r.Body).Decode(&errors)

  var err Error
  if len(errors.Errors) > 0 {
    err = errors.Errors[0]
  }
  err.RequestId = errors.RequestId
  err.StatusCode = r.StatusCode
  if err.Message == "" {
    err.Message = r.Status
  }
  return &err
	return nil
}

// query

func (kinesis *Kinesis) query(params map[string]string, resp interface{}) error {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return err
	}
	// request
	request, err := http.NewRequest("POST", fmt.Sprintf("https://kinesis.%s.amazonaws.com", kinesis.Region), strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}
	// headers
	request.Header.Set("Content-Type", "application/x-amz-json-1.1")
	request.Header.Set("X-Amz-Target", fmt.Sprintf("Kinesis_%s.%s", kinesis.Version, params[ACTION_KEY]))
	request.Header.Set("User-Agent", "Golang Kinesis")
	// response
	response, err := kinesis.client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return buildError(response)
	}

	return json.NewDecoder(response.Body).Decode(resp)
}

// ListStreams

type StreamsResp struct {
	IsMoreDataAvailable bool
	StreamNames         []string
}

func (kinesis *Kinesis) ListStreams(filter *Filter) (resp *StreamsResp, err error) {
	params := makeParams("ListStreams")
	filter.addParams(params)
	resp = &StreamsResp{}
	err = kinesis.query(params, resp)
	if err != nil {
		return nil, err
	}
	return
}
