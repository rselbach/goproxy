package proxy

// proxy is a Go module proxy
type Proxy struct {
	dir string
}

// New initializes a new proxy
func New(dir string) *Proxy {
	return &Proxy{dir: dir}
}
