package rentals

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"tournaments/roles"
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

	srv := &http.Server{
		Handler:      app.Server.Router,
		Addr:         app.Addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()
}

func (app *App) Setup() error {
	return app.Server.Setup()
}

func NewApp(addr string) (*App, error) {
	db, err := ConnectToDB()
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
