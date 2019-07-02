// +build !exclude_graphdriver_devicemapper,linux

package register

import (
	// register the devmapper graphdriver
	_ "github.com/alcideio/moby/daemon/graphdriver/devmapper"
)
