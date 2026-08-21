package metrics

import (
	"testing"
)

func TestCounterSnapshot(t *testing.T) {
	c := NewCounter()
	c.Inc("a")
	c.Add("a", 2)
	if c.Get("a") != 3 {
		t.Fatal(c.Get("a"))
	}
	m := c.Snapshot()
	m["a"] = 99
	if c.Get("a") != 3 {
		t.Fatal("snapshot mutated counter")
	}
}

func TestUnknownCounterStartsAtZero(t *testing.T) {
	c := NewCounter()
	if c.Get("missing") != 0 {
		t.Fatal("missing counter was nonzero")
	}
}

func TestTimerProducesDuration(t *testing.T) {
	timer := StartTimer()
	if timer.Elapsed() < 0 {
		t.Fatal("negative duration")
	}
}
