package store

import (
	"fmt"
	"log"

	"pass2pass/model"
)

// StoreAsync inits store of desired type and passes input channel to it.
// It returns output (error) channel to caller.
// After call to StoreAsync, store is ready to process input entries and save
// them according to algorithm of particular store.
// storeType must be a registered store type.
func StoreAsync(
	storeType string,
	entries <-chan model.Entry,
	params ...string) <-chan error {
	s, ok := stores[storeType]
	if !ok {
		errCh := make(chan error, 1)
		defer close(errCh)
		errCh <- model.NewFatalError(fmt.Errorf("unknown store: %s", storeType))
		return errCh
	}
	log.Printf("found store function for store type: %s", storeType)
	return s(entries, params...)
}

type storeFunc func(entries <-chan model.Entry, params ...string) <-chan error

var stores = map[string]storeFunc{}
