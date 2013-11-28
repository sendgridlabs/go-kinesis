package main

import (
  "fmt"
  kinesis "github.com/sendgridlabs/go-kinesis"
)

func main() {
  fmt.Println("Begin")

  streamName := "test"

  ksis := kinesis.New("", "")
  /*
  err := ksis.CreateStream("test", 1)
  fmt.Println(err)
  */
  args := kinesis.NewArgs()
  resp2, _ := ksis.ListStreams(args)
  fmt.Printf("%v", resp2)
  fmt.Println("")

  args = kinesis.NewArgs()
  args.Add("StreamName", streamName)
  resp3, _ := ksis.DescribeStream(args)
  fmt.Printf("%v", resp3)
  fmt.Println("")


  args = kinesis.NewArgs()
  args.Add("StreamName", streamName)
  args.Add("ShardId", resp3.StreamDescription.Shards[0].ShardId)
  resp4, _ := ksis.GetShardIterator(args)
  fmt.Printf("%v", resp4)
  fmt.Println("")

  /*
  err = ksis.DeleteStream("test")
  fmt.Println(err)
  */
}