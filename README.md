# kubernetes-sidecar-adapter
[![codecov](https://codecov.io/gh/roboscale/kubernetes-sidecar-adapter/branch/main/graph/badge.svg?token=N0KX7K5CGW)](https://codecov.io/gh/roboscale/kubernetes-sidecar-adapter)

sidecar adapter that keeps sidecars synced with main container in a lightweight approach.

## Usage
When building an image, add a following script to Dockerfile.

### Main Container (ROS 2 Galactic)
```bash
#!/bin/bash
# adapter_main_ros
sleep infinity
```

```dockerfile
FROM ros:galactic-ros-base-focal
SHELL [ "/bin/bash", "-c" ]
COPY ./adapter_main_ros.sh .
RUN chmod +x ./adapter_main_ros.sh
```

### Sidecar (Ubuntu)
```bash
#!/bin/bash
# adapter_sidecar_ubuntu
sleep infinity
```

```dockerfile
FROM ubuntu:focal
SHELL [ "/bin/bash", "-c" ]
COPY ./adapter_sidecar_ubuntu.sh .
RUN chmod +x ./adapter_sidecar_ubuntu.sh
```

### Kubernetes
After building both, create a pod in your Kubernetes cluster with field `sharedProcessNamespace` as `true`.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: adapter-test
spec:
  shareProcessNamespace: true
  containers:
  - name: ros
    image: <ros-with-adapter-script>
    command: ["./adapter_main_ros.sh"]

  - name: sc1
    image: <ubuntu-with-adapter-script>
    command: ["./adapter_sidecar_ubuntu.sh"]

  - name: sc2
    image: <ubuntu-with-adapter-script>
    command: ["./adapter_sidecar_ubuntu.sh"]

  - name: sc3
    image: <ubuntu-with-adapter-script>
    command: ["./adapter_sidecar_ubuntu.sh"]

  - name: adapter
    image: <this-projects-image>
    command: ["./adapter"]

  restartPolicy: Never
 ```
 
 If you want to see the logs, start `adapter` container with `sleep infinity` command and exec into it with `kubectl exec -it pod/adapter-test -n <namespace> -- bash` and run `/adapter`.
