package step

import (
	"os/exec"
)

type StepActions interface {
	Execute()
}

type Step struct {
	Name         string
	Command      string
	Path         string
	IsPathInside bool
}

func (s Step) Execute() (string, error) {
	initialCommand := "/bin/bash"
	arguments := []string{"-c"}
	setDir := "cd " + s.Path + " &&"
	arguments = append(arguments, setDir+s.Command)

	command := exec.Command(initialCommand, arguments...)
	stdout, err := command.Output()
	if err != nil {
		return string(stdout), err
	}

	return string(stdout), nil
}
