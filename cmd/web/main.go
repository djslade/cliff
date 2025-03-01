package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	app := &application{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	app.logger.Info("starting server", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, app.routes())

	app.logger.Error(err.Error())
	os.Exit(1)

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
