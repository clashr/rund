package main

import (
	"net/http"
	"os"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/opencontainers/runc/libcontainer"
	_ "github.com/opencontainers/runc/libcontainer/nsenter"

	"github.com/clashr/rund/handlers"
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

func main() {
	viper.SetDefault("Addr", ":3000")
	viper.SetDefault("RootfsDir", "/var/lib/rund/fs")
	viper.SetDefault("ContainerDir", "/var/lib/rund/container")
	viper.SetDefault("MaxSize", 1024*1024)

	logrus.Info("starting rund...")
	viper.SetConfigName("rund")
	viper.AddConfigPath("/etc/rund")
	viper.AddConfigPath(".")
	logrus.Info("searching for config...")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("failed to read config: %s\n", err)
	}
	viper.SetEnvPrefix("rund")
	viper.AutomaticEnv()

	addr := viper.GetString("Addr")

	logrus.Info("preparing routes...")
	http.HandleFunc("/run", handlers.MainController)

	logrus.Infof("serving at addr: %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logrus.Fatalf("failed to start server: %s", err)
	}
}
