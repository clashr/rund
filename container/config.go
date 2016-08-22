package container

import (
	"errors"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/opencontainers/runc/libcontainer/configs"
	"github.com/opencontainers/runc/libcontainer/utils"
)

func CreateConfig(hostname string) *configs.Config {
	var config = &configs.Config{
		Capabilities: []string{"CAP_DAC_OVERRIDE"},
		Cgroups: &configs.Cgroup{
			Name:   hostname,
			Parent: "system",
			Resources: &configs.Resources{
				MemorySwappiness: nil,
				AllowAllDevices:  nil,
				AllowedDevices:   configs.DefaultAllowedDevices,
			},
		},
		MaskPaths:     []string{"/proc/kcore"},
		ReadonlyPaths: []string{"/proc/sys", "/proc/sysrq-trigger", "/proc/irq", "/proc/bus"},
		Devices:       configs.DefaultAutoCreatedDevices,
		Hostname:      hostname,
		Mounts: []*configs.Mount{
			{
				Source:      "proc",
				Destination: "/proc",
				Device:      "proc",
				Flags:       defaultMountFlags,
			},
			{
				Source:      "tmpfs",
				Destination: "/dev",
				Device:      "tmpfs",
				Flags:       syscall.MS_NOSUID | syscall.MS_STRICTATIME,
				Data:        "mode=755",
			},
			{
				Source:      "devpts",
				Destination: "/dev/pts",
				Device:      "devpts",
				Flags:       syscall.MS_NOSUID | syscall.MS_NOEXEC,
				Data:        "newinstance,ptmxmode=0666,mode=0620,gid=5",
			},
			{
				Device:      "tmpfs",
				Source:      "shm",
				Destination: "/dev/shm",
				Data:        "mode=1777,size=65536k",
				Flags:       defaultMountFlags,
			},
			{
				Source:      "mqueue",
				Destination: "/dev/mqueue",
				Device:      "mqueue",
				Flags:       defaultMountFlags,
			},
			{
				Source:      "sysfs",
				Destination: "/sys",
				Device:      "sysfs",
				Flags:       defaultMountFlags | syscall.MS_RDONLY,
			},
		},
		Networks: []*configs.Network{
			{
				Type:    "loopback",
				Address: "127.0.0.1/0",
				Gateway: "localhost",
			},
		},
	}
	return config
}

func SetNamespaces(c *configs.Config) error {
	for _, ns := range configs.NamespaceTypes() {
		if !configs.IsNamespaceSupported(ns) {
			logrus.Infof("set-namespaces: %s not supported", configs.NsName(ns))
			c.Namespaces.Remove(ns)
			continue
		}
		c.Namespaces.Add(ns, "")
	}

	if c.Namespaces == nil {
		return errors.New("no namespaces supported")
	}

	return nil
}

func SetDefaultNamespaces(c *configs.Config) {
	c.Namespaces = []configs.Namespace{
		{Type: configs.NEWNS},
		{Type: configs.NEWUTS},
		{Type: configs.NEWIPC},
		{Type: configs.NEWPID},
		{Type: configs.NEWNET},
	}
}

func CreateDefaultConfig(path, hostname string) *configs.Config {
	c := CreateConfig(hostname)
	SetDefaultNamespaces(c)
	SetRootfs(c, path)
	return c
}

// Set Memory Cap in Bytes
func SetMemoryCap(c *configs.Config, chars int64) {
	c.Cgroups.Resources.Memory = chars
	c.Cgroups.Resources.MemorySwap = chars
}

func SetTimeout(c *configs.Config, dur time.Duration) {
	c.Rlimits = append(c.Rlimits, configs.Rlimit{
		Type: syscall.RLIMIT_CPU,
		Hard: uint64(dur.Seconds()),
		Soft: uint64(dur.Seconds()),
	})
}

func SetRootfs(c *configs.Config, path string) error {
	var err error
	c.Rootfs, err = utils.ResolveRootfs(utils.CleanPath(path))
	return err
}
