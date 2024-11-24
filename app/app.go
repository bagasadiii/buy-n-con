package app

import (
	"github.com/bagasadiii/buy-n-con/handler"
	"github.com/bagasadiii/buy-n-con/internal/middleware"
	router "github.com/julienschmidt/httprouter"
)

func SetupRouter(item handler.ItemHandlerImpl, user handler.UserHandlerImpl)*router.Router{
	r := router.New()

	r.POST("/api/register", user.Register)
	r.POST("/api/login", user.Login)
	r.GET("/api/u/:username", user.GetUserByUsername)

	r.POST("/api/u/:username/items", middleware.Auth(item.CreateItem))
	r.GET("/api/u/:username/items/:item_id", item.GetItemByID)
	r.GET("/api/u/:username/items", item.GetAllItems)
	r.PATCH("/api/u/:username/items/:item_id", middleware.Auth(item.UpdateItem))
	r.DELETE("/api/u/:username/items/:item_id", middleware.Auth(item.DeleteItem))

	return r
}