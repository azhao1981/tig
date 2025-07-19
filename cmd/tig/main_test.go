package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	assert.NotEmpty(t, Version)
	assert.NotEmpty(t, Build)
}

func TestRun(t *testing.T) {
	// This is a basic test to ensure the run function doesn't crash
	// In a real scenario, we would use a mock terminal
	// For now, we'll just test that the function exists and returns an error
	// when there's no terminal
	err := run()
	assert.Error(t, err)
}