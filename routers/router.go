package routers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lscantillo/twitter-clone-api/handlers"
	"github.com/lscantillo/twitter-clone-api/middlewares"
	"github.com/rs/cors"
)

// App struct to hold the router
type App struct {
	Router *mux.Router
}

// Initialize the app function to initialize the router
func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run the app function to start the server
func (a *App) Run(addr string) {
	log.Println("Server running...")
	handler := cors.AllowAll().Handler(a.Router)
	err := http.ListenAndServe(addr, handler)
	if err != nil {
		return
	}
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/", middlewares.CheckDB(handlers.HomeHandler)).Methods("GET")
	a.Router.HandleFunc("/register", middlewares.CheckDB(handlers.RegisterHandler)).Methods("POST")
	a.Router.HandleFunc("/login", middlewares.CheckDB(handlers.LoginHandler)).Methods("POST")
	a.Router.HandleFunc("/profile", middlewares.CheckDB(middlewares.ValidateJWT(handlers.ProfileHandler))).Methods("GET")
	// a.Router.HandleFunc("/read_characters/{type}/{items}/{itemsPerWorker}", controller.ReadCharacters).Methods("GET")
}
