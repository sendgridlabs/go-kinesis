package kinesis

import (
	"net/http"
	"strings"
	"testing"
)

var testSignFactoryData = []struct {
	AWS_KEY    string
	AWS_SECRET string
	TOKEN      string
	DateHeader string
	AuthHeader string
}{
	{"ASWKEY", "AWSSECRET", "TOKEN1", "Thu, 28 Nov 2013 15:04:05 GMT", "AWS4-HMAC-SHA256 Credential=ASWKEY/20131128/us-east-1/kinesis/aws4_request, SignedHeaders=content-type;date;host;user-agent;x-amz-target, Signature=6c21aca39f1d4afd383fbc45dd3a580192036162f74bf9fda6cad6c6fb7cde2f"},
	{"ASWKEY2", "AWSSECRET2", "TOKEN2", "Thu, 28 Nov 2013 15:04:05 GMT", "AWS4-HMAC-SHA256 Credential=ASWKEY2/20131128/us-east-1/kinesis/aws4_request, SignedHeaders=content-type;date;host;user-agent;x-amz-target, Signature=488ee09d2d56e747beb5653064d7976cb67136a2afa6013d82ff36d6ae95d263"},
	{"ASWNEWKEY", "AWSSECRET", "TOKEN3", "Thu, 28 Nov 2013 15:04:05 GMT", "AWS4-HMAC-SHA256 Credential=ASWNEWKEY/20131128/us-east-1/kinesis/aws4_request, SignedHeaders=content-type;date;host;user-agent;x-amz-target, Signature=6c21aca39f1d4afd383fbc45dd3a580192036162f74bf9fda6cad6c6fb7cde2f"},
	{"ASWKEY", "AWSSECRET", "TOKEN4", "Mon, 25 Nov 2013 15:04:05 GMT", "AWS4-HMAC-SHA256 Credential=ASWKEY/20131125/us-east-1/kinesis/aws4_request, SignedHeaders=content-type;date;host;user-agent;x-amz-target, Signature=cec25de1e72db69dd48ff4895dc4022e31dc5933209d5bce61286779d49a95e5"},
}

func TestSign(t *testing.T) {
	for _, data := range testSignFactoryData {
		request, err := http.NewRequest("POST", "https://kinesis.us-east-1.amazonaws.com", strings.NewReader("{}"))
		if err != nil {
			t.Errorf("NewRequest Error %v", err)
		}

		request.Header.Set("Content-Type", "application/x-amz-json-1.1")
		request.Header.Set("X-Amz-Target", "")
		request.Header.Set("User-Agent", "Golang Kinesis")

		request.Header.Set("Date", data.DateHeader)
		err = Sign(NewAuth(data.AWS_KEY, data.AWS_SECRET, data.TOKEN), request)
		if err != nil {
			t.Errorf("Error on sign (%v)", err)
			continue
		}
		if request.Header.Get("Authorization") != data.AuthHeader {
			t.Errorf("Get this header (%v), but expect this (%v)", request.Header.Get("Authorization"), data.AuthHeader)
		}
	}
}
