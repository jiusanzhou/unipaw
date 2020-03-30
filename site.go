package unipaw

import (
	"net/http"

	"go.zoe.im/x/httputil"
)

// Site ...
type Site struct {
	Manager *Manager

	Name      string
	Version   string
	ServePath string
	Paths     map[string]*Path  // TODO: replace with sync.Map
	Params    map[string]*Param // TODO: replace with sync.Map
	BodyHook  func([]byte) (interface{}, error)
	Flush     func(w http.ResponseWriter, data interface{}, err error)

	Host      string
	URL       string
	UserAgent string // client
	Network   string // network type
}

// Init site
func (site *Site) Init() {
	for _, path := range site.Paths {
		path.Site = site
		path.Init()
	}
}

// Register paths for site
func (site *Site) Register(paths ...*Path) {
	for _, path := range paths {
		site.Paths[path.Name] = path
	}
}

// RegisterWithOptions register path with opts
func (site *Site) RegisterWithOptions(opts ...PathOption) *Path {
	var path = NewPath(opts...)
	site.Paths[path.Name] = path
	return path
}

// SiteOption is function to set field of Site
type SiteOption func(s *Site)

var (
	defaultVersion   = "v0"
	defaultUserAgent = "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36"
	defaultFlush     = func(w http.ResponseWriter, data interface{}, err error) {
		if err != nil {
			httputil.NewResponse(w).WithError(err).Flush()
			return
		}
		httputil.NewResponse(w).WithData(data).Flush()
	}
)

// NewSite returns a new Site
func NewSite(opts ...SiteOption) *Site {
	s := &Site{
		Version:   defaultVersion,
		UserAgent: defaultUserAgent,
		Flush:     defaultFlush,
		Paths:     make(map[string]*Path),
		Params:    make(map[string]*Param),
	}

	for _, o := range opts {
		o(s)
	}

	return s
}

// SiteName option to set name of Site
func SiteName(name string) SiteOption {
	return func(s *Site) {
		s.Name = name
	}
}

// SiteVersion option to set version of site
func SiteVersion(ver string) SiteOption {
	return func(s *Site) {
		s.Version = ver
	}
}

// SiteServePath option to set serve_path of Site
func SiteServePath(path string) SiteOption {
	return func(s *Site) {
		s.ServePath = path
	}
}

// SitePath option to add path of Site
func SitePath(opts ...PathOption) SiteOption {
	path := NewPath(opts...)
	return func(s *Site) {
		s.Paths[path.Name] = path
	}
}

// SiteParam ...
func SiteParam(opts ...ParamOption) SiteOption {
	var param = NewParam(opts...)
	return func(s *Site) {
		s.Params[param.Name] = param
	}
}

// SiteBodyHook option to set body_book of Site
func SiteBodyHook(f func([]byte) (interface{}, error)) SiteOption {
	return func(s *Site) {
		s.BodyHook = f
	}
}

// SiteFlush option to set flush method of Site
func SiteFlush(f func(w http.ResponseWriter, data interface{}, err error)) SiteOption {
	return func(s *Site) {
		s.Flush = f
	}
}

// SiteHost option to set host of Site
func SiteHost(host string) SiteOption {
	return func(s *Site) {
		s.Host = host
	}
}

// SiteURL option to set URL of Site
func SiteURL(url string) SiteOption {
	return func(s *Site) {
		s.URL = url
	}
}

// SiteUserAgent option to set UserAgent of Site
func SiteUserAgent(userAgent string) SiteOption {
	return func(s *Site) {
		s.UserAgent = userAgent
	}
}

// SiteNetwork option to set network of site
func SiteNetwork(network string) SiteOption {
	return func(s *Site) {

	}
}
