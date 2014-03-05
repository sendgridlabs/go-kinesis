package main

import (
  "fmt"
  "time"
  kinesis "github.com/sendgridlabs/go-kinesis"
)

func main() {
  fmt.Println("Begin")

  streamName := "test"
  // set env variables AWS_ACCESS_KEY and AWS_SECRET_KEY
  ksis := kinesis.New("", "")

  err := ksis.CreateStream(streamName, 2)
  if err != nil {
    fmt.Printf("CreateStream ERROR: %v\n", err)
  }

  args := kinesis.NewArgs()
  resp2, _ := ksis.ListStreams(args)
  fmt.Printf("ListStreams: %v\n", resp2)

  resp3 := &kinesis.DescribeStreamResp{}

  timeout := make(chan bool, 30)
  for {

    args = kinesis.NewArgs()
    args.Add("StreamName", streamName)
    resp3, _ = ksis.DescribeStream(args)
    fmt.Printf("DescribeStream: %v\n", resp3)

    if resp3.StreamDescription.StreamStatus != "ACTIVE" {
      time.Sleep(4 * time.Second)
      timeout <- true
    } else {
      break
    }

  }


  for i := 0; i < 10; i++ {
    args = kinesis.NewArgs()
    args.Add("StreamName", streamName)
    args.Add("Data", []byte(fmt.Sprintf("Hello AWS Kinesis %d", i)))
    args.Add("PartitionKey", fmt.Sprintf("partitionKey-%d", i))
    resp4, err := ksis.PutRecord(args)
    if err != nil {
      fmt.Printf("PutRecord err: %v\n", err)
    } else {
      fmt.Printf("PutRecord: %v\n", resp4)
    }
  }

  for _, shard := range resp3.StreamDescription.Shards {

    args = kinesis.NewArgs()
    args.Add("StreamName", streamName)
    args.Add("ShardId", shard.ShardId)
    args.Add("ShardIteratorType", "TRIM_HORIZON")
    resp10, _ := ksis.GetShardIterator(args)

    shardIterator := resp10.ShardIterator

    for {
      args = kinesis.NewArgs()
      args.Add("ShardIterator", shardIterator)
      resp11, err := ksis.GetRecords(args)

      if len(resp11.Records) > 0 {
        fmt.Printf("GetRecords Data BEGIN\n")
        for _, d := range resp11.Records {
          fmt.Printf("GetRecords Data: %v\n", string(d.Data))
        }
        fmt.Printf("GetRecords Data END\n")
      }

      if len(resp11.Records) == 0 || err != nil {
        break
      } else if resp11.NextShardIterator == "" {
        break
      }

      shardIterator = resp11.NextShardIterator
    }
  }

  err1 := ksis.DeleteStream("test")
  if err1 != nil {
    fmt.Printf("DeleteStream ERROR: %v\n", err1)
  }

  fmt.Println("End")
}