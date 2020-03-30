package unipaw

import (
	"log"
	"net/http"

	"github.com/parnurzeal/gorequest"
)

const (
	// INQUERY ...
	INQUERY = iota
	// INDATA ...
	INDATA
	// INHEADER ...
	INHEADER
	// INCOOKIE ...
	INCOOKIE // not use
)

// TODO: remove gorequest
func doRequest(site *Site, path *Path, oldr *http.Request) (rdata interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			log.Printf("request error: %s\n", err)
		}
	}()

	url := PathJoin(site.Host, site.URL, path.URL)

	// set method and url
	r := gorequest.New().CustomMethod(path.Method, url)

	// set query, data, header and cookie if which is neccessary
	query := make(map[string]interface{})

	data := path.PostData()

	// set params to query, data, header and cookies.
	// p should get value from r
	for _, p := range site.Params {
		err = p.Into(p, oldr, query, data)
		if err != nil {
			return
		}
	}
	for _, p := range path.Params {
		err = p.Into(p, oldr, query, data)
		if err != nil {
			return
		}
	}

	// append other query

	// combine them all with r
	path.RequestHook(path, oldr, query, data)

	r.Query(query).Send(data)

	// set user-agent special
	r.Set("User-Agent", site.UserAgent)

	// do request
	_, body, errs := r.EndBytes()

	// with error back
	if len(errs) > 0 {
		return nil, errs[0]
	}

	return path.BodyHook(body)
}

// NewHandler create http handler do request
func NewHandler(site *Site, path *Path) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := doRequest(site, path, r)
		if err != nil {
			site.Flush(w, nil, err)
			return
		}
		site.Flush(w, data, err)
	}
}
