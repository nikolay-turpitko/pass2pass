package store

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func init() {
	stores["cmd"] = newTemplateBasedStore(&cmdStore{})
}

type cmdStore struct {
	cmd string
}

func (s *cmdStore) init(params ...string) error {
	if len(params) != 1 {
		return fmt.Errorf("illegal arguments: command required")
	}
	s.cmd = params[0]
	return nil
}

func (s cmdStore) save(path, content string) error {
	log.Printf("executing command: `%s` with argument: `%s`", s.cmd, path)
	cmd := exec.Command(s.cmd, path)
	cmd.Stdin = strings.NewReader(content)
	_, err := cmd.Output()
	return err
}
