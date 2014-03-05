// Package provide GOlang API for http://aws.amazon.com/kinesis/
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
// Structure for kinesis client
type Kinesis struct {
  client  *Client
  Region  string
  Version string
}

// Initialize new client for AWS Kinesis
func New(access_key, secret_key string) *Kinesis {
  keys := &Auth{
    AccessKey: access_key,
    SecretKey: secret_key,
  }
  return &Kinesis{client: NewClient(keys), Version: "20131202", Region: "us-east-1"}
}

// Create params object for request
func makeParams(action string) map[string]string {
  params := make(map[string]string)
  params[ACTION_KEY] = action
  return params
}
// RequestArgs store params for request
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

// Error represent error from Kinesis API
type Error struct {
  // HTTP status code (200, 403, ...)
  StatusCode int
  // error code ("UnsupportedOperation", ...)
  Code string
  // The human-oriented error message
  Message   string
  RequestId string
}
// Return error message from error object
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

// Query by AWS API
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
  defer response.Body.Close()

  if response.StatusCode != 200 {
    return buildError(response)
  }

  if resp == nil {
    return nil
  }

  return json.NewDecoder(response.Body).Decode(resp)
}

// CreateStream adds a new Amazon Kinesis stream to your AWS account
// StreamName is a name of stream, ShardCount is number of shards
// more info http://docs.aws.amazon.com/kinesis/latest/APIReference/API_CreateStream.html
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

// DeleteStream deletes a stream and all of its shards and data from your AWS account
// StreamName is a name of stream
// more info http://docs.aws.amazon.com/kinesis/latest/APIReference/API_DeleteStream.html
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

// MergeShards merges two adjacent shards in a stream and combines them into a single shard to reduce the stream's capacity to ingest and transport data
// more info http://docs.aws.amazon.com/kinesis/latest/APIReference/API_MergeShards.html
func (kinesis *Kinesis) MergeShards(args *RequestArgs) error {
  params := makeParams("MergeShards")
  err := kinesis.query(params, args.params, nil)
  if err != nil {
    return err
  }
  return nil
}

// SplitShard splits a shard into two new shards in the stream, to increase the stream's capacity to ingest and transport data
// more info http://docs.aws.amazon.com/kinesis/latest/APIReference/API_SplitShard.html
func (kinesis *Kinesis) SplitShard(args *RequestArgs) error {
  params := makeParams("SplitShard")
  err := kinesis.query(params, args.params, nil)
  if err != nil {
    return err
  }
  return nil
}

// ListStreamsResp stores the information that provides by ListStreams API call
type ListStreamsResp struct {
  IsMoreDataAvailable bool
  StreamNames         []string
}

// ListStreams returns an array of the names of all the streams that are associated with the AWS account making the ListStreams request
// more info http://docs.aws.amazon.com/kinesis/latest/APIReference/API_ListStreams.html
func (kinesis *Kinesis) ListStreams(args *RequestArgs) (resp *ListStreamsResp, err error) {
  params := makeParams("ListStreams")
  resp = &ListStreamsResp{}
  err = kinesis.query(params, args.params, resp)
  if err != nil {
    return nil, err
  }
  return
}

// DescribeStreamShards stores the information about list of shards inside DescribeStreamResp
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
// DescribeStreamResp stores the information that provides by DescribeStream API call
type DescribeStreamResp struct {
  StreamDescription struct {
    IsMoreDataAvailable     bool
    Shards                  []DescribeStreamShards
    StreamARN               string
    StreamName              string
    StreamStatus            string
  }
}

// DescribeStream returns the following information about the stream: the current status of the stream,
// the stream Amazon Resource Name (ARN), and an array of shard objects that comprise the stream.
// For each shard object there is information about the hash key and sequence number ranges that
// the shard spans, and the IDs of any earlier shards that played in a role in a MergeShards or
// SplitShard operation that created the shard
// more info http://docs.aws.amazon.com/kinesis/latest/APIReference/API_DescribeStream.html
func (kinesis *Kinesis) DescribeStream(args *RequestArgs) (resp *DescribeStreamResp, err error) {
  params := makeParams("DescribeStream")
  resp = &DescribeStreamResp{}
  err = kinesis.query(params, args.params, resp)
  if err != nil {
    return nil, err
  }
  return
}

// GetShardIteratorResp stores the information that provides by GetShardIterator API call
type GetShardIteratorResp struct {
  ShardIterator         string
}

// GetShardIterator returns a shard iterator
// more info http://docs.aws.amazon.com/kinesis/latest/APIReference/API_GetShardIterator.html
func (kinesis *Kinesis) GetShardIterator(args *RequestArgs) (resp *GetShardIteratorResp, err error) {
  params := makeParams("GetShardIterator")
  resp = &GetShardIteratorResp{}
  err = kinesis.query(params, args.params, resp)
  if err != nil {
    return nil, err
  }
  return
}

// GetNextRecordsRecords stores the information that provides by GetNextRecordsResp
type GetRecordsRecords struct {
  Data                      []byte
  PartitionKey              string
  SequenceNumber            string
}
// GetNextRecordsResp stores the information that provides by GetNextRecords API call
type GetRecordsResp struct {
  NextShardIterator           string
  Records                     []GetRecordsRecords
}

// GetRecords returns one or more data records from a shard
// more info http://docs.aws.amazon.com/kinesis/latest/APIReference/API_GetRecords.html
func (kinesis *Kinesis) GetRecords(args *RequestArgs) (resp *GetRecordsResp, err error) {
  params := makeParams("GetRecords")
  resp = &GetRecordsResp{}
  err = kinesis.query(params, args.params, resp)
  if err != nil {
    return nil, err
  }
  return
}

// PutRecordResp stores the information that provides by PutRecord API call
type PutRecordResp struct {
  SequenceNumber          string
  ShardId                 string
}

// PutRecord puts a data record into an Amazon Kinesis stream from a producer
// more info http://docs.aws.amazon.com/kinesis/latest/APIReference/API_PutRecord.html
func (kinesis *Kinesis) PutRecord(args *RequestArgs) (resp *PutRecordResp, err error) {
  params := makeParams("PutRecord")
  resp = &PutRecordResp{}
  err = kinesis.query(params, args.params, resp)
  if err != nil {
    return nil, err
  }
  return
}
