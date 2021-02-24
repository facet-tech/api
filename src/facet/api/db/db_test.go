package db

import (
	"strings"
	"testing"
)

// Sample Test Case
func TestCreateId(t *testing.T) {
	key := "test-key"
	generatedID := CreateRandomId(key)
	prefixStr := strings.Split(generatedID, "~")
	if prefixStr[0] != key || len(prefixStr) != 2 {
		t.Error()
	}
}
