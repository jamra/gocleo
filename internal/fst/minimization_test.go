package fst

import (
	"testing"
)

func TestMinimizingBuilder_Basic(t *testing.T) {
	builder := NewMinimizingBuilder()
	
	// Add some test data
	testData := []struct {
		key   string
		value uint64
	}{
		{"apple", 1},
		{"apply", 2},
		{"banana", 3},
	}
	
	for _, item := range testData {
		err := builder.Add([]byte(item.key), item.value)
		if err != nil {
			t.Fatalf("Failed to add key %s: %v", item.key, err)
		}
	}
	
	fst, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build FST: %v", err)
	}
	
	// Test lookups
	for _, item := range testData {
		value, found := fst.Get([]byte(item.key))
		if !found {
			t.Errorf("Key %s not found", item.key)
		}
		if value != item.value {
			t.Errorf("Key %s: expected value %d, got %d", item.key, item.value, value)
		}
	}
	
	// Test non-existent key
	if fst.Contains([]byte("app")) {
		t.Errorf("Key 'app' should not exist")
	}
	
	t.Logf("FST built with %d states", fst.NumStates())
}

func TestMinimizingBuilder_EmptyKey(t *testing.T) {
	builder := NewMinimizingBuilder()
	
	// Empty keys should be rejected
	err := builder.Add([]byte(""), 1)
	if err == nil {
		t.Fatalf("Expected error for empty key")
	}
}

func TestMinimizingBuilder_OrderValidation(t *testing.T) {
	builder := NewMinimizingBuilder()
	
	// Add keys in correct order
	err := builder.Add([]byte("a"), 1)
	if err != nil {
		t.Fatalf("Failed to add 'a': %v", err)
	}
	
	err = builder.Add([]byte("b"), 2)
	if err != nil {
		t.Fatalf("Failed to add 'b': %v", err)
	}
	
	// Try to add out of order - should fail
	err = builder.Add([]byte("a"), 3)
	if err == nil {
		t.Fatalf("Expected error when adding duplicate key")
	}
}
