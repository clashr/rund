package container

import (
	"os"
	"os/exec"
	"path"
	"runtime"
	"syscall"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/opencontainers/runc/libcontainer"
	_ "github.com/opencontainers/runc/libcontainer/nsenter"
)

func init() {
	if len(os.Args) > 1 && os.Args[1] == "init" {
		runtime.GOMAXPROCS(1)
		runtime.LockOSThread()
		factory, _ := libcontainer.New("")
		if err := factory.StartInitialization(); err != nil {
			logrus.Fatalf("init: %s\n", err)
		}
		panic("init: container init failed to exec")
	}
}

var (
	wd     string
	rootfs string
)

func TestInit(t *testing.T) {
	wd, _ = os.Getwd()
	rootfs = path.Join(wd, "../_integration/rootfs")
	Init(path.Join(wd, "../_containers"))
}

func TestCreateConfig(t *testing.T) {
	config := CreateConfig("random-hostname")
	if config.Rootfs != "" {
		t.Error("rootfs not nil on empty on new config")
	}
}

func TestSetNamespaces(t *testing.T) {
	config := CreateConfig("random-hostname")
	if err := SetNamespaces(config); err != nil {
		t.Error(err)
	}
}

func TestSetMemoryCap(t *testing.T) {
	config := CreateConfig("random-hostname")
	SetMemoryCap(config, 1024*1024)
	// TODO: run a container
}

func TestSetTimeout(t *testing.T) {
	config := CreateConfig("random-hostname")
	SetMemoryTimeout(config, 60*time.Second)
	// TODO: run a container
}

func TestBadExit(t *testing.T) {
	container, err := CreateDefaultContainer(rootfs)
	if err != nil {
		t.Error(err)
	}
	defer container.Destroy()

	process := CreateProcess("/bad_exit")
	if err := container.Start(process); err != nil {
		t.Errorf("err in starting: %s\n", err)
	}
	if err := container.Exec(); err != nil {
		t.Errorf("err in exec: %s\n", err)
	}
	_, err = process.Wait()
	if err == nil {
		t.Error("expected exit code 2, got 0.")
	}

	var exitstatus int
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			exitstatus = status.ExitStatus()
			if exitstatus == 2 {
				return
			}
		}
	}
	t.Errorf("expected exit code 2, got %d.", exitstatus)
}

func TestIntegration(t *testing.T) {
	container, err := CreateDefaultContainer(rootfs)
	if err != nil {
		t.Error(err)
	}
	defer container.Destroy()

	process := CreateProcess("/hello")
	if err := container.Start(process); err != nil {
		t.Errorf("err in starting: %s\n", err)
	}
	if err := container.Exec(); err != nil {
		t.Errorf("err in exec: %s\n", err)
	}
	_, err = process.Wait()
	if err != nil {
		t.Fatal(err)
	}
}
