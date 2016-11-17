package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
)

// Log is a very simple access log handler
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s\t%s\t%s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	var err error
	var path string
	var port string

	if len(os.Args) >= 2 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			fmt.Fprintf(os.Stderr, "\nSYNTAX:  %s [ PATH ] [ ADDR ]\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "\nEXAMPLES:\n\n")
			fmt.Fprintf(os.Stderr, "  %s\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "  %s /var/www\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "  %s /var/www :8080\n\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "Environment Variables:\n\n")
			fmt.Fprintf(os.Stderr, "  TLS_CERT_PATH:  if set, path to PEM encoded certificate chain\n")
			fmt.Fprintf(os.Stderr, "  TLS_KEY_PATH:   if set, path to PEM encoded key\n\n")
			fmt.Fprintf(os.Stderr, "Note: if either of the variables are set, this becomes an HTTPS server.\n\n")
			os.Exit(2)
		}
	}

	// get port or default to localhost:8080
	if len(os.Args) >= 3 {
		port = os.Args[2]
	} else {
		port = "localhost:8080"
	}

	// get path or use current working directory
	if len(os.Args) >= 2 {
		path = os.Args[1]
	} else {
		path, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
			os.Exit(2)
		}

	}

	// get environment variables and see if we want to use TLS
	useTLS := false
	tlsCertFile := os.Getenv("TLS_CERT_PATH")
	tlsKeyFile := os.Getenv("TLS_KEY_PATH")
	if len(tlsCertFile) > 0 && len(tlsKeyFile) > 0 {
		useTLS = true
	}

	// https
	if useTLS {
		tlsCfg := tls.Config{
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS12,
		}

		mux := http.NewServeMux()
		mux.Handle("/", Log(http.StripPrefix("/", http.FileServer(http.Dir(path)))))

		s := &http.Server{
			Addr:      port,
			TLSConfig: &tlsCfg,
			Handler:   mux,
		}

		fmt.Printf("Starting HTTPS webserver with file root at '%s', with ListenAndServe at '%s'\n", path, port)
		if err = s.ListenAndServeTLS(tlsCertFile, tlsKeyFile); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
			os.Exit(1)
		}
		// http
	} else {
		http.Handle("/", Log(http.StripPrefix("/", http.FileServer(http.Dir(path)))))
		fmt.Printf("Starting HTTP webserver with file root at '%s', with ListenAndServe at '%s'\n", path, port)
		if err = http.ListenAndServe(port, nil); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
			os.Exit(1)
		}
	}
}
