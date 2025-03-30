package main

import "net/http"

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/assets", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	server.ListenAndServe()
}
