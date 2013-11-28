package kinesis

import (
  "encoding/json"
  "fmt"
  "net/http"
  "strings"
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

// params
func makeParams(action string) map[string]string {
  params := make(map[string]string)
  params[ACTION_KEY] = action
  return params
}

type RequestArgs struct {
  params     map[string]interface{}
}

// NewFilter creates a new Filter.
func NewArgs() *RequestArgs {
  return &RequestArgs{make(map[string]interface{})}
}

// Add appends a filtering parameter with the given name and value(s).
func (f *RequestArgs) Add(name string, value interface{}) {
  f.params[name] = value
}

// errors

type Error struct {
  // HTTP status code (200, 403, ...)
  StatusCode int
  // error code ("UnsupportedOperation", ...)
  Code string
  // The human-oriented error message
  Message   string
  RequestId string
}

func (err *Error) Error() string {
  if err.Code == "" {
    return err.Message
  }
  return fmt.Sprintf("%s (%s)", err.Message, err.Code)
}

type jsonErrors struct {
  Message   string
}

func buildError(r *http.Response) error {
  errors := jsonErrors{}
  json.NewDecoder(r.Body).Decode(&errors)

  var err Error
  err.Message = errors.Message
  err.StatusCode = r.StatusCode
  if err.Message == "" {
    err.Message = r.Status
  }
  return &err
}

// query

func (kinesis *Kinesis) query(params map[string]string, data interface{}, resp interface{}) error {
  jsonData, err := json.Marshal(data)
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

  if resp == nil {
    return nil
  }

  return json.NewDecoder(response.Body).Decode(resp)
}

// CreateStream

func (kinesis *Kinesis) CreateStream(StreamName string, ShardCount int) error {
  params := makeParams("CreateStream")
  requestParams := struct {
    StreamName string
    ShardCount int
  } {
    StreamName,
    ShardCount,
  }
  err := kinesis.query(params, requestParams, nil)
  if err != nil {
    return err
  }
  return nil
}

// DeleteStream

func (kinesis *Kinesis) DeleteStream(StreamName string) error {
  params := makeParams("DeleteStream")
  requestParams := struct {
    StreamName string
  } {
    StreamName,
  }
  err := kinesis.query(params, requestParams, nil)
  if err != nil {
    return err
  }
  return nil
}

// MergeShards

func (kinesis *Kinesis) MergeShards(args *RequestArgs) error {
  params := makeParams("MergeShards")
  err := kinesis.query(params, args.params, nil)
  if err != nil {
    return err
  }
  return nil
}

// SplitShard

func (kinesis *Kinesis) SplitShard(args *RequestArgs) error {
  params := makeParams("SplitShard")
  err := kinesis.query(params, args.params, nil)
  if err != nil {
    return err
  }
  return nil
}

// ListStreams

type ListStreamsResp struct {
  IsMoreDataAvailable bool
  StreamNames         []string
}

func (kinesis *Kinesis) ListStreams(args *RequestArgs) (resp *ListStreamsResp, err error) {
  params := makeParams("ListStreams")
  resp = &ListStreamsResp{}
  err = kinesis.query(params, args.params, resp)
  if err != nil {
    return nil, err
  }
  return
}

// DescribeStream

type DescribeStreamShards struct {
  AdjacentParentShardId     string
  HashKeyRange struct {
    EndingHashKey           string
    StartingHashKey         string
  }
  ParentShardId             string
  SequenceNumberRange struct {
    EndingHashKey           string
    StartingHashKey         string
  }
  ShardId                   string
}

type DescribeStreamResp struct {
  StreamDescription struct {
    IsMoreDataAvailable     bool
    Shards                  []DescribeStreamShards
    StreamARN               string
    StreamName              string
    StreamStatus            string
  }
}

func (kinesis *Kinesis) DescribeStream(args *RequestArgs) (resp *DescribeStreamResp, err error) {
  params := makeParams("DescribeStream")
  resp = &DescribeStreamResp{}
  err = kinesis.query(params, args.params, resp)
  if err != nil {
    return nil, err
  }
  return
}

// GetShardIterator

type GetShardIteratorResp struct {
  ShardIterator         string
}

func (kinesis *Kinesis) GetShardIterator(args *RequestArgs) (resp *GetShardIteratorResp, err error) {
  params := makeParams("GetShardIterator")
  resp = &GetShardIteratorResp{}
  err = kinesis.query(params, args.params, resp)
  if err != nil {
    return nil, err
  }
  return
}

// GetNextRecords

type GetNextRecordsRecords struct {
  Data                      []byte
  PartitionKey              string
  SequenceNumber            string
}

type GetNextRecordsResp struct {
  NextShardIterator           string
  Records                     []GetNextRecordsRecords
}

func (kinesis *Kinesis) GetNextRecords(args *RequestArgs) (resp *GetNextRecordsResp, err error) {
  params := makeParams("GetNextRecords")
  resp = &GetNextRecordsResp{}
  err = kinesis.query(params, args.params, resp)
  if err != nil {
    return nil, err
  }
  return
}

// PutRecord

type PutRecordResp struct {
  SequenceNumber          string
  ShardId                 string
}

func (kinesis *Kinesis) PutRecord(args *RequestArgs) (resp *PutRecordResp, err error) {
  params := makeParams("PutRecord")
  resp = &PutRecordResp{}
  err = kinesis.query(params, args.params, resp)
  if err != nil {
    return nil, err
  }
  return
}
