package main

import (
	"log"
	"strings"
	"time"

	"github.com/roboscale/kubernetes-sidecar-adapter/pkg/container"
	"github.com/roboscale/kubernetes-sidecar-adapter/pkg/step"
	"github.com/shirou/gopsutil/v3/process"
)

func main() {

	for {
		main := container.Container{}
		sidecars := []container.Container{}
		processes, _ := process.Processes()
		for _, v := range processes {

			if command, _ := v.CmdlineSlice(); len(command) > 1 && strings.Contains(command[1], "adapter") {
				cont, err := container.New(int(v.Pid), command[1])
				if err != nil {
					panic(err)
				}

				switch cont.Type {
				case "main":
					main = cont
				case "sidecar":
					sidecars = append(sidecars, cont)
				}

			}
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
