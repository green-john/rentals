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
	server *Server
	addr   string
}

// Helper method to enable easier testing
func (app *App) dropDB() {
	app.server.db.DropTableIfExists(DbModels...)
}

// Serves http requests. app.server must be initialized,
// otherwise an error is thrown
func (app *App) ServeHTTP() error {
	if !app.server.initialized {
		return errors.New("app must be initialized first")
	}

	srv := &http.Server{
		Handler:      app.server.router,
		Addr:         app.addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()
}

func (app *App) Setup() error {
	return app.server.Setup()
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
		db:          db,
		router:      router,
		authN:       authN,
		authZ:       authZ,
		initialized: false,
	}

	app := &App{
		server: server,
		addr:   addr,
	}

	return app, nil
}
