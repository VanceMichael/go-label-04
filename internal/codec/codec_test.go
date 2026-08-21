package codec

import (
	"testing"
)

func TestCodecRoundTrip(t *testing.T) {
	input := map[string]any{"name": "x", "n": 2}
	raw, err := Encode(input)
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := Decode(raw, &out); err != nil || out["name"] != "x" {
		t.Fatalf("%v %#v", err, out)
	}
	compressed, err := Compress(raw)
	if err != nil || len(compressed) == 0 {
		t.Fatal(err)
	}
}
