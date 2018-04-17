package tmplfunc

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// Funcs is a map of functions used in templates.
var Funcs = map[string]interface{}{

	"toLower": func() interface{} { return unicode.ToLower },
	"toTitle": func() interface{} { return unicode.ToTitle },
	"toUpper": func() interface{} { return unicode.ToUpper },

	"map": strings.Map,

	"exclude": exclude,
	"replace": replace,

	"extractMeta": extractMeta,
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

// extractMeta attempts to extract yaml-like meta data from supplied string.
// It returns map of the extracted key-value pairs, putting all extra lines
// into the tag with name supplied in extraTag argument.
func extractMeta(extraTag string, s string) (map[string]string, error) {
	//TODO
	return nil, nil
}
