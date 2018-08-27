package servecmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"goproxy/internal/commandline/base"
)

var (
	addr string // address to bind to
	tls  bool   // whether to enable TLS
	cert string // certificate file
	key  string // key file
)

var CmdServe = &base.Command{
	Run:       runRuleSave,
	UsageLine: "serve [-http=:8080] [-dir=./repo] [-tls] [-tls-cert=] [-tls-key=]",
	Short:     "serve packages over HTTP(S)",
	Long: `
Serve starts a Go module proxy.

    -http           address of the proxy
    -tls            whether to use TLS
    -tls-cert       TLS certificate file
    -tls-key        TLS key file
`,
}

func init() {
	CmdServe.Flag.BoolVar(&tls, "tls", false, "enable TLS")
	CmdServe.Flag.StringVar(&cert, "tls-cert", "", "TLS certificate file")
	CmdServe.Flag.StringVar(&key, "tls-key", "", "TLS key file")
	CmdServe.Flag.StringVar(&addr, "http", ":8080", "address of the server")
}

func runRuleSave(cmd *base.Command, args []string) {
	http.HandleFunc("/", handler)
	base.Logf("Serving files from %v at %s", base.Dir, addr)
	if tls {
		if cert == "" || key == "" {
			fmt.Fprintln(os.Stderr, "need both certificate and key files in order to use TLS")
			os.Exit(2)
		}
		panic(http.ListenAndServeTLS(addr, cert, key, nil))
	}
	panic(http.ListenAndServe(addr, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	pkgPath := filepath.Join("", r.RequestURI)

	f, err := os.Open(pkgPath)
	if err != nil {
		http.Error(w, "cannot find package", http.StatusNotFound)
		return
	}
	defer f.Close()

	if _, err := io.Copy(w, f); err != nil {
		http.Error(w, fmt.Sprintf("could not write response: %s", err), http.StatusInternalServerError)
	}
}
