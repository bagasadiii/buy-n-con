package main

import (
	"log"
	"net/http"

	"github.com/bagasadiii/buy-n-con/app"
	"github.com/bagasadiii/buy-n-con/handler"
	"github.com/bagasadiii/buy-n-con/internal/config"
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
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userServ := service.NewUserService(userRepo)
	userHand := handler.NewUserHandler(userServ)

	itemRepo := repository.NewItemRepository()
	itemServ := service.NewItemService(itemRepo, db)
	itemHand := handler.NewItemHandler(itemServ)

	r := app.SetupRouter(itemHand, userHand)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("failed to run server")
	}
	log.Println("server running")
}