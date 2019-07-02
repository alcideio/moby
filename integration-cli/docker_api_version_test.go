package main

import (
	"encoding/json"
	"net/http"

	"github.com/alcideio/moby/api/types"
	"github.com/alcideio/moby/dockerversion"
	"github.com/alcideio/moby/integration-cli/checker"
	"github.com/alcideio/moby/integration-cli/request"
	"github.com/go-check/check"
)

func (s *DockerSuite) TestGetVersion(c *check.C) {
	status, body, err := request.SockRequest("GET", "/version", nil, daemonHost())
	c.Assert(status, checker.Equals, http.StatusOK)
	c.Assert(err, checker.IsNil)

	var v types.Version

	c.Assert(json.Unmarshal(body, &v), checker.IsNil)

	c.Assert(v.Version, checker.Equals, dockerversion.Version, check.Commentf("Version mismatch"))
}
