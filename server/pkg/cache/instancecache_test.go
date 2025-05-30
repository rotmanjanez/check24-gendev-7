package cache

import (
	"context"
	"testing"
	"time"
)

type testValue struct {
	Data string
}

func (v *testValue) MarshalBinary() ([]byte, error) {
	return []byte(v.Data), nil
}

func (v *testValue) UnmarshalBinary(data []byte) error {
	v.Data = string(data)
	return nil
}

func TestInstanceCache_SetGetDelete(t *testing.T) {
	cache := NewInstanceCache("test")
	ctx := context.Background()
	key := "foo"
	val := &testValue{Data: "bar"}

	// Set
	err := cache.Set(ctx, key, val, 0)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get
	got := &testValue{}
	found, err := cache.Get(ctx, key, got)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !found {
		t.Fatalf("Expected key to be found")
	}
	if got.Data != val.Data {
		t.Errorf("Expected %q, got %q", val.Data, got.Data)
	}

	// Delete
	err = cache.Delete(ctx, key)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	found, err = cache.Get(ctx, key, got)
	if err != nil {
		t.Fatalf("Get after delete failed: %v", err)
	}
	if found {
		t.Errorf("Expected key to be deleted")
	}
}

func TestInstanceCache_SetIfNotExists(t *testing.T) {
	cache := NewInstanceCache("test")
	ctx := context.Background()
	key := "foo"
	val := &testValue{Data: "bar"}

	// Should set
	set, err := cache.SetIfNotExists(ctx, key, val, 0)
	if err != nil {
		t.Fatalf("SetIfNotExists failed: %v", err)
	}
	if !set {
		t.Errorf("Expected SetIfNotExists to set value")
	}

	// Should not set again
	set, err = cache.SetIfNotExists(ctx, key, &testValue{Data: "baz"}, 0)
	if err != nil {
		t.Fatalf("SetIfNotExists failed: %v", err)
	}
	if set {
		t.Errorf("Expected SetIfNotExists to not set value if exists")
	}

	// Value should still be original
	got := &testValue{}
	found, err := cache.Get(ctx, key, got)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !found || got.Data != val.Data {
		t.Errorf("Expected value %q, got %q", val.Data, got.Data)
	}
}

func TestInstanceCache_TTL(t *testing.T) {
	cache := NewInstanceCache("test")
	ctx := context.Background()
	key := "foo"
	val := &testValue{Data: "bar"}

	// Set with short TTL
	err := cache.Set(ctx, key, val, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got := &testValue{}
	found, err := cache.Get(ctx, key, got)
	if err != nil || !found {
		t.Fatalf("Expected to find key before TTL expires")
	}

	time.Sleep(60 * time.Millisecond)
	found, err = cache.Get(ctx, key, got)
	if err != nil {
		t.Fatalf("Get after TTL failed: %v", err)
	}
	if found {
		t.Errorf("Expected key to expire after TTL")
	}
}

func TestInstanceCache_Persist(t *testing.T) {
	cache := NewInstanceCache("test")
	ctx := context.Background()
	key := "foo"
	err := cache.Persist(ctx, key)
	if err != nil {
		t.Errorf("Persist should be a no-op and not error, got: %v", err)
	}
}
