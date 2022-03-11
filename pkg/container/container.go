package container

import (
	"log"
	"strconv"
	"strings"

	"github.com/roboscale/kubernetes-sidecar-adapter/pkg/step"
)

type ContainerActions interface {
	Configure()
}

type Container struct {
	Type  string
	Name  string
	Pid   int
	Path  string
	Steps *[]step.Step
}

func New(pid int, flag string) (Container, error) {
	// ./adapter_main_ros.sh
	parts := strings.Split(flag, "_")
	c := Container{}
	c.Pid = pid
	c.Type = parts[0]
	c.Name = parts[1]
	c.Path = "/proc/" + strconv.Itoa(c.Pid) + "/root"

	return c, nil
}

func (c *Container) Configure() ([]string, error) {
	containerPathPlaceholder := ":::container:path:::"

	aggOutput := []string{}
	if c.Steps != nil {
		for _, step := range *c.Steps {
			step.Command = strings.Replace(step.Command, containerPathPlaceholder, c.Path, -1)
			if !step.IsPathInside {
				step.Path = c.Path + step.Path
			}

			log.Println("Executing in container -> " + c.Name)
			log.Println("\t" + step.Command)
			out, err := step.Execute()
			if err != nil {
				aggOutput = append(aggOutput, err.Error())
				return aggOutput, err
			}
			aggOutput = append(aggOutput, out)
		}
	}

	return aggOutput, nil
}

func (c Container) Equals(second Container) bool {
	if c.Name == second.Name &&
		c.Type == second.Type &&
		c.Pid == second.Pid {
		return true
	}

	return false
}

func AllEquals(first []Container, second []Container) bool {
	if len(first) != len(second) {
		return false
	}

	for _, f := range first {

		hasEqual := false
		for _, s := range second {
			if f.Equals(s) {
				hasEqual = true
				break
			}
		}

		if !hasEqual {
			return false
		}
	}

	return true
}
