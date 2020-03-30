/*
 * Copyright (c) 2020 wellwell.work, LLC by Zoe
 *
 * Licensed under the Apache License 2.0 (the "License");
 * You may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package unipaw

import (
	"net/http"

	"go.zoe.im/x"
)

// Path serve server and do request
type Path struct {
	Site *Site

	Name        string
	Version     string
	Params      map[string]*Param
	ServeMethod string
	BodyHook    func([]byte) (interface{}, error) // how to flush to response data
	PostData    func() interface{}
	Flush       func(w http.ResponseWriter, data interface{}, err error)

	// TODO: use pure http.Request
	RequestHook func(*Path, *http.Request, map[string]interface{}, interface{})

	URL     string
	Method  string
	Headers map[string]string
}

// Init path
func (path *Path) Init() {
	for _, param := range path.Params {
		param.Path = path
		param.Init()
		path.Flush = path.Site.Flush
		path.Version = path.Site.Version
	}
}

// NewHandler create http handler do request
func (path *Path) NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := doRequest(path.Site, path, r)
		path.Flush(w, data, err)
	}
}

// TODO: register param, register with param

// Do call this path request
// use ...params or PathOption
func (path *Path) Do(opts ...PathOption) (interface{}, error) {
	// copy a new path to handle those options
	var npath = new(Path)
	*npath = *path

	// deep copy map
	npath.Params = make(map[string]*Param)
	for k, v := range path.Params {
		// it's ok don't copy param
		npath.Params[k] = v
	}
	npath.Headers = make(map[string]string)
	for k, v := range path.Headers {
		npath.Headers[k] = v
	}

	for _, o := range opts {
		o(npath)
	}

	// take every tings from path except params

	// r *http.Request
	// TODO:
	var r = &http.Request{}

	return doRequest(npath.Site, npath, r)
}

var (
	defaultMethod   = "GET"
	defaultbodyHook = func(body []byte) (interface{}, error) {
		return x.Bytes2Str(body), nil
	}

	defaultPostDataGetter = func() interface{} {
		return nil
	}

	defaultRequestHook = func(p *Path, r *http.Request, query map[string]interface{}, data interface{}) {
		// only process which can be **incremental update**, like: form, header and cookies

		// copy form
		// for k, v := range r.URL.Query() {
		// 	query[k] = v[0]
		// }

		// copy header ???

		// copy cookies ???

		// copy data ???
	}
)

// PathOption ...
type PathOption func(*Path)

// NewPath new a path
func NewPath(opts ...PathOption) *Path {
	p := &Path{
		Method:   defaultMethod,
		Params:   make(map[string]*Param),
		Headers:  make(map[string]string),
		PostData: defaultPostDataGetter,

		BodyHook:    defaultbodyHook,
		RequestHook: defaultRequestHook,
	}

	for _, o := range opts {
		o(p)
	}

	return p
}

// PathName ...
func PathName(name string) PathOption {
	return func(p *Path) {
		p.Name = name
	}
}

// PathVersion ...
func PathVersion(ver string) PathOption {
	return func(p *Path) {
		p.Version = ver
	}
}

// PathParam ...
func PathParam(opts ...ParamOption) PathOption {
	var param = NewParam(opts...)
	return func(p *Path) {
		p.Params[param.Name] = param
	}
}

// PathServeMethod ...
func PathServeMethod(method string) PathOption {
	return func(p *Path) {
		p.ServeMethod = method
	}
}

// PathBodyHook ...
func PathBodyHook(f func([]byte) (interface{}, error)) PathOption {
	return func(p *Path) {
		p.BodyHook = f
	}
}

// PathURL ...
func PathURL(url string) PathOption {
	return func(p *Path) {
		p.URL = url
	}
}

// PathMethod ...
func PathMethod(method string) PathOption {
	return func(p *Path) {
		p.Method = method
	}
}

// PathHeader ...
func PathHeader(key, value string) PathOption {
	return func(p *Path) {
		p.Headers[key] = value
	}
}

// PathPostData ...
func PathPostData(f func() interface{}) PathOption {
	return func(p *Path) {
		p.PostData = f
	}
}
