package rentals

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type App struct {
	Srv *http.Server
	DB  *gorm.DB
}

func (app *App) dropDB() {
	app.DB.DropTableIfExists(DbModels...)
}

func (app *App) Serve() error {
	return app.Srv.ListenAndServe()
}

func (app *App) getServerURL() string {
	return app.Srv.Addr
}

func NewApp(port int) (*App, error) {
	db, err := ConnectToDB()
	if err != nil {
		return nil, fmt.Errorf("[NewApp] error in ConnectToDB(): %v", err)
	}

	router, err := initRouter(db)
	if err != nil {
		return nil, fmt.Errorf("[NewApp] error in initRouter(): %v", err)
	}

	srv := &http.Server{
		Handler:      router.Router,
		Addr:         fmt.Sprintf("localhost:%d", port), // TODO dont hardcode this value
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	app := &App{
		Srv: srv,
		DB:  db,
	}

	return app, nil
}

func initRouter(db *gorm.DB) (*GorillaRouter, error) {
	mux := NewGorillaRouter()

	resources := []Resource{&UserResource{db}}

	for _, resource := range resources {
		CreateRoutes(resource, mux)
	}

	return mux, nil
}
