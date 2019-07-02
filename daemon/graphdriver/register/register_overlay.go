// +build !exclude_graphdriver_overlay,linux

package register

import (
	// register the overlay graphdriver
	_ "github.com/alcideio/moby/daemon/graphdriver/overlay"
	_ "github.com/alcideio/moby/daemon/graphdriver/overlay2"
)
