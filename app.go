package rentals

import (
	"errors"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"rentals/roles"
	"time"
)

type App struct {
	Server *Server
	Addr   string
}

// Helper method to enable easier testing
func (app *App) DropDB() {
	app.Server.Db.DropTableIfExists(DbModels...)
}

// Serves http requests. app.Server must be Initialized,
// otherwise an error is thrown
func (app *App) ServeHTTP() error {
	if !app.Server.Initialized {
		return errors.New("app must be Initialized first")
	}

	// Enable CORS for testing purposes. This should be
	// configured properly for production
	allOrigins := handlers.AllowedOrigins([]string{"*"})
	allMethods := handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"})
	allHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})

	// Log all requests
	r := handlers.LoggingHandler(os.Stderr, app.Server.Router)

	srv := &http.Server{
		Handler:      handlers.CORS(allOrigins, allMethods, allHeaders)(r),
		Addr:         app.Addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()
}

func (app *App) Setup() error {
	return app.Server.Setup()
}

func NewApp(addr string, testing bool) (*App, error) {
	db, err := ConnectToDB(testing)
	if err != nil {
		return nil, fmt.Errorf("[NewApp] error in ConnectToDB(): %v", err)
	}

	router := mux.NewRouter()
	authN := NewDbAuthenticator(db)
	authZ := roles.NewAuthorizer()

	server := &Server{
		Db:          db,
		Router:      router,
		AuthN:       authN,
		AuthZ:       authZ,
		Initialized: false,
	}

	app := &App{
		Server: server,
		Addr:   addr,
	}

	return app, nil
}
