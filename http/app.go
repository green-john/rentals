package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"rentals"
	"rentals/roles"
)

type App struct {
	Server *Server
	Addr   string
}

// Helper method to enable easier testing
func (app *App) DropDB() {
	app.Server.Db.DropTableIfExists(rentals.DbModels...)
}

func (app *App) Setup() error {
	return app.Server.Setup()
}

func (app *App) ServeHTTP() error {
	return app.Server.ServeHTTP(app.Addr)
}

func NewApp(addr string, testing bool) (*App, error) {
	db, err := rentals.ConnectToDB(testing)
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
