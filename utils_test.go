package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	key, value := "key", "value"
	os.Setenv(key, value)

	assert.Equal(t, value, getenv(key, ""))
	assert.Equal(t, "fallback", getenv("not_existed_key", "fallback"))
}
