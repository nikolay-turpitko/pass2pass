package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"pass2pass/model"
	"pass2pass/parse"
	"pass2pass/paths"
	"pass2pass/store"
)

func main() {
	log.Println("pass2pass started")

	var (
		parserType   string
		storeType    string
		pathCleaner  string
		pathFilter   string
		pathReplacer string
	)

	flag.StringVar(&parserType, "parser", "lastpass", "parser type for input data from stdin")
	flag.StringVar(&storeType, "store", "cmd", "store type for output")
	flag.StringVar(&pathCleaner, "path-cleaner", "", "command executed to clean path")
	flag.StringVar(&pathFilter, "path-filter", "", "command executed to filter/exclude path")
	flag.StringVar(&pathReplacer, "path-replacer", "", "command executed to replace path")

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "\nUsage of %s:\n\n", os.Args[0])
		fmt.Fprintf(w, "pass2pass [-parser PARSER_TYPE] [-store STORE_TYPE] [-path-cleaner PATH_CLEANER_CMD] [-path-filter PATH_FILTER_CMD] [-path-replacer PATH_REPLACER_CMD] ARGS_FOR_STORE ...\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(w, "\nExample:\n\ncat passwords.csv | pass2pass -parser lastpass -store cmd './templates/dir-scheme' './cmd-pass.sh'\n")
		fmt.Fprintf(w, `
Currently implemented store types (shown with required arguments):
plain <template-dir> <out-dir> - stores entries in plain-text files within output directory
cmd <template-dir> <command> - invokes command (such as pass) to store entries

<template-dir> is a directory with Go text templates, used to configure out entries layout.

To create destination entry path Name and Grouping fields of the input record are cleaned
with command, provided in -path-cleaner flag, than $name and $group substrings are replaced
with those cleaned values within template path, after that -path-filter and -path-replacer
command are applied to result. This allows to customize path creation rules and
rearrange entries during import.

`)
	}

	flag.Parse()

	paths.Init(pathCleaner, pathFilter, pathReplacer)

	dataCh, errCh1 := parse.ParseAsync(parserType, os.Stdin)
	errCh2 := store.StoreAsync(storeType, dataCh, flag.Args()...)

	log.Println("pass2pass waits for data on stdin")

	var wg sync.WaitGroup
	processErrorsAsync(&wg, errCh1, false)
	processErrorsAsync(&wg, errCh2, true)
	wg.Wait()

	log.Println("pass2pass finished")
}

// processErrorsAsync reads errors from channel and processes them.
//
// Subpackages should not log errors themselves, but should return errors
// to the main package, which handles them (currently just log to stderr).
func processErrorsAsync(
	wg *sync.WaitGroup,
	errs <-chan error,
	printStats bool) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		errlog := log.New(os.Stderr, "ERR  ", 0)
		count := map[string]int{}
		for err := range errs {
			count["ALL"]++
			if err == nil {
				count["SUCCESS"]++
			} else {
				count["ERRORS"]++
				switch err := err.(type) {
				case model.FatalError:
					if err.IsFatal() {
						errlog.Fatalf("fatal error: %v", err)
					}
				case *exec.ExitError:
					errlog.Printf("error: %v\n%s", err, err.Stderr)
				default:
					errlog.Printf("error: %v", err)
				}
			}
		}
		if printStats {
			log.Println("Entries processed:")
			log.Printf("- All: %d", count["ALL"])
			log.Printf("- Successful: %d", count["SUCCESS"])
			log.Printf("- Errors: %d", count["ERRORS"])
		}
	}()
}

func init() {
	// Configure default logger.
	// Default logger can be used by subpackages for information messages.
	// Subpackages shouldn't change configuration of default logger.
	log.SetPrefix("INFO ")
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
}
