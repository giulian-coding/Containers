# Containers

This Project show cases the simple implementation of Containers in Golang.


## Requirements
- Go `1.27.1`
- This Code was tested in an Ubuntu 26.04 LTS environment, so this Version is recomended for the best experience. Note: It wont work on Windows, but you can use WSL2 to run it.

### Get Filesystem

```bash
wget https://dl-cdn.alpinelinux.org/alpine/v3.19/releases/x86_64/alpine-minirootfs-3.19.0-x86_64.tar.gz
mkdir -p giutainer-root
tar -xzf alpine-minirootfs-3.19.0-x86_64.tar.gz -C giutainer-root
```

## Build and Run

```bash
go build -o giutainer cmd/main.go
sudo ./giutainer run sh
```

## Cleanup

```bash
sudo rm -rf giutainer-root
sudo rmdir /sys/fs/cgroup/giutainer/ 2>/dev/null

## TODOs
- [ ] Implement a Container Management System
- [ ] Implement a Network Management System
- [ ] Implement a Volume Management System
- [ ] Implement a Connection to my personal Assistance System Ion