package transformer

import (
	"context"
	"testing"
)

func TestTransform_ValidFilter(t *testing.T) {
	payload := `{"foo": "bar", "baz": "bat"}`
	filter := ".foo"

	result, err := Transform(context.Background(), payload, filter)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := `"bar"`
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}

	filter = ".missing"
	result, err = Transform(context.Background(), payload, filter)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != "null" {
		t.Errorf("expected null result, got %s", result)
	}

	payload = `[{"age":30, "name":"John"}, {"age":25, "name":"Jane"}]`
	filter = ".[]"

	// Act
	result, err = Transform(context.Background(), payload, filter)

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected = `[{"age":30,"name":"John"},{"age":25,"name":"Jane"}]`
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}
