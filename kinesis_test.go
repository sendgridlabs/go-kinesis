package kinesis

import "testing"

// Trivial test to make sure that Kinesis implements KinesisClient.
func TestInterfaceIsImplemented(t *testing.T) {
	var client KinesisClient
	auth := &Auth{
		AccessKey: "BAD_ACCESS_KEY",
		SecretKey: "BAD_SECRET_KEY",
	}

	client = New(auth, USEast1)
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
