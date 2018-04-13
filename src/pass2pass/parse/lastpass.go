package parse

import (
	"io"

	"github.com/gocarina/gocsv"

	"pass2pass/model"
)

func init() {
	parsers["lastpass"] = lastpass
	parsers["lpass"] = lastpass
}

func lastpass(in io.Reader) (<-chan model.Entry, <-chan error) {
	dataCh := make(chan model.Entry, 100)
	proxyCh := make(chan record, 100)
	errCh := make(chan error, 1)
	go func() {
		defer close(dataCh)
		for r := range proxyCh {
			if r.URL == "http://sn" {
				r.URL = ""
			}
			dataCh <- model.Entry{
				URL:      r.URL,
				Username: r.Username,
				Password: r.Password,
				Extra:    r.Extra,
				Name:     r.Name,
				Grouping: r.Grouping,
			}
		}
	}()
	go func() {
		defer close(errCh)
		errCh <- gocsv.UnmarshalToChan(in, proxyCh)
	}()
	return dataCh, errCh
}

type record struct {
	URL      string `csv:"url"`
	Username string `csv:"username"`
	Password string `csv:"password"`
	Extra    string `csv:"extra"`
	Name     string `csv:"name"`
	Grouping string `csv:"grouping"`
}
