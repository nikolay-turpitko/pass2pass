package parse

import (
	"html"
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
				URL:      html.UnescapeString(r.URL),
				Username: html.UnescapeString(r.Username),
				Password: html.UnescapeString(r.Password),
				Extra:    html.UnescapeString(r.Extra),
				Name:     html.UnescapeString(r.Name),
				Grouping: html.UnescapeString(r.Grouping),
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
