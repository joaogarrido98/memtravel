package main

import (
	"log"
	"net/http"
	"time"

	"memtravel/configs"
	"memtravel/middleware"
	"memtravel/routes"
)

func main() {
	// trips
	http.HandleFunc("GET /trips/get", middleware.AuthMiddleware(routes.GetTripsHandler))
	http.HandleFunc("POST /trips/add", middleware.AuthMiddleware(routes.AddTripHandler))
	http.HandleFunc("GET /trips/upcoming", middleware.AuthMiddleware(routes.GetUpcomingTripsHandler))
	http.HandleFunc("GET /trips/previous", middleware.AuthMiddleware(routes.GetPreviousTripsHandler))
	http.HandleFunc("POST /trips/edit/{id}", middleware.AuthMiddleware(routes.EditTripHandler))
	http.HandleFunc("POST /trips/remove/{id}", middleware.AuthMiddleware(routes.RemoveTripHandler))

	// favourites
	http.HandleFunc("POST /favourites/add", middleware.AuthMiddleware(routes.AddFavouritesHandler))
	http.HandleFunc("POST /favourites/remove/{id}", middleware.AuthMiddleware(routes.RemoveFavouritesHandler))

	// account
	http.HandleFunc("POST /account/login", routes.LoginHandler)
	http.HandleFunc("POST /account/register", routes.RegisterHandler)
	http.HandleFunc("POST /account/close/{id}", middleware.AuthMiddleware(routes.CloseAccountHandler))
	http.HandleFunc("GET /account/information/view/{id}", middleware.AuthMiddleware(routes.AccountInformationHandler))
	http.HandleFunc("POST /account/information/edit/{id}", middleware.AuthMiddleware(routes.AccountInformationEditHandler))

	// ratings
	http.HandleFunc("POST /ratings/add", middleware.AuthMiddleware(routes.AddRatingHandler))

	// friends
	http.HandleFunc("POST /friends/request", middleware.AuthMiddleware(routes.FriendRequestHandler))
	http.HandleFunc("POST /friends/accept", middleware.AuthMiddleware(routes.AcceptFriendHandler))
	http.HandleFunc("POST /friends/remove/{id}", middleware.AuthMiddleware(routes.RemoveFriendHandler))
	http.HandleFunc("GET /friends/get", middleware.AuthMiddleware(routes.GetFriendsHandler))
	http.HandleFunc("GET /friends/getspecific/{id}", middleware.AuthMiddleware(routes.GetFriendHandler))
	http.HandleFunc("GET /friends/search", middleware.AuthMiddleware(routes.SearchFriendsHandler))

	server := &http.Server{
		Addr:         configs.Envs.Port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Println("Server listening on :8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}
