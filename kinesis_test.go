package kinesis

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

const localEndpoint = "http://127.0.0.1:4567"

// Trivial test to make sure that Kinesis implements KinesisClient.
func TestInterfaceIsImplemented(t *testing.T) {
	auth := NewAuth("BAD_ACCESS_KEY", "BAD_SECRET_KEY")

	client := New(auth, USEast1)
	if client == nil {
		t.Error("Client is nil")
	}
}

func TestRegions(t *testing.T) {
	if EUWest1.Name != "eu-west-1" {
		t.Errorf("%q != %q", EUWest1.Name, "eu-west-1")
	}
	if USWest2.Name != "us-west-2" {
		t.Errorf("%q != %q", USWest2.Name, "us-west-2")
	}
	if USEast1.Name != "us-east-1" {
		t.Errorf("%q != %q", USEast1.Name, "us-east-1")
	}
}

func TestAddRecord(t *testing.T) {
	args := NewArgs()

	args.AddRecord(
		[]byte("data"),
		"partition_key",
	)

	if len(args.Records) != 1 {
		t.Errorf("%q != %q", len(args.Records), 1)
	}
}

func TestListStreams(t *testing.T) {
	auth := NewAuth("BAD_ACCESS_KEY", "BAD_SECRET_KEY")
	client := NewWithEndpoint(auth, USEast1, localEndpoint)
	resp, err := client.ListStreams(NewArgs())
	if resp == nil {
		t.Error("resp == nil")
	}
	if err != nil {
		t.Errorf("%q != nil", err)
	}
}

func TestCreateStream(t *testing.T) {
	auth := NewAuth("BAD_ACCESS_KEY", "BAD_SECRET_KEY")
	client := NewWithEndpoint(auth, USEast1, localEndpoint)

	streamName := "test2"

	err := client.CreateStream(streamName, 1)
	if err != nil {
		t.Errorf("%q != nil", err)
	}

	err = waitForStreamStatus(client, streamName, "ACTIVE")
	if err != nil {
		t.Errorf("%q != nil", err)
	}

	client.DeleteStream(streamName)
	err = waitForStreamDeletion(client, streamName)
	if err != nil {
		t.Errorf("%q != nil", err)
	}
}

// Older, lower-level way to use PutRecord
func TestPutRecordWithAddData(t *testing.T) {
	auth := NewAuth("BAD_ACCESS_KEY", "BAD_SECRET_KEY")
	client := NewWithEndpoint(auth, USEast1, localEndpoint)

	streamName := "pizza"
	err := createStream(client, streamName, 1)

	if err != nil {
		t.Errorf("%q != nil", err)
	}

	args := NewArgs()
	args.Add("StreamName", streamName)
	args.AddData([]byte("The cheese is old and moldy, where is the bathroom?"))
	args.Add("PartitionKey", "key-1")

	resp, err := client.PutRecord(args)
	if resp == nil {
		t.Error("resp == nil")
	}
	if err != nil {
		t.Errorf("%q != nil", err)
	}

	client.DeleteStream(streamName)
	err = waitForStreamDeletion(client, streamName)
	if err != nil {
		t.Errorf("%q != nil", err)
	}
}

// Newer, higher-level way to use PutRecord
func TestPutRecordWithAddRecord(t *testing.T) {
	auth := NewAuth("BAD_ACCESS_KEY", "BAD_SECRET_KEY")
	client := NewWithEndpoint(auth, USEast1, localEndpoint)

	streamName := "pizza"

	err := createStream(client, streamName, 1)
	if err != nil {
		t.Errorf("%q != nil", err)
	}

	args := NewArgs()
	args.Add("StreamName", streamName)
	args.AddRecord([]byte("The cheese is old and moldy, where is the bathroom?"), "key-1")
	resp, err := client.PutRecord(args)

	if resp == nil {
		t.Error("resp == nil")
	}
	if err != nil {
		t.Errorf("%q != nil", err)
	}

	client.DeleteStream(streamName)
	err = waitForStreamDeletion(client, streamName)
	if err != nil {
		t.Errorf("%q != nil", err)
	}
}

// waitForStreamStatus will poll for a stream status repeatedly, once every MS, for up to 1000 MS,
// blocking until the stream has the desired status. It will return an error if the stream never
// achieves the desired status. If a stream doesn’t exist then an error will be returned.
func waitForStreamStatus(client *Kinesis, streamName string, statusToAwait string) error {
	args := NewArgs()
	args.Add("StreamName", streamName)
	var resp3 *DescribeStreamResp
	var err error

	for i := 1; i < 1000; i++ {
		resp3, err = client.DescribeStream(args)
		if err != nil {
			return err
		}

		if resp3.StreamDescription.StreamStatus == statusToAwait {
			break
		} else {
			time.Sleep(1 * time.Millisecond)
		}
	}

	if resp3 == nil {
		return errors.New("Could not get Stream Description")
	}

	if resp3.StreamDescription.StreamStatus != statusToAwait {
		return errors.New(fmt.Sprintf("Timed out waiting for stream to enter status %v; last status was %v.", statusToAwait, resp3.StreamDescription.StreamStatus))
	}

	return nil
}

// waitForStreamDeletion will poll for a stream status repeatedly, once every MS, for up to 1000 MS,
// blocking until the stream has been deleted. It will return an error if the stream is never deleted
// or some other error occurs. If it succeeds then the return value will be nil.
func waitForStreamDeletion(client *Kinesis, streamName string) error {
	err := waitForStreamStatus(client, streamName, "FOO")
	if !strings.Contains(err.Error(), "not found") {
		return err
	}
	return nil
}

// helper
func createStream(client *Kinesis, streamName string, partitions int) error {
	err := client.CreateStream(streamName, partitions)
	if err != nil {
		return err
	}

	err = waitForStreamStatus(client, streamName, "ACTIVE")
	if err != nil {
		return err
	}

	return nil
}
