package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/djslade/cliff/internal/models"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/subosito/gotenv"
)

type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	gotenv.Load(".env")

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsm := flag.String("dsm", os.Getenv("MYSQL_DSM"), "MySQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsm)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	app.logger.Info("starting server", slog.String("addr", *addr))

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
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

func openDB(dsm string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsm)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
