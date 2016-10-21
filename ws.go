package main

import (
	"fmt"
	"net/http"
	"os"
)

// Log is a very simple access log handler
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %d %s %s", r.RemoteAddr, r.Response.StatusCode, r.Method, r.URL)
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

	fmt.Printf("Starting webserver with file root at '%s', with ListenAndServe at '%s'\n", path, port)
	http.Handle("/", Log(http.StripPrefix("/", http.FileServer(http.Dir(path)))))
	if err = http.ListenAndServe(port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}
