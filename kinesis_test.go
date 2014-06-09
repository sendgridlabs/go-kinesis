package kinesis

import "testing"

// Trivial test to make sure that Kinesis implements KinesisClient.
func TestInterfaceIsImplemented(t *testing.T) {
  var client KinesisClient
  client = New("BAD_ACCESS_KEY", "BAD_SECRET_KEY")
  if client == nil {
    t.Error("Client is nil")
  }
}
