package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/bagasadiii/buy-n-con/internal/middleware"
	"github.com/bagasadiii/buy-n-con/internal/model"
	"github.com/bagasadiii/buy-n-con/internal/service"
	"github.com/go-playground/validator/v10"
	router "github.com/julienschmidt/httprouter"
)

type UserHandlerImpl interface {
	Register(w http.ResponseWriter, r *http.Request, _ router.Params)
	Login(w http.ResponseWriter, r *http.Request, _ router.Params)
	GetUserByUsername(w http.ResponseWriter, r *http.Request, p router.Params)
}
type UserHandler struct {
	serv service.UserServiceImpl
	valid *validator.Validate
}
func NewUserHandler(serv service.UserServiceImpl)UserHandlerImpl{
	return &UserHandler{
		serv:serv,
		valid: validator.New(),
	}
}

func(h *UserHandler)Register(w http.ResponseWriter, r *http.Request, _ router.Params){
	var input model.RegisterInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		res := helper.BadRequestErr("Bad request", err)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	if err := h.valid.Struct(&input); err != nil {
		res := helper.BadRequestErr("Fill required form", err)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	user, err := h.serv.RegisterService(r.Context(), &input)
	if err != nil {
		res := helper.InternalErr("Internal server error while register", err)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	res := helper.Response{
		Status: http.StatusCreated,
		Message: "user created",
		Data: &user,
		Err: nil,
	}
	helper.JSONResponse(w, res.Status, res)
}

func(h *UserHandler)Login(w http.ResponseWriter, r *http.Request, _ router.Params){
	var input model.LoginInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		res := helper.BadRequestErr("Bad request", err)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	if err := h.valid.Struct(&input); err != nil {
		res := helper.BadRequestErr("Fill required form", err)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	token, err := h.serv.LoginService(r.Context(), &input)
	if err != nil {
		res := helper.BadRequestErr("Invalid Login", err)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	validate := middleware.ValidateToken(token)
	if validate.Err != nil {
		res := helper.InternalErr("Validation token error", validate.Err)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	res := helper.Response{
        Status: http.StatusOK,
        Message: "Login successful",
        Data: validate,
        Err: nil,
    }
    helper.JSONResponse(w, res.Status, res)
}

func(h *UserHandler)GetUserByUsername(w http.ResponseWriter, r *http.Request, p router.Params){
	username := p.ByName("username")
	user, err := h.serv.GetUserService(r.Context(), username)
	if err != nil {
		res := helper.BadRequestErr("Bad request", err)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	res := helper.Response{
		Status: http.StatusOK,
		Message: "OK",
		Data: user,
		Err: nil,
	}
    helper.JSONResponse(w, res.Status, res)
}