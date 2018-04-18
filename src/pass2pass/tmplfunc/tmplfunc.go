package tmplfunc

import (
	"bufio"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"github.com/fatih/structs"
	"gopkg.in/yaml.v2"
)

// Funcs is a map of functions used in templates.
var Funcs = map[string]interface{}{

	"toLower": func() interface{} { return unicode.ToLower },
	"toTitle": func() interface{} { return unicode.ToTitle },
	"toUpper": func() interface{} { return unicode.ToUpper },

	"map": strings.Map,

	"exclude": exclude,
	"replace": replace,

	"excludeKey":  excludeKey,
	"extractMeta": extractMeta,
	"mapKeys":     mapKeys,
	"toMap":       structs.Map,
	"yaml":        renderYAML,

	"hostname": hostname,
}

// exclude is kind of like `grep -v` - it returns its last argument if it does
// not match any others its arguments. All arguments except last are regular
// expressions, last argument is a tested value.
func exclude(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("exclude: at least 2 args required")
	}
	s := args[len(args)-1]
	ptrn := strings.Join(args[:len(args)-1], "|")
	r, err := regexp.Compile(ptrn)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		s := scanner.Text()
		if !r.MatchString(s) {
			b.WriteString(s)
			b.WriteRune('\n')
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	s = b.String()
	if s == "" {
		return "", nil
	}
	return s[:len(s)-1], nil
}

// replace is kind of like `sed` - it returns its last argument with substrings
// replaced by its others arguments. All arguments except last are pairs of
// regular expressions and replacement strings. Last argument is an input
// string to work on.
func replace(args ...string) (string, error) {
	if len(args)%2 != 1 {
		return "", fmt.Errorf("replace: even number of args (%d)", len(args))
	}
	s := args[len(args)-1]
	ptrns := args[:len(args)-1]
	for i := 0; i < len(ptrns); i += 2 {
		r, err := regexp.Compile(ptrns[i])
		if err != nil {
			return "", err
		}
		s = r.ReplaceAllString(s, ptrns[i+1])
	}
	return s, nil
}

// renderYAML converts struct or map to YAML representation.
func renderYAML(s interface{}) (string, error) {
	if s, ok := s.(map[string]interface{}); ok && len(s) == 0 {
		return "", nil
	}
	b, err := yaml.Marshal(s)
	return string(b), err
}

// excludeKey deletes keys from the map. Map accepted as the last argument, all
// other arguments must be strings and treated as excluded keys.
func excludeKey(args ...interface{}) map[string]interface{} {
	s := args[len(args)-1].(map[string]interface{})
	keys := args[:len(args)-1]
	for _, k := range keys {
		delete(s, k.(string))
	}
	return s
}

// mapKeys processes keys of the map with given predicate function.
func mapKeys(f func(rune) rune, m map[string]interface{}) map[string]interface{} {
	n := make(map[string]interface{}, len(m))
	for k, v := range m {
		n[strings.Map(f, k)] = v
	}
	return n
}

// extractMeta attempts to extract yaml-like meta data from supplied string.
// It returns map of the extracted key-value pairs, putting all extra lines
// into the tag with name supplied in extraTag argument.
func extractMeta(extraTag string, s string) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	// try to parse as valid yaml first
	err := yaml.Unmarshal([]byte(s), &m)
	if err == nil {
		return m, nil
	}
	// if fail, try to leave only safe-looking lines and parse again
	safeLooking := `(^\s*-\s*[\p{L}\d_])|(^[\s\p{L}\d_]*[\p{L}\d_]+:\s*[^/])|(^[\s\p{L}\d_]*[\p{L}\d_]+:\s*$)`
	r, err := regexp.Compile(safeLooking)
	if err != nil {
		return nil, err
	}
	quotingRequired := `^([\s\p{L}\d_]*[\p{L}\d_]+:\s*)([^\s"].*)$`
	q, err := regexp.Compile(quotingRequired)
	if err != nil {
		return nil, err
	}
	var b, bExtra strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		s := scanner.Text()
		if r.MatchString(s) {
			s = q.ReplaceAllString(s, `$1"$2"`)
			b.WriteString(strings.TrimSpace(strings.Replace(s, ":", ": ", 1)))
			b.WriteRune('\n')
		} else {
			bExtra.WriteString(s)
			bExtra.WriteRune('\n')
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(b.String()), &m)
	if err == nil {
		m[extraTag] = strings.TrimSpace(bExtra.String())
		return m, nil
	}
	// if fail anyway, put whole input as extra tag
	return map[string]interface{}{extraTag: strings.TrimSpace(s)}, nil
}

// hostname returns host name of URL. It returns original string if it cannot
// parse it as URL.
func hostname(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return s
	}
	if strings.Contains(s, "localhost") {
		return ""
	}
	matched, _ := regexp.MatchString(`\b(?:\d{1,3}\.){3}\d{1,3}\b`, s)
	if matched {
		return ""
	}
	return u.Hostname()
}
