package main

import (
  "fmt"
  kinesis "github.com/sendgridlabs/go-kinesis"
)

func main() {
  fmt.Println("Begin")

  ksis := kinesis.New("", "")
  filter := kinesis.NewFilter()
  resp, err := ksis.ListStreams(filter)

  fmt.Println(resp)
  fmt.Println(err)

}