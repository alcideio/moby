package stack

import (
	"bytes"
	"testing"

	"github.com/alcideio/moby/cli/compose/convert"
	"github.com/alcideio/moby/cli/internal/test"
	"github.com/alcideio/moby/pkg/testutil/assert"
	"golang.org/x/net/context"
)

func TestPruneServices(t *testing.T) {
	ctx := context.Background()
	namespace := convert.NewNamespace("foo")
	services := map[string]struct{}{
		"new":  {},
		"keep": {},
	}
	client := &fakeClient{services: []string{objectName("foo", "keep"), objectName("foo", "remove")}}
	dockerCli := test.NewFakeCli(client, &bytes.Buffer{})
	dockerCli.SetErr(&bytes.Buffer{})

	pruneServices(ctx, dockerCli, namespace, services)

	assert.DeepEqual(t, client.removedServices, buildObjectIDs([]string{objectName("foo", "remove")}))
}
