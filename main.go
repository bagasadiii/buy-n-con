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
	"github.com/rs/cors"
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

	postRepo := repository.NewPostRepository()
	postServ := service.NewServiceImpl(postRepo, db)
	postHand := handler.NewPostHandler(postServ)

	route := app.Routes{
		User: userHand,
		Item: itemHand,
		Post: postHand,
	}

	r := app.SetupRouter(&route)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowCredentials: true,
		Debug: true,
	})

	if err := http.ListenAndServe(":8080", c.Handler(r)); err != nil {
		log.Fatalf("failed to run server: %v\n", err)
	}
	log.Println("server running")
}