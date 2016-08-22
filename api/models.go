package api

import (
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"time"
)

type RunRequest struct {
	BinUri    string        `json:bin_uri`    // location of binary .tar
	RootfsUri string        `json:rootfs_uri` // location of rootfs .tar
	Timeout   time.Duration `json:timeout`    // cpu time limit
	MemLimit  int64         `json:mem_limit`  // memory limit in bytes
	Command   []string      `json:command`    // startup command
	Stdin     string        `json:stdin`      // stdin
}

type RunResponse struct {
	RuntimeStats *cgroups.Stats `json:runtime_stats`
	// stdout is graded, stderr is not
	Stdout string `json:stdout`
	Stderr string `json:stderr`
}
