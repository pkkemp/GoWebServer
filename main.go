package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	time2 "time"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Write "Hello, world!" to the response body
	time := time2.Now()
	text := "Hello world!\n" + time.String() + "\n"
	text += "Remote IP:" + 	r.RemoteAddr + "\n"
	text+= "Forwarded For:" + r.Header.Get("X-Forwarded-For") + "\n"
	text+= "Forwarding Protocol: " + r.Header.Get("X-Forwarded-Proto") + "\n"
	text+= "CDN Loop: " + r.Header.Get("CDN-Loop")
	io.WriteString(w, text)
}

func main() {
	// Set up a /hello resource handler

	http.HandleFunc("/hello", helloHandler)

	// Create a CA certificate pool and add cert.pem to it
	caCert, err := ioutil.ReadFile("cert.pem")
	cloudflare, err := ioutil.ReadFile("cloudflare.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	caCertPool.AppendCertsFromPEM(cloudflare)

	// Create the TLS Config with the CA pool and enable Client certificate validation
	tlsConfig := &tls.Config{
		ClientCAs: caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	// Create a Server instance to listen on port 8443 with the TLS config
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}

	// Listen to HTTPS connections with the server certificate and wait
	log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
}