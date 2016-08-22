package container

import (
	"github.com/opencontainers/runc/libcontainer"
)

func CreateProcess(name string, arg ...string) *libcontainer.Process {
	var args []string
	args = append(args, name)
	args = append(args, arg...)
	return &libcontainer.Process{Args: args}
}
