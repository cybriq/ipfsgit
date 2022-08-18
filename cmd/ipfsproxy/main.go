package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/elazarl/goproxy"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

var port = flag.Int("port", 8008, "Port for proxy to listen on")
var address = flag.String("address", "0.0.0.0",
	"network address to bind pull proxy to")

type Handler struct{}

func initLog() {
	infoHandle, warningHandle, errorHandle := os.Stderr, os.Stderr, os.Stderr
	Info = log.New(infoHandle, "INFO ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(warningHandle, "WARNING ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle, "ERROR ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	initLog()
	flag.Parse()
	h := &Handler{}

	proxy := goproxy.NewProxyHttpServer()
	// proxy.Verbose = true
	proxy.OnRequest().HandleConnectFunc(h.handleConnect)
	proxy.OnRequest().DoFunc(h.handleRequest)

	log.Fatal(http.ListenAndServe(
		fmt.Sprintf("%s:%d", *address, *port), proxy))
}

func (h *Handler) handleRequest(r *http.Request,
	ctx *goproxy.ProxyCtx) (req *http.Request, res *http.Response) {
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
	// Info.Println(spew.Sdump(res))
	// Info.Println(spew.Sdump(r))
	res.Header.Set("Location", Original)
	r.Header.Set("Location", Original)
	// res.Request.URL, _ = url.ParseRequestURI(URI)
	// TODO: search and replace IPFS web interface and web pages in general
	//  to identify IPNS relative link locations and rewrite the content so
	//  the web interface works with this proxy.
	return r, res
}
func (h *Handler) handleConnect(host string,
	ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
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
