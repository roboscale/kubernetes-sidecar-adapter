package container

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/roboscale/kubernetes-sidecar-adapter/pkg/step"
)

type ContainerActions interface {
	Configure()
}

type Container struct {
	Type    string
	Name    string
	Pid     int
	Path    string
	Adapter string
	Steps   []step.Step
}

func New(pid int, adapter string) (Container, error) {
	// ./adapter_main_ros.sh

	c := Container{}
	c.Pid = pid
	c.Adapter = adapter

	if pid == 0 || !strings.Contains(adapter, "adapter") {
		return Container{}, errors.New("container needs pid and adapter")
	}

	plain := adapter[2 : len(adapter)-3]
	parts := strings.Split(plain, "_")

	c.Type = parts[1]
	c.Name = parts[2]
	c.Path = "/proc/" + strconv.Itoa(c.Pid) + "/root"

	return c, nil
}

func (c *Container) Configure() ([]string, error) {
	containerPathPlaceholder := ":::container:path:::"

	aggOutput := []string{}
	for key, step := range c.Steps {
		if strings.Contains(step.Command, containerPathPlaceholder) {
			c.Steps[key].Command = strings.Replace(c.Steps[key].Command, containerPathPlaceholder, c.Path, -1)
		}
		if !step.IsPathInside {
			c.Steps[key].Path = c.Path + step.Path
		}

		log.Println("Executing in container -> " + c.Name)
		log.Println("\t" + c.Steps[key].Command)
		out, err := c.Steps[key].Execute()
		if err != nil {
			aggOutput = append(aggOutput, err.Error())
			return aggOutput, err
		}
		aggOutput = append(aggOutput, out)
	}

	return aggOutput, nil
}
