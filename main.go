package main

import (
	"log"
	"net/http"

	"github.com/bagasadiii/buy-n-con/handler"
	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/bagasadiii/buy-n-con/internal/config"
	"github.com/bagasadiii/buy-n-con/internal/middleware"
	"github.com/bagasadiii/buy-n-con/internal/repository"
	"github.com/bagasadiii/buy-n-con/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("failed to get .env file: ", err)
	}
	db := config.DBConnection()
	userRepo := repository.NewUserRepository(db)
	userServ := service.NewUserService(userRepo)
	userHand := handler.NewUserHandler(userServ)

	http.HandleFunc("/register", userHand.Register)
	http.HandleFunc("/login", userHand.Login)
	http.HandleFunc("/hello", middleware.Auth(hello))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("failed to run server")
	}
	log.Println("server running")
}
func hello(w http.ResponseWriter, r *http.Request){
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
    if !ok {
        http.Error(w, "userID not found in context", http.StatusInternalServerError)
        return
    }

    username, ok := r.Context().Value(middleware.UsernameKey).(string)
    if !ok {
        http.Error(w, "username not found in context", http.StatusInternalServerError)
        return
    }

	res := helper.Response{
		Status: http.StatusOK,
		Message: "hello kontol",
		Data: map[string]string{
			"userid": userID,
			"username": username,
		},
		Err: nil,
	}
	helper.JSONResponse(w, res.Status, res)
}