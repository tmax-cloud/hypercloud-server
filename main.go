package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	claim "github.com/tmax-cloud/hypercloud-server/claim"
	cluster "github.com/tmax-cloud/hypercloud-server/cluster"
	"k8s.io/klog"
)

var (
	port     int
	certFile string
	keyFile  string
)

func main() {
	flag.IntVar(&port, "port", 80, "hypercloud server port")
	// flag.StringVar(&certFile, "certFile", "/run/secrets/tls/server.crt", "hypercloud server cert")
	// flag.StringVar(&keyFile, "keyFile", "/run/secrets/tls/server.key", "x509 Private key file for TLS connection")
	flag.Parse()

	// keyPair, err := tls.LoadX509KeyPair(certFile, keyFile)
	// if err != nil {
	// 	klog.Errorf("Failed to load key pair: %s", err)
	// }

	mux := http.NewServeMux()

	mux.HandleFunc("/api/master/clusterclaim", serveClusterClaim)
	mux.HandleFunc("/api/master/cluster", serveCluster)
	mux.HandleFunc("/api/master/cluster/owner", serveClusterOwner)
	mux.HandleFunc("/api/master/cluster/member", serveClusterMember)
	mux.HandleFunc("/api/master/test/", serveTest)

	svr := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
		// TLSConfig: &tls.Config{Certificates: []tls.Certificate{keyPair}},
	}

	klog.Info("Starting hypercloud5 server...")

	go func() {
		if err := svr.ListenAndServe(); err != nil { //HTTPS로 서버 시작
			// if err := svr.ListenAndServeTLS("", ""); err != nil { //HTTPS로 서버 시작
			klog.Errorf("Failed to listen and serve hypercloud5 server: %s", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	klog.Info("OS shutdown signal received...")
	svr.Shutdown(context.Background())
}

func serveTest(w http.ResponseWriter, r *http.Request) {
	klog.Infof("Http request: method=%s, uri=%s", r.Method, r.URL.Path)
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	klog.Info("Request body: \n", string(body))
}

func serveClusterClaim(res http.ResponseWriter, req *http.Request) {
	klog.Infof("Http request: method=%s, uri=%s", req.Method, req.URL.Path)
	switch req.Method {
	case http.MethodGet:
		//curl -XPUT 172.22.6.2:32319/api/master/clusterclaim?userId=sangwon_cho@tmax.co.kr
		claim.List(res, req)
	case http.MethodPost:
	case http.MethodPut:
		//curl -XPUT 172.22.6.2:32319/api/master/clusterclaim?userId=sangwon_cho@tmax.co.kr\&clusterClaim=test-d5n92\&admit=true
		claim.Put(res, req)
	case http.MethodDelete:
	default:
	}
}

func serveCluster(res http.ResponseWriter, req *http.Request) {
	klog.Infof("Http request: method=%s, uri=%s", req.Method, req.URL.Path)
	switch req.Method {
	case http.MethodGet:
		cluster.List(res, req)
	case http.MethodPost:
	case http.MethodPut:
		// invite multiple users
		cluster.Put(res, req)
	case http.MethodDelete:
	default:
	}
}

func serveClusterOwner(res http.ResponseWriter, req *http.Request) {
	klog.Infof("Http request: method=%s, uri=%s", req.Method, req.URL.Path)
	switch req.Method {
	case http.MethodGet:
		cluster.ListOwner(res, req)
	case http.MethodPost:
	case http.MethodPut:
	case http.MethodDelete:
	default:
	}
}

func serveClusterMember(res http.ResponseWriter, req *http.Request) {
	klog.Infof("Http request: method=%s, uri=%s", req.Method, req.URL.Path)
	switch req.Method {
	case http.MethodGet:
		cluster.ListMember(res, req)
	case http.MethodPost:
	case http.MethodPut:
	case http.MethodDelete:
	default:
	}
}
