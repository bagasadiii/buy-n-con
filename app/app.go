package app

import (
	"github.com/bagasadiii/buy-n-con/handler"
	mw "github.com/bagasadiii/buy-n-con/internal/middleware"
	router "github.com/julienschmidt/httprouter"
)
type Routes struct {
	User handler.UserHandlerImpl
	Item handler.ItemHandlerImpl
	Post handler.PostHandlerImpl
}
func SetupRouter(route *Routes)*router.Router{
	r := router.New()

	r.POST("/api/register", route.User.Register)
	r.POST("/api/login", route.User.Login)
	r.GET("/api/u/:username", route.User.GetUserByUsername)

	r.POST("/api/u/:username/items", mw.Auth(route.Item.CreateItem))
	r.GET("/api/u/:username/items/:item_id", route.Item.GetItemByID)
	r.GET("/api/u/:username/items", route.Item.GetAllItems)
	r.PATCH("/api/u/:username/items/:item_id", mw.Auth(route.Item.UpdateItem))
	r.DELETE("/api/u/:username/items/:item_id", mw.Auth(route.Item.DeleteItem))

	r.POST("/api/u/:username/post", mw.Auth(route.Post.CreatePost))
	r.GET("/api/u/:username/post/:post_id", route.Post.GetPostByID)
	r.GET("/api/u/:username/post", route.Post.GetAllPosts)
	r.PATCH("/api/u/:username/post/:post_id", mw.Auth(route.Post.UpdatePost))
	r.DELETE("/api/u/:username/post/:post_id", mw.Auth(route.Post.DeletePost))
	return r
}