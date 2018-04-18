package paths

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"pass2pass/model"
	"pass2pass/tmplfunc"
)

// Init stores command lines, used to implement path manipulation tools.
func Init(fieldCleanerName, fieldCleanerGroup, pathFilter, pathReplacer string) error {
	switch {
	case fieldCleanerName == "":
		CleanerName = noopCleaner{true}
	case isExecutable(fieldCleanerName):
		CleanerName = cmdCleaner{true, fieldCleanerName}
	default:
		t, err := template.New(filepath.Base(fieldCleanerName)).
			Funcs(tmplfunc.Funcs).
			ParseFiles(fieldCleanerName)
		if err != nil {
			return err
		}
		CleanerName = templateCleaner{t}
	}
	switch {
	case fieldCleanerGroup == "":
		CleanerGroup = noopCleaner{false}
	case isExecutable(fieldCleanerGroup):
		CleanerGroup = cmdCleaner{false, fieldCleanerGroup}
	default:
		t, err := template.New(filepath.Base(fieldCleanerGroup)).
			Funcs(tmplfunc.Funcs).
			ParseFiles(fieldCleanerGroup)
		if err != nil {
			return err
		}
		CleanerGroup = templateCleaner{t}
	}
	switch {
	case pathFilter == "":
		Filter = noopFilter{}
	case isExecutable(pathFilter):
		Filter = cmdFilter{pathFilter}
	default:
		t, err := template.New(filepath.Base(pathFilter)).
			Funcs(tmplfunc.Funcs).
			ParseFiles(pathFilter)
		if err != nil {
			return err
		}
		Filter = templateFilter{t}
	}
	switch {
	case pathReplacer == "":
		Replacer = noopReplacer{}
	case isExecutable(pathReplacer):
		Replacer = cmdReplacer{pathReplacer}
	default:
		t, err := template.New(filepath.Base(pathReplacer)).
			Funcs(tmplfunc.Funcs).
			ParseFiles(pathReplacer)
		if err != nil {
			return err
		}
		Replacer = templateReplacer{t}
	}
	return nil
}

// PathFromFields constructs path of target entry from its fields and the path
// of the template. It can invoke external command (if provided) to clean
// fields before using them to construct path.
func PathFromFields(templPath string, entry model.Entry) (string, error) {
	name, err := CleanerName.Do(entry)
	if err != nil {
		return "", err
	}
	group, err := CleanerGroup.Do(entry)
	if err != nil {
		return "", err
	}
	templReplacer := strings.NewReplacer(
		"$name", name,
		"$group", group)
	return filepath.Clean(templReplacer.Replace(templPath)), nil
}

type cleaner interface {
	Do(entry model.Entry) (string, error)
}

type filter interface {
	Do(path string) (bool, error)
}

type replacer interface {
	Do(path string) (string, bool, error)
}

var (
	CleanerName  cleaner
	CleanerGroup cleaner
	Filter       filter
	Replacer     replacer
)

type noopCleaner struct {
	name bool
}

func (c noopCleaner) Do(entry model.Entry) (string, error) {
	if c.name {
		return entry.Name, nil
	}
	return entry.Grouping, nil
}

type cmdCleaner struct {
	name    bool
	cmdLine string
}

func (c cmdCleaner) Do(entry model.Entry) (string, error) {
	cmd := exec.Command(c.cmdLine)
	if c.name {
		cmd.Stdin = strings.NewReader(entry.Name)
	} else {
		cmd.Stdin = strings.NewReader(entry.Grouping)
	}
	out, err := cmd.Output()
	return string(out), err
}

type templateCleaner struct {
	tmpl *template.Template
}

func (c templateCleaner) Do(entry model.Entry) (string, error) {
	var b bytes.Buffer
	err := c.tmpl.Execute(&b, entry)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

type noopFilter struct{}

func (noopFilter) Do(path string) (bool, error) {
	return false, nil
}

type cmdFilter struct {
	cmdLine string
}

func (f cmdFilter) Do(path string) (bool, error) {
	cmd := exec.Command(f.cmdLine)
	cmd.Stdin = strings.NewReader(path)
	out, err := cmd.Output()
	return len(out) == 0, err
}

type templateFilter struct {
	tmpl *template.Template
}

func (f templateFilter) Do(path string) (bool, error) {
	var b bytes.Buffer
	err := f.tmpl.Execute(&b, path)
	if err != nil {
		return false, err
	}
	return b.Len() == 0, nil
}

type noopReplacer struct{}

func (noopReplacer) Do(path string) (string, bool, error) {
	return path, false, nil
}

type cmdReplacer struct {
	cmdLine string
}

func (r cmdReplacer) Do(path string) (string, bool, error) {
	cmd := exec.Command(r.cmdLine)
	cmd.Stdin = strings.NewReader(path)
	out, err := cmd.Output()
	new := string(out)
	return new, new != path, err
}

type templateReplacer struct {
	tmpl *template.Template
}

func (r templateReplacer) Do(path string) (string, bool, error) {
	var b bytes.Buffer
	err := r.tmpl.Execute(&b, path)
	if err != nil {
		return "", false, err
	}
	new := b.String()
	return new, new != path, nil
}

func isExecutable(s string) bool {
	stat, err := os.Stat(s)
	if os.IsNotExist(err) {
		return true // probably, it's a shell command
	}
	_, err = exec.LookPath(s)
	if err == nil {
		return true // executable file
	}
	// if it's a regular file or a symlink, try it as template
	return !stat.Mode().IsRegular() && stat.Mode()&os.ModeSymlink == 0
}
