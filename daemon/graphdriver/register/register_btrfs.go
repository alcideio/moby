// +build !exclude_graphdriver_btrfs,linux

package register

import (
	// register the btrfs graphdriver
	_ "github.com/alcideio/moby/daemon/graphdriver/btrfs"
)
