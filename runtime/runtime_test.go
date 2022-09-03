package runtime

import (
	"testing"
	"strings"
	"github.com/stretchr/testify/assert"
)

func TestGetCallPath(t *testing.T) {
	path := GetCallPath()
	assert.True(t, strings.HasSuffix(path, "runtime_test.go"), "except get true but got false")
}
