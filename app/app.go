package app

import (
	"net/http"

	"github.com/bagasadiii/buy-n-con/handler"
	"github.com/bagasadiii/buy-n-con/internal/middleware"
	router "github.com/julienschmidt/httprouter"
)

func SetupRouter(item handler.ItemHandlerImpl, user handler.UserHandlerImpl)*router.Router{
	r := router.New()

	r.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }
    })

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