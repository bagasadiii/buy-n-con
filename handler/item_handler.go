package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/bagasadiii/buy-n-con/internal/model"
	"github.com/bagasadiii/buy-n-con/internal/service"
)

type ItemHandlerImpl interface{
	CreateItem(w http.ResponseWriter, r *http.Request)
}
type ItemHandler struct {
	serv service.ItemServiceImpl
}
func NewItemHandler(serv service.ItemServiceImpl)ItemHandlerImpl{
	return &ItemHandler{serv:serv}
}

func(h *ItemHandler)CreateItem(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()

	var input model.CreateItemInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		res := helper.BadRequestErr("Bad request", err)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	item, err := h.serv.CreateItemService(ctx, &input)
	if err != nil {
		res := helper.InternalErr("Failed to create user", err)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	res := helper.Response{
		Status: http.StatusCreated,
		Message: "item created",
		Data: &item,
		Err: nil,
	}
	helper.JSONResponse(w, res.Status, res)
}