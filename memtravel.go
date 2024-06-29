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

	authMiddleware := middleware.CreateStack(middleware.LogMiddleware, middleware.AuthMiddleware)

	// trips
	http.HandleFunc("GET /trips/get", authMiddleware(handlers.GetTripsHandler))
	http.HandleFunc("POST /trips/add", authMiddleware(handlers.AddTripHandler))
	http.HandleFunc("GET /trips/upcoming", authMiddleware(handlers.GetUpcomingTripsHandler))
	http.HandleFunc("GET /trips/previous", authMiddleware(handlers.GetPreviousTripsHandler))
	http.HandleFunc("POST /trips/edit/{id}", authMiddleware(handlers.EditTripHandler))
	http.HandleFunc("POST /trips/remove/{id}", authMiddleware(handlers.RemoveTripHandler))

	// favourites
	http.HandleFunc("POST /favourites/add", authMiddleware(handlers.AddFavouritesHandler))
	http.HandleFunc("POST /favourites/remove/{id}", authMiddleware(handlers.RemoveFavouritesHandler))

	// account
	http.HandleFunc("POST /account/login", middleware.LogMiddleware(handlers.LoginHandler))
	http.HandleFunc("POST /account/register", middleware.LogMiddleware(handlers.RegisterHandler))
	http.HandleFunc("POST /account/password/recover", middleware.LogMiddleware(handlers.PasswordRecoverHandler))
	http.HandleFunc("POST /account/password/change", authMiddleware(handlers.PasswordChangeHandler))
	http.HandleFunc("POST /account/close/{id}", authMiddleware(handlers.CloseAccountHandler))
	http.HandleFunc("GET /account/information/view/{id}", authMiddleware(handlers.AccountInformationHandler))
	http.HandleFunc("POST /account/information/edit/{id}", authMiddleware(handlers.AccountInformationEditHandler))

	// ratings
	http.HandleFunc("POST /ratings/add", authMiddleware(handlers.AddRatingHandler))

	// friends
	http.HandleFunc("POST /friends/request", authMiddleware(handlers.FriendRequestHandler))
	http.HandleFunc("POST /friends/accept", authMiddleware(handlers.AcceptFriendHandler))
	http.HandleFunc("POST /friends/remove/{id}", authMiddleware(handlers.RemoveFriendHandler))
	http.HandleFunc("GET /friends/get", authMiddleware(handlers.GetFriendsHandler))
	http.HandleFunc("GET /friends/getspecific/{id}", authMiddleware(handlers.GetFriendHandler))
	http.HandleFunc("GET /friends/search", authMiddleware(handlers.SearchFriendsHandler))

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
