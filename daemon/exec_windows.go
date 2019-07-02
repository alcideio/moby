package daemon

import (
	"github.com/alcideio/moby/container"
	"github.com/alcideio/moby/daemon/exec"
	"github.com/alcideio/moby/libcontainerd"
)

func execSetPlatformOpt(c *container.Container, ec *exec.Config, p *libcontainerd.Process) error {
	// Process arguments need to be escaped before sending to OCI.
	p.Args = escapeArgs(p.Args)
	p.User.Username = ec.User
	return nil
}
