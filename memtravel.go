package main

import (
	"log"
	"net/http"
	"time"

	"memtravel/configs"
	"memtravel/db"
	"memtravel/handlers"
	"memtravel/middleware"
)

func main() {
	// open a connection to the database
	database, err := db.DBConnect()
	if err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}

	log.Println("database connected")

	defer database.Close()

	// create a new handler which has database object available
	handler := handlers.NewHandler(database)

	// create the middlewares we need
	authMiddleware := middleware.CreateStack(middleware.LogMiddleware, middleware.AuthMiddleware)

	// account
	http.HandleFunc("POST /account/login", middleware.LogMiddleware(handler.LoginHandler))
	http.HandleFunc("POST /account/register", middleware.LogMiddleware(handler.RegisterHandler))
	http.HandleFunc("POST /account/password/recover", middleware.LogMiddleware(handler.PasswordRecoverHandler))
	http.HandleFunc("POST /account/password/change", authMiddleware(handler.PasswordChangeHandler))
	http.HandleFunc("POST /account/close", authMiddleware(handler.CloseAccountHandler))
	http.HandleFunc("GET /account/information/view/{id}", authMiddleware(handler.AccountInformationHandler))
	http.HandleFunc("POST /account/information/edit", authMiddleware(handler.AccountInformationEditHandler))
	http.HandleFunc("GET /account/activate/{code}", middleware.LogMiddleware(handler.ActivateAccountHandler))

	// trips
	http.HandleFunc("POST /trips/add", authMiddleware(handler.AddTripHandler))
	http.HandleFunc("GET /trips/upcoming", authMiddleware(handler.GetUpcomingTripsHandler))
	http.HandleFunc("GET /trips/previous", authMiddleware(handler.GetPreviousTripsHandler))
	http.HandleFunc("POST /trips/edit/{id}", authMiddleware(handler.EditTripHandler))
	http.HandleFunc("POST /trips/remove/{id}", authMiddleware(handler.RemoveTripHandler))

	// favourites
	http.HandleFunc("POST /favourites/add", authMiddleware(handler.AddFavouritesHandler))
	http.HandleFunc("POST /favourites/remove/{id}", authMiddleware(handler.RemoveFavouritesHandler))

	// ratings
	http.HandleFunc("POST /ratings/add", authMiddleware(handler.AddRatingHandler))

	// friends
	http.HandleFunc("POST /friends/request", authMiddleware(handler.FriendRequestHandler))
	http.HandleFunc("POST /friends/accept", authMiddleware(handler.AcceptFriendHandler))
	http.HandleFunc("POST /friends/remove/{id}", authMiddleware(handler.RemoveFriendHandler))
	http.HandleFunc("GET /friends/get", authMiddleware(handler.GetFriendsHandler))
	http.HandleFunc("GET /friends/getspecific/{id}", authMiddleware(handler.GetFriendHandler))
	http.HandleFunc("GET /friends/search", authMiddleware(handler.SearchFriendsHandler))

	// organise
	// http.HandleFunc("POST /organise/request", authMiddleware(handler.FriendRequestHandler))

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
