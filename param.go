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
)

// Param contains to serve http and send request
type Param struct {
	Path *Path

	Name     string
	Required bool
	Typ      string                                     // string, int, ...
	Value    func(r *http.Request) (interface{}, error) // TODO:
	Query    string                                     // query key
	In       string                                     // no query data header cookie

	Into func(*Param, *http.Request, map[string]interface{}, interface{}) error
}

// Init init param
func (param *Param) Init() {
	// do nothing at all
}

var (
	defaultRequired = true
	defaultValue    = func(r *http.Request) (interface{}, error) {
		return "", nil
	}
	defaultInto = func(p *Param, r *http.Request, query map[string]interface{}, data interface{}) error {
		// TODO:
		// else query
		// name must be request param name
		// set to where, today only query
		// TODO: set to data
		v, ok := r.URL.Query()[p.Name]
		if ok {
			query[p.Query] = v[0]
			return nil
		}
		var err error

		query[p.Query], err = p.Value(r)

		// add option to ingore error from value call
		return err
	}
)

// ParamOption is a option to set field of param
type ParamOption func(p *Param)

// NewParam ...
func NewParam(opts ...ParamOption) *Param {
	param := &Param{
		Required: defaultRequired,
		Into:     defaultInto,
		Value:    defaultValue,
	}

	for _, o := range opts {
		o(param)
	}

	if param.Query == "" {
		param.Query = param.Name
	}

	return param
}

// ParamName is a option to set name of Param
func ParamName(name string) ParamOption {
	return func(p *Param) {
		p.Name = name
	}
}

// ParamRequired is a option set required of Param
func ParamRequired(required bool) ParamOption {
	return func(p *Param) {
		p.Required = required
	}
}

// ParamType ...
func ParamType(typ string) ParamOption {
	return func(p *Param) {
		p.Typ = typ
	}
}

// ParamValue ...
func ParamValue(f interface{}) ParamOption {
	return func(p *Param) {
		// TODO: checkout the type of f
		fn, ok := f.(func(r *http.Request) (interface{}, error))
		if ok {
			p.Value = fn
			return
		}
		p.Value = func(r *http.Request) (interface{}, error) { return f, nil }
	}
}

// ParamQuery ...
func ParamQuery(name string) ParamOption {
	return func(p *Param) {
		p.Query = name
	}
}

// ParamIn ...
func ParamIn(in string) ParamOption {
	return func(p *Param) {
		p.In = in
	}
}

// ParamInto ...
func ParamInto(f func(*Param, *http.Request, map[string]interface{}, interface{}) error) ParamOption {
	return func(p *Param) {
		p.Into = f
	}
}
