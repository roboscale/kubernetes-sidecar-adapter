package container

import (
	"testing"

	"github.com/roboscale/kubernetes-sidecar-adapter/pkg/step"
	"github.com/stretchr/testify/assert"
)

func TestContainerCreation(t *testing.T) {

	cont, err := New(1, "main_x")
	expectedCont := Container{
		Type:  "main",
		Name:  "x",
		Pid:   1,
		Path:  "/proc/1/root",
		Steps: nil,
	}

	assert.Equal(t, cont, expectedCont)
	assert.Nil(t, err)

}

func TestContainerConfigure(t *testing.T) {

	cont, _ := New(1, "main_x")

	out, err := cont.Configure()
	assert.Equal(t, out, []string{})
	assert.Nil(t, err)

	steps := []step.Step{
		{
			Name:         "step",
			Command:      "ls",
			Path:         "/",
			IsPathInside: true,
		},
	}
	cont.Steps = &steps

	out, err = cont.Configure()

	assert.Equal(t, len(out), 1)
	assert.Nil(t, err)

	cont.Steps = &[]step.Step{
		{
			Name:         "step",
			Command:      "ls",
			Path:         "/",
			IsPathInside: false,
		},
	}

	_, err = cont.Configure()

	assert.NotNil(t, err)
}
