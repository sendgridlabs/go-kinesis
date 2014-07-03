package kinesis

import "testing"

// Trivial test to make sure that Kinesis implements KinesisClient.
func TestInterfaceIsImplemented(t *testing.T) {
	var client KinesisClient
	client = New("BAD_ACCESS_KEY", "BAD_SECRET_KEY", USEast)
	if client == nil {
		t.Error("Client is nil")
	}
}

func TestRegions(t *testing.T) {
	if EUWest.Name != "eu-west-1" {
		t.Errorf("%q != %q", EUWest.Name, "eu-west-1")
	}
	if USWest2.Name != "us-west-2" {
		t.Errorf("%q != %q", USWest2.Name, "us-west-2")
	}
	if USEast.Name != "us-east-1" {
		t.Errorf("%q != %q", USEast.Name, "us-east-1")
	}
}
