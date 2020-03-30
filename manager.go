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
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	// DefaultManager set for default
	DefaultManager = NewManager()
)

// Init init default
func Init(opts ...Option) {
	DefaultManager.Init(opts...)
}

// Register register to default
func Register(sites ...*Site) {
	DefaultManager.Register(sites...)
}

// RegisterWithOptions register site with options
func RegisterWithOptions(opts ...SiteOption) *Site {
	return DefaultManager.RegisterWithOptions(opts...)
}

// Manager ...
type Manager struct {
	inited bool
	r      *mux.Router
	// version -> name -> *Site
	sites map[string]map[string]*Site

	prefix string
}

// RegisterWithOptions register site with options
func (m *Manager) RegisterWithOptions(opts ...SiteOption) *Site {
	var site = NewSite(opts...)
	m.Register(site)
	return site
}

// Register site directlly
func (m *Manager) Register(sites ...*Site) {
	for _, site := range sites {
		var names, ok = m.sites[site.Version]
		if !ok {
			names = make(map[string]*Site)
			m.sites[site.Version] = names
		}

		if _, ok := m.sites[site.Version][site.Name]; ok {
			// TODO: warning
		}
		m.sites[site.Version][site.Name] = site
		log.Printf("registed site name: %s, version: %s\n", site.Name, site.Version)
	}
}

// Init ...
func (m *Manager) Init(opts ...Option) error {
	for _, o := range opts {
		o(m)
	}

	if m.inited {
		return nil
	}

	m.inited = true

	m.r.HandleFunc(m.prefix, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("uniapi server"))
	})

	log.Println("init manager")

	// load all handlers
	for _, names := range m.sites {
		for _, site := range names {
			site.Init()

			for _, p := range site.Paths {
				ver := site.Version
				if p.Version != "" {
					ver = p.Version
				}

				spath := site.Name
				if site.ServePath != "" {
					spath = site.ServePath
				}

				lpath := PathJoin(m.prefix, ver, spath, p.Name)

				log.Printf(
					"server http handler siteName: %s, version: %s, name: %s, path: %s\n",
					site.Name, ver, p.Name, lpath,
				)

				m.r.HandleFunc(lpath, NewHandler(site, p))
			}
		}
	}

	return nil
}

func (m *Manager) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.r.ServeHTTP(w, req)
}

// Do process on special site and path
func (m *Manager) Do(site, path, version string, w http.ResponseWriter, req *http.Request) {

}

// NewManager TODO: add manager option
func NewManager(opts ...Option) *Manager {
	m := &Manager{
		r:     mux.NewRouter(),
		sites: make(map[string]map[string]*Site),
	}

	for _, o := range opts {
		o(m)
	}

	return m
}

// Option ....
type Option func(m *Manager)

// Prefix ...
func Prefix(prefix string) Option {
	return func(m *Manager) {
		m.prefix = prefix
	}
}
