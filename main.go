package main

import (
	"errors"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/roboscale/kubernetes-sidecar-adapter/pkg/container"
	"github.com/roboscale/kubernetes-sidecar-adapter/pkg/step"
	"github.com/shirou/gopsutil/v3/process"
)

func main() {

	for {
		main := container.Container{}
		processes, _ := process.Processes()
		sidecarsMap := make(map[string]container.Container)
		for _, v := range processes {
			if v.Pid == 1 {
				continue
			}
			commandStr := "ls /proc/" + strconv.Itoa(int(v.Pid)) + "/root/etc/container"
			command := exec.Command("/bin/bash", "-c", commandStr)
			stdout, err := command.Output()
			outstrRaw := string(stdout)
			if err != nil {
				continue
			}

			containerFlag := strings.ReplaceAll(outstrRaw, "\n", "")
			if _, ok := sidecarsMap[containerFlag]; ok {
				continue
			}

			cont, err := container.New(int(v.Pid), containerFlag)
			if err != nil {
				panic(err)
			}

			if cont.Type == "main" {
				main = cont
			} else if cont.Type == "sidecar" {
				sidecarsMap[containerFlag] = cont
			} else {
				panic(errors.New("container undetected: " + containerFlag))
			}
		}

		if main.Pid == 0 {
			panic(errors.New("no main container"))
		}

		sidecars := []container.Container{}
		for _, v := range sidecarsMap {
			sidecars = append(sidecars, v)
		}

		log.Println("Main Container: " + main.Name + "/tPath:" + main.Path)
		sidecarLogs := ""
		for _, s := range sidecars {
			sidecarLogs += "Sidecar Container: " + s.Name + "\tPath:" + s.Path + "\n"
		}
		log.Println(sidecarLogs)

		// fmt.Printf("%+v\n", main)
		// fmt.Printf("%+v\n", sidecars)

		step1 := step.Step{
			Name:    "python_link",
			Command: "ln -sf " + main.Path + "/usr/bin/python3 usr/bin/python3",
		}

		step2 := step.Step{
			Name:    "ros_opt_link",
			Command: "ln -sf " + main.Path + "/opt/ros opt/",
		}

		step3 := step.Step{
			Name:         "ros_opt_link",
			Command:      "./traverser " + main.Path + "/usr/lib " + main.Path + "/usr/include && cat /libs.conf > :::container:path:::/etc/ld.so.conf.d/randomLibs.conf && chroot :::container:path::: /sbin/ldconfig",
			IsPathInside: true,
			Path:         "/",
		}

		steps := []step.Step{step1, step2, step3}

		for sc := range sidecars {
			sidecars[sc].Steps = &steps
			out, err := sidecars[sc].Configure()
			if err != nil {
				log.Println(out)
				panic(err)
			}
			log.Println(out)
		}

		time.Sleep(time.Minute * 1)
	}

}
