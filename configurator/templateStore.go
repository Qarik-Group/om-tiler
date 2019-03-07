package configurator

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

type templateStore struct {
	base  string
	store http.FileSystem
}

func (ts *templateStore) lookup(dir string, file string) (io.Reader, error) {
	if filepath.Ext(file) == "" {
		file = fmt.Sprintf("%s.yml", file)
	}
	return ts.store.Open(filepath.Join(ts.base, dir, file))
}

func (ts *templateStore) batchLookup(dir string, files []string, out *[]io.Reader, ignoreMissing bool) error {
	for _, file := range files {
		f, err := ts.lookup(dir, file)
		if err != nil {
			if ignoreMissing {
				continue
			} else {
				return err
			}
		}
		*out = append(*out, f)
	}
	return nil
}
