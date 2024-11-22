package main

import (
	"log"
	"net/http"

	"github.com/bagasadiii/buy-n-con/handler"
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

	itemRepo := repository.NewItemRepository()
	itemServ := service.NewItemService(itemRepo, db)
	itemHand := handler.NewItemHandler(itemServ)

	http.HandleFunc("/register", userHand.Register)
	http.HandleFunc("/login", userHand.Login)
	http.HandleFunc("/item", middleware.Auth(itemHand.CreateItem))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("failed to run server")
	}
	log.Println("server running")
}