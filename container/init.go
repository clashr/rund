package container

import (
	"github.com/Sirupsen/logrus"
	"github.com/opencontainers/runc/libcontainer"
)

var factory libcontainer.Factory

func Init(path string) error {
	var err error

	if factory != nil {
		return nil
	}

	factory, err = libcontainer.New(path, libcontainer.Cgroupfs)
	if err != nil {
		return err
	}

	logrus.Info("container: container factory initialized")
	return nil
}
