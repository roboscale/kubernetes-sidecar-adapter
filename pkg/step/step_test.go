package step

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStepExecution(t *testing.T) {

	step := Step{
		Name:         "test",
		Command:      "ls",
		Path:         ".",
		IsPathInside: true,
	}
	out, err := step.Execute()
	assert.NotEqual(t, out, "")
	assert.Nil(t, err)

}
