package builder

import (
	"os"

	"github.com/alcideio/moby/pkg/archive"
	"github.com/alcideio/moby/pkg/gitutils"
)

// MakeGitContext returns a Context from gitURL that is cloned in a temporary directory.
func MakeGitContext(gitURL string) (ModifiableContext, error) {
	root, err := gitutils.Clone(gitURL)
	if err != nil {
		return nil, err
	}

	c, err := archive.Tar(root, archive.Uncompressed)
	if err != nil {
		return nil, err
	}

	defer func() {
		// TODO: print errors?
		c.Close()
		os.RemoveAll(root)
	}()
	return MakeTarSumContext(c)
}
