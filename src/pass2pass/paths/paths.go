package paths

import (
	"os/exec"
	"path/filepath"
	"strings"

	"pass2pass/model"
)

var (
	cleaner, filter, replacer string
)

// Init stores command lines, used to implement path manipulation tools.
func Init(pathCleaner, pathFilter, pathReplacer string) {
	cleaner, filter, replacer = pathCleaner, pathFilter, pathReplacer
}

func clean(s string) (string, error) {
	if cleaner == "" {
		return s, nil
	}
	cmd := exec.Command(cleaner)
	cmd.Stdin = strings.NewReader(s)
	out, err := cmd.Output()
	return string(out), err
}

// PathFromFields constructs path of target entry from its fields and the path
// of the template. It can invoke external command (if provided) to clean
// fields before using them to construct path.
func PathFromFields(templPath string, r model.Entry) (string, error) {
	name, err := clean(r.Name)
	if err != nil {
		return "", err
	}
	group, err := clean(r.Grouping)
	if err != nil {
		return "", err
	}
	templReplacer := strings.NewReplacer(
		"$name", name,
		"$group", group)
	return filepath.Clean(templReplacer.Replace(templPath)), nil
}

// Exclude checks if path should be excluded. It invokes external command (if
// provided) to do its job.
func Exclude(path string) (bool, error) {
	if filter == "" {
		return false, nil
	}
	cmd := exec.Command(filter)
	cmd.Stdin = strings.NewReader(path)
	out, err := cmd.Output()
	return len(out) == 0, err
}

// Transform invokes external command (if provided) to replace substring within
// entry's path.
func Transform(path string) (string, bool, error) {
	if replacer == "" {
		return path, false, nil
	}
	cmd := exec.Command(replacer)
	cmd.Stdin = strings.NewReader(path)
	out, err := cmd.Output()
	new := string(out)
	return new, new != path, err
}
