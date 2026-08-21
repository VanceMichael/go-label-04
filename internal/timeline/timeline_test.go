package timeline

import (
	"testing"
	"time"
)

func TestTimeline(t *testing.T) {
	base := time.Now().UTC()
	tl := New()
	if err := tl.Append(Event{Kind: "loaded", Actor: "a", At: base.Add(time.Hour)}); err != nil {
		t.Fatal(err)
	}
	if err := tl.Append(Event{Kind: "departed", Actor: "a", At: base.Add(2 * time.Hour)}); err != nil {
		t.Fatal(err)
	}
	d, err := tl.Duration("loaded", "departed")
	if err != nil || d != time.Hour {
		t.Fatalf("%v %v", err, d)
	}
	if _, err := tl.Latest("missing"); err == nil {
		t.Fatal("missing event")
	}
}
