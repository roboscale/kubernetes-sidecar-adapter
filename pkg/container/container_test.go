package container

import (
	"errors"
	"testing"

	"github.com/roboscale/kubernetes-sidecar-adapter/pkg/step"
	"github.com/stretchr/testify/assert"
)

func TestContainerCreation(t *testing.T) {

	cont, err := New(1, "./adapter_main_x.sh")
	expectedCont := Container{
		Type:    "main",
		Name:    "x",
		Pid:     1,
		Path:    "/proc/1/root",
		Adapter: "./adapter_main_x.sh",
		Steps:   nil,
	}

	assert.Equal(t, cont, expectedCont)
	assert.Nil(t, err)

	cont, err = New(0, "./adapter_sidecar_x.sh")

	expectedErr := errors.New("container needs pid and adapter")
	expectedCont = Container{}

	assert.Equal(t, cont, expectedCont, "they should be equal")
	assert.Equal(t, err, expectedErr, "they should be equal")

}

func TestContainerConfigure(t *testing.T) {

	cont, _ := New(1, "./adapter_main_x.sh")

	out, err := cont.Configure()
	assert.Equal(t, out, []string{})
	assert.Nil(t, err)

	cont.Steps = &[]step.Step{
		{
			Name:         "step",
			Command:      "ls",
			Path:         "/",
			IsPathInside: true,
		},
	}

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
