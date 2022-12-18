package main

// Packages to import
import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/hbourgeot/portfolio/internal/forms"

	_ "github.com/go-sql-driver/mysql"
)

// Struct for dependency injection
type portfolio struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	messages       *forms.MessageModel
	users          *forms.UserModel
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// Portfolio flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "root:Wini.h16b.@/portfolio?parseTime=true", "MySQL data source name")

	flag.Parse()

	// Loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Connection to database
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close() // finally, close the db connection

	// Form decoder
	formDecoder := form.NewDecoder()

	// Session manager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// Assign values to portfolio struct
	app := &portfolio{
		errorLog:       errorLog,
		infoLog:        infoLog,
		formDecoder:    formDecoder,
		messages:       &forms.MessageModel{DB: db},
		users:          &forms.UserModel{DB: db},
		sessionManager: sessionManager,
	}

	// Server config
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Println("Starting server on port 4000")
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
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
