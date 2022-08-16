package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/elazarl/goproxy"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

type Handler struct{}

func initLog(
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Info = log.New(infoHandle,
		"INFO ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	initLog(os.Stdout, os.Stdout, os.Stderr)

	h := &Handler{}

	proxy := goproxy.NewProxyHttpServer()
	// proxy.Verbose = true
	proxy.OnRequest().HandleConnectFunc(h.handleConnect)
	proxy.OnRequest().DoFunc(h.handleRequest)

	log.Fatal(http.ListenAndServe(":8008", proxy))
}

func (h *Handler) handleRequest(r *http.Request, ctx *goproxy.ProxyCtx) (req *http.Request, res *http.Response) {
	var host string
	if !strings.Contains(r.URL.Host, ":") {
		host = r.URL.Host + ":80"
	} else {
		host = r.URL.Host
	}
	Original := r.URL.String()
	Info.Printf("handling request for host '%s' URL '%s'", host, Original)
	split := strings.Split(r.URL.RequestURI(), "/")
	var URI string
	if len(split) > 2 {
		URI = "/" + strings.Join(split[2:], "/")
	}
	Info.Println(split)
	URL := "http://" + split[1] + "." + r.URL.Host + ".localhost:8080" + URI
	Info.Printf("rewriting to '%s'", URL)
	var err error
	res, err = http.Get(URL)
	if err != nil {
		Error.Printf("error: '%s'", err.Error())
		return r, nil
	}
	// res.Header.Set("", URI)
	Info.Println(spew.Sdump(res))
	Info.Println(spew.Sdump(r))
	res.Header.Set("Location", Original)
	r.Header.Set("Location", Original)
	// res.Request.URL, _ = url.ParseRequestURI(URI)
	return r, res
}
func (h *Handler) handleConnect(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	Info.Printf("connecting to '%s'", host)
	switch host {
	case "ipns", "ipfs", "ipld":
		Info.Printf("redirecting to local IPFS service '%s'", host)
	default:
		Info.Printf("connecting to '%s'", host)
	}
	// Info.Println(spew.Sdump(ctx))
	return goproxy.OkConnect, host
}
