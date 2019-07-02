package daemon

import (
	// Importing packages here only to make sure their init gets called and
	// therefore they register themselves to the logdriver factory.
	_ "github.com/alcideio/moby/daemon/logger/awslogs"
	_ "github.com/alcideio/moby/daemon/logger/etwlogs"
	_ "github.com/alcideio/moby/daemon/logger/fluentd"
	_ "github.com/alcideio/moby/daemon/logger/jsonfilelog"
	_ "github.com/alcideio/moby/daemon/logger/logentries"
	_ "github.com/alcideio/moby/daemon/logger/splunk"
	_ "github.com/alcideio/moby/daemon/logger/syslog"
)
