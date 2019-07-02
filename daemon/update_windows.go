// +build windows

package daemon

import (
	"github.com/alcideio/moby/api/types/container"
	"github.com/alcideio/moby/libcontainerd"
)

func toContainerdResources(resources container.Resources) libcontainerd.Resources {
	var r libcontainerd.Resources
	return r
}
