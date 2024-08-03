package main

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type App struct {
	Router *mux.Router
}

func (a *App) Initialize() {

	log.Info("In Initialize")
	a.Router = mux.NewRouter()
	a.initializeRoutes()

}

func (a *App) Run(addr string) {
	log.Print(fmt.Sprintf("Server running on port [%s]", addr))
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
