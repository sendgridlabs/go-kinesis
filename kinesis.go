package kinesis

import (
  "encoding/json"
  "net/http"
  "fmt"
  "strings"
  "time"
)

var (
  timeNow = time.Now
)

type Kinesis struct {
  client  *Client
  Region string
  Version string
}

// New

func New(access_key, secret_key string) *Kinesis {
  keys := &Auth{
    AccessKey: access_key,
    SecretKey: secret_key,
  }
  return &Kinesis{ client: NewClient(keys), Version: "20131104", Region: "us-east-1" }
}

func buildError(r *http.Response) error {
  return nil
}


func (kinesis *Kinesis) query(data interface{}, params map[string]string, resp interface{}) error {
  jsonData, err := json.Marshal(data)
  if err != nil {
    return err
  }
  jsonReader := strings.NewReader(string(jsonData))

  request, err := http.NewRequest("POST", fmt.Sprintf("https://kinesis.%s.amazonaws.com", kinesis.Region), jsonReader)
  if err != nil {
    return err
  }
  // headers
  request.Header.Set("Content-Type", "application/x-amz-json-1.1")
  request.Header.Set("X-Amz-Target", fmt.Sprintf("Kinesis_%s.%s", kinesis.Version, params["Action"]))
  request.Header.Set("User-Agent", "Golang Kinesis")
  // request
  response, err := kinesis.client.Do(request)
  if err != nil {
    return err
  }

  if response.StatusCode != 200 {
    return buildError(response)
  }

  return json.NewDecoder(response.Body).Decode(resp)
}
