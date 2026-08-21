package notification

import (
	"context"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"testing"
)

func TestDelivery(t *testing.T) {
	r := New()
	sender := &MemorySender{}
	if err := r.Register("email", sender); err != nil {
		t.Fatal(err)
	}
	if err := r.Deliver(context.Background(), Message{ID: "1", Recipient: "a", Channel: "email", Body: "hello"}); err != nil {
		t.Fatal(err)
	}
	if len(sender.Sent) != 1 {
		t.Fatal(sender.Sent)
	}
	if err := r.Deliver(context.Background(), Message{ID: "2", Recipient: "a", Channel: "sms", Body: "x"}); err != domain.ErrNotFound {
		t.Fatal(err)
	}
}
