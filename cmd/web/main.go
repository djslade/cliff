package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Print("starting server on :4000")

	log.Fatal(http.ListenAndServe(":4000", mux))

	/*
		http.ListenAndServe()
		addr parameter expects format 'host:port'
		When the host portion is omitted the function will listen on the machine's
		available network interfaces. That's mostly fine unless the device has more
		than one interface and you only want to listen to one of them.

		The port doesn't have to be a number. In some cases you may see a project or
		piece of documentation use something like ':http' instead. If that's the case,
		this function will look up the relevant port number from your /etc/services file,
		and will throw an error if unable to find a match.
	*/
}
