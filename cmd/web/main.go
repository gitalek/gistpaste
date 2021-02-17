package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/gitalek/gistpaste/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

// application holds the application-wide deps
type application struct {
	errorLog      *log.Logger
	gists         *mysql.GistModel
	infoLog       *log.Logger
	session       *sessions.Session
	templateCache map[string]*template.Template
	users         *mysql.UserModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	// dsn is a data source name (or connection string)
	dsn := flag.String(
		"dsn",
		"web:password@/gistpaste?parseTime=true",
		"MySQL data source name",
	)
	// Define a new command-line flag for the session secret (a random key which
	// will be used to encrypt and authenticate session cookies).
	// It should be 32 bytes long.
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a new session manager.
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	// set true to serve all request only over HTTPS
	session.Secure = true

	app := &application{
		errorLog:      errorLog,
		gists:         &mysql.GistModel{DB: db},
		infoLog:       infolog,
		session:       session,
		templateCache: templateCache,
		users:         &mysql.UserModel{DB: db},
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	infolog.Printf("Starting server on %s\n", *addr)
	errorLog.Fatal(srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"))
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
