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
	err := db.DBConnect()
	if err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}

	log.Println("database connected")

	// trips
	http.HandleFunc("GET /trips/get", middleware.AuthMiddleware(handlers.GetTripsHandler))
	http.HandleFunc("POST /trips/add", middleware.AuthMiddleware(handlers.AddTripHandler))
	http.HandleFunc("GET /trips/upcoming", middleware.AuthMiddleware(handlers.GetUpcomingTripsHandler))
	http.HandleFunc("GET /trips/previous", middleware.AuthMiddleware(handlers.GetPreviousTripsHandler))
	http.HandleFunc("POST /trips/edit/{id}", middleware.AuthMiddleware(handlers.EditTripHandler))
	http.HandleFunc("POST /trips/remove/{id}", middleware.AuthMiddleware(handlers.RemoveTripHandler))

	// favourites
	http.HandleFunc("POST /favourites/add", middleware.AuthMiddleware(handlers.AddFavouritesHandler))
	http.HandleFunc("POST /favourites/remove/{id}", middleware.AuthMiddleware(handlers.RemoveFavouritesHandler))

	// account
	http.HandleFunc("POST /account/login", handlers.LoginHandler)
	http.HandleFunc("POST /account/register", handlers.RegisterHandler)
	http.HandleFunc("POST /account/password/recover", handlers.PasswordRecoverHandler)
	http.HandleFunc("POST /account/password/change", middleware.AuthMiddleware(handlers.PasswordChangeHandler))
	http.HandleFunc("POST /account/close/{id}", middleware.AuthMiddleware(handlers.CloseAccountHandler))
	http.HandleFunc("GET /account/information/view/{id}", middleware.AuthMiddleware(handlers.AccountInformationHandler))
	http.HandleFunc("POST /account/information/edit/{id}", middleware.AuthMiddleware(handlers.AccountInformationEditHandler))

	// ratings
	http.HandleFunc("POST /ratings/add", middleware.AuthMiddleware(handlers.AddRatingHandler))

	// friends
	http.HandleFunc("POST /friends/request", middleware.AuthMiddleware(handlers.FriendRequestHandler))
	http.HandleFunc("POST /friends/accept", middleware.AuthMiddleware(handlers.AcceptFriendHandler))
	http.HandleFunc("POST /friends/remove/{id}", middleware.AuthMiddleware(handlers.RemoveFriendHandler))
	http.HandleFunc("GET /friends/get", middleware.AuthMiddleware(handlers.GetFriendsHandler))
	http.HandleFunc("GET /friends/getspecific/{id}", middleware.AuthMiddleware(handlers.GetFriendHandler))
	http.HandleFunc("GET /friends/search", middleware.AuthMiddleware(handlers.SearchFriendsHandler))

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
