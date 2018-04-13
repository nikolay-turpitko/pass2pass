package store

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func init() {
	stores["plain"] = newTemplateBasedStore(&plainStore{})
}

type plainStore struct {
	baseDir string
}

func (s *plainStore) init(params ...string) error {
	if len(params) != 1 {
		return fmt.Errorf("illegal arguments: base dir required")
	}
	s.baseDir = params[0]
	return nil
}

func (s plainStore) save(path, content string) error {
	p := filepath.Join(s.baseDir, path)
	log.Printf("creating file: %s", p)
	err := os.MkdirAll(filepath.Dir(p), 0700)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(p, []byte(content), 0600)
}
