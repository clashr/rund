package handlers

import (
	"bytes"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/clashr/rund/api"
	"github.com/clashr/rund/container"
)

func executeRequest(r api.RunRequest) (*api.RunResponse, error) {
	logrus.Info("running request...")

	logrus.Info("setting up rootfs...")
	rootfsdir := viper.GetString("RootfsDir")
	containerdir := viper.GetString("ContainerDir")

	container.Init(containerdir)
	c, err := container.CreateDefaultContainer(rootfsdir)
	if err != nil {
		return nil, err
	}
	defer func() {
		logrus.Infof("destroying container '%s'...", c.ID())
		if err := c.Destroy(); err != nil {
			logrus.Error("cannot destroy: %s", err)
		}
	}()

	wd := c.Config().Rootfs

	if r.RootfsUri != "" {
		if err := download(r.RootfsUri, wd); err != nil {
			return nil, err
		}
	}
	if r.BinUri != "" {
		if err := download(r.BinUri, wd); err != nil {
			return nil, err
		}
	}

	logrus.Infof("applying resource caps...")
	config := c.Config()
	container.SetTimeout(&config, r.Timeout)
	container.SetMemoryCap(&config, r.MemLimit)
	c.Set(config)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	p := container.CreateProcess(r.Command[0], r.Command[1:]...)
	p.Stdin = bytes.NewBufferString(r.Stdin)
	p.Stdout = stdout
	p.Stderr = stderr

	logrus.Infof("starting process...")
	if err := c.Start(p); err != nil {
		return nil, err
	}

	logrus.Infof("running process...")
	if err := c.Exec(); err != nil {
		return nil, err
	}

	logrus.Infof("waiting on process...")
	if ps, err := p.Wait(); err != nil {
		if ps.Success() {
			return nil, err
		}
	}

	logrus.Info("process finished, retreiving stats...")
	stats, err := c.Stats()
	if err != nil {
		return nil, err
	}

	var result api.RunResponse

	result.RuntimeStats = stats.CgroupStats
	result.Stdout = stdout.String()
	result.Stderr = stderr.String()

	return &result, nil
}
