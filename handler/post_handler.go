package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/bagasadiii/buy-n-con/internal/middleware"
	"github.com/bagasadiii/buy-n-con/internal/model"
	"github.com/bagasadiii/buy-n-con/internal/service"
	"github.com/google/uuid"
	router "github.com/julienschmidt/httprouter"
)

type PostHandlerImpl interface{
	CreatePost(w http.ResponseWriter, r *http.Request, p router.Params)
	GetPostByID(w http.ResponseWriter, r *http.Request, p router.Params)
	GetAllPosts(w http.ResponseWriter, r *http.Request, p router.Params)
	UpdatePost(w http.ResponseWriter, r *http.Request, p router.Params)
	DeletePost(w http.ResponseWriter, r *http.Request, p router.Params)
}

type PostHandler struct {
	serv service.PostServiceImpl
}

func NewPostHandler(serv service.PostServiceImpl)PostHandlerImpl{
	return &PostHandler{
		serv: serv,
	}
}

func(h *PostHandler)CreatePost(w http.ResponseWriter, r *http.Request, p router.Params){
	ctx := r.Context()
	userCtx, ok := ctx.Value(middleware.UserContextKey).(*middleware.ContextKey)
	if !ok {
		res := helper.UnauthorizedErr("Unauthorized: ", nil)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	username := p.ByName("username")
	if username != userCtx.UsernameKey {
		res := helper.ForbiddenErr("Forbidden access: ", nil)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	var input model.PostInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		res := helper.InternalErr("Internal error: ", err)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	post, err := h.serv.CreatePostService(ctx, &input)
	if err != nil {
		res := helper.InternalErr("Failed to create post: ", err)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	res := helper.Response{
		Status: http.StatusCreated,
		Message: "post created",
		Data: &post,
		Err: nil,
	}
	helper.JSONResponse(w, res.Status, res)
}
func(h *PostHandler)GetPostByID(w http.ResponseWriter, r *http.Request, p router.Params){
	username := p.ByName("username")
	postID, err := uuid.Parse(p.ByName("post_id"))
	if err != nil {
		res := helper.BadRequestErr("Bad request: Invalid post ID", nil)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	input := &model.GetPostInput{
		PostID: postID,
		Owner: username,
	}
	post, err := h.serv.GetPostByIDService(r.Context(), input)
	if err != nil {
		res := helper.InternalErr("Failed to fetch post: ", err)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	res := helper.Response{
		Status: http.StatusOK,
		Message: "OK",
		Data: &post,
		Err: nil,
	}
	helper.JSONResponse(w, res.Status, res)
}
func(h *PostHandler)GetAllPosts(w http.ResponseWriter, r *http.Request, p router.Params){
	username := p.ByName("username")
	if username == "" {
		res := helper.BadRequestErr("Username is required ", nil)
		helper.JSONResponse(w, res.Status, nil)
		return
	}
	queryParams := r.URL.Query()
	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(queryParams.Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}
	pageReq := &model.PostsPageReq{
		Username: username,
		Limit:    limit,
		Offset:   offset,
	}
	posts, err := h.serv.GetAllPostService(r.Context(), pageReq)
	if err != nil {
		res := helper.InternalErr("Bad request: unable to fetch posts", nil)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	res := helper.Response{
		Status: http.StatusOK,
		Message: "posts fetched",
		Data: posts,
		Err: nil,
	}
	helper.JSONResponse(w, res.Status, res)
}
func(h *PostHandler)UpdatePost(w http.ResponseWriter, r *http.Request, p router.Params){
	ctx := r.Context()
	userCtx, ok := ctx.Value(middleware.UserContextKey).(*middleware.ContextKey)
	if !ok {
		res := helper.UnauthorizedErr("Unauthorized: ", nil)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	postID, err := uuid.Parse(p.ByName("post_id"))
	if err != nil {
		res := helper.BadRequestErr("Invalid ID: post ID parsing failed", err)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	username := p.ByName("username")
	if username != userCtx.UsernameKey {
		res := helper.ForbiddenErr("Forbidden access: ", nil)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	var input model.UpdatePostInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		res := helper.BadRequestErr("Bad request: ", err)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	getPost := model.GetPostInput{
		PostID: postID,
		Owner: username,
	}
	updatedPost, err := h.serv.UpdatePostService(ctx, &input, &getPost)
	if err != nil {
		res := helper.Response{
			Status: http.StatusInternalServerError,
			Message: "invalid id",
			Data: &updatedPost,
			Err: err.Error(),
		}
		helper.JSONResponse(w, res.Status, res)
		return
	}
	res := helper.Response{
		Status: http.StatusOK,
		Message: "post updated",
		Data: &updatedPost,
		Err: nil,
	}
	helper.JSONResponse(w, res.Status, res)
}
func(h *PostHandler)DeletePost(w http.ResponseWriter, r *http.Request, p router.Params) { 
    ctx := r.Context()
    userCtx, ok := ctx.Value(middleware.UserContextKey).(*middleware.ContextKey)
    if !ok {
        res := helper.UnauthorizedErr("Unauthorized: ", nil)
        helper.JSONResponse(w, res.Status, res)
        return
    }

    postID, err := uuid.Parse(p.ByName("post_id"))
    if err != nil {
        res := helper.BadRequestErr("Invalid ID: ", err)
        helper.JSONResponse(w, res.Status, res)
        return
    }
    username := p.ByName("username")
    if username != userCtx.UsernameKey {
        res := helper.ForbiddenErr("Forbidden access: ", nil)
        helper.JSONResponse(w, res.Status, res)
        return
    }
	input := model.GetPostInput{
		PostID: postID,
		Owner: username,
	}
    err = h.serv.DeletePostService(ctx, &input)
    if err != nil {
        res := helper.InternalErr("Failed to delete post: ", err)
        helper.JSONResponse(w, res.Status, res)
        return
    }

    res := helper.Response{
        Status:  http.StatusOK,
        Message: "post deleted successfully",
        Data:    nil,
        Err:     nil,
    }
    helper.JSONResponse(w, res.Status, res)
}