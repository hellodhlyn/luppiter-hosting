package main

import (
	"crypto/tls"

	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var httpClient = &http.Client{}

func proxyRequest(method, path string) (*http.Response, error) {
	url := os.Getenv("LUPPITER_REMOTE_HOST") + strings.ReplaceAll(path, "//", "/")
	req, _ := http.NewRequest(method, url, nil)
	return httpClient.Do(req)
}

func handle(w http.ResponseWriter, r *http.Request) {
	host := strings.Split(r.Host, ":")[0]

	var instance hostingInstance
	db.Where(&hostingInstance{Domain: host}).First(&instance)

	var backend hostingBackend
	db.Where(&hostingBackend{InstanceID: instance.ID}).First(&backend)

	props, err := backend.getProps()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pathPrefix := "/storage/" + props.BucketName + "/" + props.FilePrefix
	res, err := proxyRequest(r.Method, pathPrefix+r.URL.Path)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	if res.StatusCode != http.StatusOK && props.RedirectToIdnex {
		res, err = proxyRequest(r.Method, pathPrefix+"index.html")
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}

	for key, values := range res.Header {
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}
	io.Copy(w, res.Body)
	w.WriteHeader(res.StatusCode)
	res.Body.Close()
}

func main() {
	config := tlsConfig{}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handle)

	srv := &http.Server{
		Addr:    ":8443",
		Handler: mux,
		TLSConfig: &tls.Config{
			GetCertificate: config.GetCertificateFunc(),
		},
	}

	log.Fatal(srv.ListenAndServeTLS("", ""))
}
