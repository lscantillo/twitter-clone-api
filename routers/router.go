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
	a.Router.HandleFunc("/updateprofile", middlewares.CheckDB(middlewares.ValidateJWT(handlers.UpdateProfileHandler))).Methods("PUT")

	a.Router.HandleFunc("/tweet", middlewares.CheckDB(middlewares.ValidateJWT(handlers.CreateTweet))).Methods("POST")
	a.Router.HandleFunc("/get_tweets", middlewares.CheckDB(middlewares.ValidateJWT(handlers.GetTweets))).Methods("GET")
	a.Router.HandleFunc("/delete_tweet", middlewares.CheckDB(middlewares.ValidateJWT(handlers.DeleteTweet))).Methods("DELETE")

	a.Router.HandleFunc("/upload_avatar", middlewares.CheckDB(middlewares.ValidateJWT(handlers.UploadAvatar))).Methods("POST")
	a.Router.HandleFunc("/get_avatar", middlewares.CheckDB(handlers.GetAvatar)).Methods("GET")
	a.Router.HandleFunc("/upload_banner", middlewares.CheckDB(middlewares.ValidateJWT(handlers.UploadBanner))).Methods("POST")
	a.Router.HandleFunc("/get_banner", middlewares.CheckDB(handlers.GetBanner)).Methods("GET")

	a.Router.HandleFunc("/create_relation", middlewares.CheckDB(middlewares.ValidateJWT(handlers.CreateRelation))).Methods("POST")
	a.Router.HandleFunc("/delete_relation", middlewares.CheckDB(middlewares.ValidateJWT(handlers.DeleteRelation))).Methods("DELETE")
	a.Router.HandleFunc("/get_relation", middlewares.CheckDB(middlewares.ValidateJWT(handlers.GetRelation))).Methods("GET")

	a.Router.HandleFunc("/get_users", middlewares.CheckDB(middlewares.ValidateJWT(handlers.GetUsers))).Methods("GET")
	a.Router.HandleFunc("/read_followers_tweets", middlewares.CheckDB(middlewares.ValidateJWT(handlers.ReadTweets))).Methods("GET")

}
