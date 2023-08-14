# Go Simple Container

A minimal container implementation using Go, leveraging Linux namespaces and cgroups for process isolation. This project was created to gain a better understanding of containers and the Go programming language.

## Requirements

- Linux OS with kernel support for namespaces and cgroups(v2).
- Go
- Root or `sudo` access may be required for certain operations, especially when working with cgroups.

## Installation

1. Clone the repository:
```bash
git clone https://github.com/voukatas/go_container
cd go_container
```

2. Download a mini Linux filesystem. For example:
```bash
wget https://dl-cdn.alpinelinux.org/alpine/v3.18/releases/x86_64/alpine-minirootfs-3.18.3-x86_64.tar.gz
tar -xzvf alpine-minirootfs-3.18.3-x86_64.tar.gz -C ./alpine_fs
```

3. Build the project:
```bash
go build -o container
```

## Usage
Execute a command inside the container:

```bash
sudo ./container run <command_name>
```
For example:
```bash
sudo ./container run /bin/sh
```
## Limitations

- This is a basic and educational example of a container and lacks features found in production-ready container solutions like Docker.
- Proper cleanup and handling of cgroups and namespaces are required to avoid system issues.
- Ensure you understand the security implications before using this in a production environment.

## Learning Resources

1. [Deep Into Container: Build Your Own Container with Golang](https://faun.pub/deep-into-container-build-your-own-container-with-golang-98ef93f42923)
2. [YouTube Video: Understanding Containers](https://www.youtube.com/watch?v=8fi7uSYlOdc)
