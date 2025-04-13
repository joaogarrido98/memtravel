package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"time"

	"memtravel/configs"
	"memtravel/db"
	"memtravel/handlers"
	"memtravel/middleware"
	"memtravel/ratelimiter"
)

//go:embed templates/*.html static/*.html
var fs embed.FS

var templates = template.Must(template.ParseFS(fs, "templates/*.html"))

func main() {
	// open a connection to the database
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}

	log.Println("database connected")

	defer database.Close()
	defer ratelimiter.ShutdownLimiter()

	// create a new handler which has database and templates available
	handler := handlers.NewHandler(database, templates)

	// create the middlewares we need
	authMiddleware := middleware.CreateStack(middleware.BaseMiddleware, middleware.AuthMiddleware)

	// account deals only with user based interaction
	http.HandleFunc("POST /account/login", middleware.BaseMiddleware(handler.LoginHandler))
	http.HandleFunc("POST /account/register", middleware.BaseMiddleware(handler.RegisterHandler))
	http.HandleFunc("POST /account/password/recover", middleware.BaseMiddleware(handler.PasswordRecoverHandler))
	http.HandleFunc("POST /account/password/change", authMiddleware(handler.PasswordChangeHandler))
	http.HandleFunc("POST /account/close", authMiddleware(handler.CloseAccountHandler))
	http.HandleFunc("POST /account/privacystatus", authMiddleware(handler.PrivacyStatusHandler))
	http.HandleFunc("POST /account/update/country", authMiddleware(handler.UpdateCountryHandler))
	http.HandleFunc("GET /account/activate/{code}", middleware.BaseMiddleware(handler.ActivateAccountHandler))

	// friends deals with anything that is part of the social interaction
	http.HandleFunc("POST /friends/request/{type}", authMiddleware(handler.FriendRequestHandler))
	http.HandleFunc("POST /friends/remove", authMiddleware(handler.RemoveFriendHandler))
	http.HandleFunc("GET /friends/all", authMiddleware(handler.GetFriendsHandler))

	// users deals with any search/user account view
	http.HandleFunc("GET /users/search", authMiddleware(handler.SearchUsersHandler))
	http.HandleFunc("GET /users/account/view", authMiddleware(handler.GetUserHandler))
	http.HandleFunc("POST /users/account/edit", authMiddleware(handler.UserEditHandler))

	// trips deals with anything that is related with the trips
	http.HandleFunc("POST /trips/add", authMiddleware(handler.AddTripHandler))
	http.HandleFunc("GET /trips/upcoming", authMiddleware(handler.GetUpcomingTripsHandler))
	http.HandleFunc("GET /trips/previous", authMiddleware(handler.GetPreviousTripsHandler))
	http.HandleFunc("POST /trips/edit/{id}", authMiddleware(handler.EditTripHandler))
	http.HandleFunc("POST /trips/remove/{id}", authMiddleware(handler.RemoveTripHandler))
	http.HandleFunc("GET /trips/stats", authMiddleware(handler.RemoveTripHandler))

	// pinned
	http.HandleFunc("POST /pinned/add/{tpid}", authMiddleware(handler.AddPinnedHandler))
	http.HandleFunc("POST /pinned/remove/{tpid}", authMiddleware(handler.RemovePinnedHandler))

	// ratings
	http.HandleFunc("POST /ratings/add", authMiddleware(handler.AddRatingHandler))

	// organise
	http.HandleFunc("POST /organise/request", authMiddleware(handler.FriendRequestHandler))

	// country
	http.HandleFunc("GET /country/all", middleware.BaseMiddleware(handler.GetAllCountries))

	// legal
	http.HandleFunc("GET /legal/termsandconditions", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/termsandconditions.html")
	})

	server := &http.Server{
		Addr:         configs.Envs.Port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Println("Server listening on :8080")
	err = server.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}
