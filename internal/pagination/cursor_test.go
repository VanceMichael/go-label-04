package pagination

import (
	"testing"
)

func TestCursorRoundTrip(t *testing.T) {
	raw := Encode(Cursor{CreatedAt: 10, ID: "x"})
	c, err := Decode(raw)
	if err != nil || c.CreatedAt != 10 || c.ID != "x" {
		t.Fatalf("%v %#v", err, c)
	}
	if _, err := Decode("bad"); err == nil {
		t.Fatal("bad cursor accepted")
	}
}
