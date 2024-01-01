package main

import (
	"database/sql" // New import
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"snippetbox/internal/models"

	_ "github.com/go-sql-driver/mysql" // New import

	"github.com/go-playground/form/v4" // New import
)

// Define an application struct to hold the application-wide dependencies
// for the web application. For now we'll only include fields for the
// two custom loggers, but we'll add more to it as the build progresses.

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	// dsn as Data Source Name, we pass dsn as a parameter to sql.Open()
	dsn := flag.String("dsn", "user:password@/notes?parseTime=true", "MySQL Database Connection")

	// Use log.New() to create a logger for writing information messages.
	// Note that the flags are joined using the bitwise OR operator |.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create a logger for writing error messages in the same way, but use stderr as
	// the destination and use the log.Lshortfile flag to include the relevant file name and line number.
	errorLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a decoder instance...
	formDecoder := form.NewDecoder()

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	// Initialize a new http.Server struct
	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		// Call the new app.routes() method to get the servemux containing our routes.
		Handler: app.routes(),
	}

	// The value returned from the flag.String() function is a pointer to the flag value,
	// not the value itself. So we need to dereference the pointer before using it.
	infoLog.Printf("Starting server on %s", *addr)

	err = server.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	// The sql.Open() function initializes a new sql.DB object,
	// which is essential for pool of database connections.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err

	}
	return db, nil
}
