package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/bagasadiii/buy-n-con/internal/middleware"
	"github.com/bagasadiii/buy-n-con/internal/model"
	"github.com/bagasadiii/buy-n-con/internal/service"
)

type UserHandlerImpl interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}
type UserHandler struct {
	serv service.UserServiceImpl
}
func NewUserHandler(serv service.UserServiceImpl)UserHandlerImpl{
	return &UserHandler{serv:serv}
}

func(h *UserHandler)Register(w http.ResponseWriter, r *http.Request){
	var input model.RegisterInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		res := helper.BadRequestErr("Bad request", err)
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

func(h *UserHandler)Login(w http.ResponseWriter, r *http.Request){
	var input model.LoginInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		res := helper.BadRequestErr("Bad request", err)
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
        Data: map[string]string{
            "token": token,
            "username": validate.Username,
            "user_id": validate.ID.String(),
        },
        Err: nil,
    }
    helper.JSONResponse(w, res.Status, res)
}