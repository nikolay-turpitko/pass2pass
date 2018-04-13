package store

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"pass2pass/model"
	"pass2pass/paths"
)

// This file contains primitives for store based on templates.

// MaxConcurrentEntries limits maximum concurrently processed entries.
const MaxConcurrentEntries = 20

// newTemplateBasedStore creates store, based on templates, with provided
// callback function, which is used to actually store entry.
func newTemplateBasedStore(store templateBasedStore) storeFunc {
	return func(entries <-chan model.Entry, params ...string) <-chan error {
		errCh := make(chan error, 1)
		go func() {
			defer close(errCh)
			if len(params) < 1 {
				errCh <- model.NewFatalError(fmt.Errorf("illegal arguments: template dir required"))
				return
			}
			templateDir := params[0]
			err := store.init(params[1:]...)
			if err != nil {
				errCh <- model.NewFatalError(err)
				return
			}
			templates, err := prepareTemplates(templateDir)
			if err != nil {
				errCh <- model.NewFatalError(err)
				return
			}
			// http://jmoiron.net/blog/limiting-concurrency-in-go/
			sem := make(chan struct{}, MaxConcurrentEntries)
			for entry := range entries {
				entry := entry
				sem <- struct{}{}
				go func() {
					defer func() { <-sem }()
					errCh <- processEntry(templates, store, entry)
				}()
			}
			for i := 0; i < cap(sem); i++ {
				sem <- struct{}{}
			}
		}()
		return errCh
	}
}

// prepareTemplates searches templateDir for template files, parses them and
// returns slice of parsed templates.
func prepareTemplates(templateDir string) (map[string]*template.Template, error) {
	log.Printf("processing templates, folder: %s", templateDir)
	templates := map[string]*template.Template{}
	err := filepath.Walk(
		templateDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				log.Printf("processing template, file: %s", path)
				tp, err := filepath.Rel(templateDir, path)
				if err != nil {
					return err
				}
				t, err := template.ParseFiles(path)
				if err != nil {
					return err
				}
				templates[tp] = t
			}
			return nil
		})
	log.Printf("processed templates, number of templates: %d, error: %v", len(templates), err)
	return templates, err
}

// processEntry handles single password manager's entry.
// It invokes storeCallback to persist string, produced by applying template to
// entry data.
func processEntry(
	templates map[string]*template.Template,
	store templateBasedStore,
	entry model.Entry) error {
	for tp, t := range templates {
		p, err := paths.PathFromFields(tp, entry)
		if err != nil {
			return err
		}
		exclude, err := paths.Exclude(p)
		if err != nil {
			return err
		}
		if exclude {
			log.Printf("exclude path: %s", p)
			continue
		}
		log.Printf("processing path: %s", p)
		p2, changed, err := paths.Transform(p)
		if err != nil {
			return err
		}
		if changed {
			log.Printf("path: %s replaced with %s", p, p2)
			p = p2
		}
		var b bytes.Buffer
		err = t.Execute(&b, entry)
		if err != nil {
			return err
		}
		outStr := strings.TrimSpace(b.String())
		if outStr != "" {
			if err := store.save(p, outStr); err != nil {
				return err
			}
		}
		log.Printf("processed path: %s", p)
	}
	return nil
}

// templateBasedStore is an interface of template-based store.
type templateBasedStore interface {
	init(params ...string) error
	save(path, content string) error
}
