package container

import (
	"path"

	"github.com/Sirupsen/logrus"

	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runc/libcontainer/utils"

	"github.com/clashr/go-fileutils"
)

func CreateDefaultContainer(dir string) (libcontainer.Container, error) {
	hostname, _ := utils.GenerateRandomName("rund", 10)
	rootfs := path.Join(dir, hostname)

	if err := fileutils.CreateIfNotExists(rootfs, true); err != nil {
		return nil, err
	}

	config := CreateDefaultConfig(rootfs, hostname)

	container, err := factory.Create(hostname, config)
	if err != nil {
		return nil, err
	}
	logrus.Infof("create-default-container: container '%s' created", hostname)

	return container, nil
}
