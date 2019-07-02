package image

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/alcideio/moby/pkg/testutil/assert"
	"github.com/opencontainers/go-digest"
)

func defaultFSStoreBackend(t *testing.T) (StoreBackend, func()) {
	tmpdir, err := ioutil.TempDir("", "images-fs-store")
	assert.NilError(t, err)

	fsBackend, err := NewFSStoreBackend(tmpdir)
	assert.NilError(t, err)

	return fsBackend, func() { os.RemoveAll(tmpdir) }
}

func TestFSGetInvalidData(t *testing.T) {
	store, cleanup := defaultFSStoreBackend(t)
	defer cleanup()

	id, err := store.Set([]byte("foobar"))
	assert.NilError(t, err)

	dgst := digest.Digest(id)

	err = ioutil.WriteFile(filepath.Join(store.(*fs).root, contentDirName, string(dgst.Algorithm()), dgst.Hex()), []byte("foobar2"), 0600)
	assert.NilError(t, err)

	_, err = store.Get(id)
	assert.Error(t, err, "failed to verify")
}

func TestFSInvalidSet(t *testing.T) {
	store, cleanup := defaultFSStoreBackend(t)
	defer cleanup()

	id := digest.FromBytes([]byte("foobar"))
	err := os.Mkdir(filepath.Join(store.(*fs).root, contentDirName, string(id.Algorithm()), id.Hex()), 0700)
	assert.NilError(t, err)

	_, err = store.Set([]byte("foobar"))
	assert.Error(t, err, "failed to write digest data")
}

func TestFSInvalidRoot(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "images-fs-store")
	assert.NilError(t, err)
	defer os.RemoveAll(tmpdir)

	tcases := []struct {
		root, invalidFile string
	}{
		{"root", "root"},
		{"root", "root/content"},
		{"root", "root/metadata"},
	}

	for _, tc := range tcases {
		root := filepath.Join(tmpdir, tc.root)
		filePath := filepath.Join(tmpdir, tc.invalidFile)
		err := os.MkdirAll(filepath.Dir(filePath), 0700)
		assert.NilError(t, err)

		f, err := os.Create(filePath)
		assert.NilError(t, err)
		f.Close()

		_, err = NewFSStoreBackend(root)
		assert.Error(t, err, "failed to create storage backend")

		os.RemoveAll(root)
	}

}

func TestFSMetadataGetSet(t *testing.T) {
	store, cleanup := defaultFSStoreBackend(t)
	defer cleanup()

	id, err := store.Set([]byte("foo"))
	assert.NilError(t, err)

	id2, err := store.Set([]byte("bar"))
	assert.NilError(t, err)

	tcases := []struct {
		id    digest.Digest
		key   string
		value []byte
	}{
		{id, "tkey", []byte("tval1")},
		{id, "tkey2", []byte("tval2")},
		{id2, "tkey", []byte("tval3")},
	}

	for _, tc := range tcases {
		err = store.SetMetadata(tc.id, tc.key, tc.value)
		assert.NilError(t, err)

		actual, err := store.GetMetadata(tc.id, tc.key)
		assert.NilError(t, err)

		if bytes.Compare(actual, tc.value) != 0 {
			t.Fatalf("Metadata expected %q, got %q", tc.value, actual)
		}
	}

	_, err = store.GetMetadata(id2, "tkey2")
	assert.Error(t, err, "failed to read metadata")

	id3 := digest.FromBytes([]byte("baz"))
	err = store.SetMetadata(id3, "tkey", []byte("tval"))
	assert.Error(t, err, "failed to get digest")

	_, err = store.GetMetadata(id3, "tkey")
	assert.Error(t, err, "failed to get digest")
}

func TestFSInvalidWalker(t *testing.T) {
	store, cleanup := defaultFSStoreBackend(t)
	defer cleanup()

	fooID, err := store.Set([]byte("foo"))
	assert.NilError(t, err)

	err = ioutil.WriteFile(filepath.Join(store.(*fs).root, contentDirName, "sha256/foobar"), []byte("foobar"), 0600)
	assert.NilError(t, err)

	n := 0
	err = store.Walk(func(id digest.Digest) error {
		assert.Equal(t, id, fooID)
		n++
		return nil
	})
	assert.NilError(t, err)
	assert.Equal(t, n, 1)
}

func TestFSGetSet(t *testing.T) {
	store, cleanup := defaultFSStoreBackend(t)
	defer cleanup()

	type tcase struct {
		input    []byte
		expected digest.Digest
	}
	tcases := []tcase{
		{[]byte("foobar"), digest.Digest("sha256:c3ab8ff13720e8ad9047dd39466b3c8974e592c2fa383d4a3960714caef0c4f2")},
	}

	randomInput := make([]byte, 8*1024)
	_, err := rand.Read(randomInput)
	assert.NilError(t, err)

	// skipping use of digest pkg because it is used by the implementation
	h := sha256.New()
	_, err = h.Write(randomInput)
	assert.NilError(t, err)

	tcases = append(tcases, tcase{
		input:    randomInput,
		expected: digest.Digest("sha256:" + hex.EncodeToString(h.Sum(nil))),
	})

	for _, tc := range tcases {
		id, err := store.Set([]byte(tc.input))
		assert.NilError(t, err)
		assert.Equal(t, id, tc.expected)
	}

	for _, tc := range tcases {
		data, err := store.Get(tc.expected)
		assert.NilError(t, err)
		if bytes.Compare(data, tc.input) != 0 {
			t.Fatalf("expected data %q, got %q", tc.input, data)
		}
	}
}

func TestFSGetUnsetKey(t *testing.T) {
	store, cleanup := defaultFSStoreBackend(t)
	defer cleanup()

	for _, key := range []digest.Digest{"foobar:abc", "sha256:abc", "sha256:c3ab8ff13720e8ad9047dd39466b3c8974e592c2fa383d4a3960714caef0c4f2a"} {
		_, err := store.Get(key)
		assert.Error(t, err, "failed to get digest")
	}
}

func TestFSGetEmptyData(t *testing.T) {
	store, cleanup := defaultFSStoreBackend(t)
	defer cleanup()

	for _, emptyData := range [][]byte{nil, {}} {
		_, err := store.Set(emptyData)
		assert.Error(t, err, "invalid empty data")
	}
}

func TestFSDelete(t *testing.T) {
	store, cleanup := defaultFSStoreBackend(t)
	defer cleanup()

	id, err := store.Set([]byte("foo"))
	assert.NilError(t, err)

	id2, err := store.Set([]byte("bar"))
	assert.NilError(t, err)

	err = store.Delete(id)
	assert.NilError(t, err)

	_, err = store.Get(id)
	assert.Error(t, err, "failed to get digest")

	_, err = store.Get(id2)
	assert.NilError(t, err)

	err = store.Delete(id2)
	assert.NilError(t, err)

	_, err = store.Get(id2)
	assert.Error(t, err, "failed to get digest")
}

func TestFSWalker(t *testing.T) {
	store, cleanup := defaultFSStoreBackend(t)
	defer cleanup()

	id, err := store.Set([]byte("foo"))
	assert.NilError(t, err)

	id2, err := store.Set([]byte("bar"))
	assert.NilError(t, err)

	tcases := make(map[digest.Digest]struct{})
	tcases[id] = struct{}{}
	tcases[id2] = struct{}{}
	n := 0
	err = store.Walk(func(id digest.Digest) error {
		delete(tcases, id)
		n++
		return nil
	})
	assert.NilError(t, err)
	assert.Equal(t, n, 2)
	assert.Equal(t, len(tcases), 0)
}

func TestFSWalkerStopOnError(t *testing.T) {
	store, cleanup := defaultFSStoreBackend(t)
	defer cleanup()

	id, err := store.Set([]byte("foo"))
	assert.NilError(t, err)

	tcases := make(map[digest.Digest]struct{})
	tcases[id] = struct{}{}
	err = store.Walk(func(id digest.Digest) error {
		return errors.New("what")
	})
	assert.Error(t, err, "what")
}
