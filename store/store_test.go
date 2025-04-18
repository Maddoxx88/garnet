package store

import (
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	db := New()
	db.Set("foo", "bar", 0)

	val, ok := db.Get("foo")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if val != "bar" {
		t.Errorf("expected 'bar', got '%s'", val)
	}
}

func TestSetWithTTL(t *testing.T) {
	db := New()
	db.Set("temp", "123", 1) // TTL = 1 second

	time.Sleep(2 * time.Second)

	_, ok := db.Get("temp")
	if ok {
		t.Fatal("expected key to be expired")
	}
}

func TestDel(t *testing.T) {
	db := New()
	db.Set("key", "val", 0)
	deleted := db.Del("key")
	if !deleted {
		t.Error("expected key to be deleted")
	}

	_, ok := db.Get("key")
	if ok {
		t.Error("expected key to be gone")
	}
}

func TestTTLLoopCleansUp(t *testing.T) {
	db := New()
	db.Set("ephemeral", "bye", 1)
	db.StartTTLLoop(500 * time.Millisecond)

	time.Sleep(2 * time.Second)

	_, ok := db.Get("ephemeral")
	if ok {
		t.Fatal("key should have been cleaned up")
	}
}
