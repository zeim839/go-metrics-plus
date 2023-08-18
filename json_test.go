package metrics

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestRegistryMarshallJSON(t *testing.T) {
	b := &bytes.Buffer{}
	enc := json.NewEncoder(b)
	r := NewRegistry()
	r.Register("counter", NewCounter(nil))
	enc.Encode(r)
	if s := b.String(); s != "{\"counter\":{\"count\":0,\"labels\":{}}}\n" {
		t.Fatalf(s)
	}
}

func TestRegistryWriteJSONOnce(t *testing.T) {
	r := NewRegistry()
	r.Register("counter", NewCounter(nil))
	b := &bytes.Buffer{}
	WriteJSONOnce(r, b)
	if s := b.String(); s != "{\"counter\":{\"count\":0,\"labels\":{}}}\n" {
		t.Fail()
	}
}
