package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bagasadiii/buy-n-con/helper"
	"github.com/bagasadiii/buy-n-con/internal/middleware"
	"github.com/bagasadiii/buy-n-con/internal/model"
	"github.com/bagasadiii/buy-n-con/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	router "github.com/julienschmidt/httprouter"
)

type ItemHandlerImpl interface{
	CreateItem(w http.ResponseWriter, r *http.Request, p router.Params)
	GetItemByID(w http.ResponseWriter, r *http.Request, p router.Params)
	GetAllItems(w http.ResponseWriter, r *http.Request, p router.Params)
	UpdateItem(w http.ResponseWriter, r *http.Request, p router.Params)
	DeleteItem(w http.ResponseWriter, r *http.Request, p router.Params)
}
type ItemHandler struct {
	serv service.ItemServiceImpl
	valid *validator.Validate
}
func NewItemHandler(serv service.ItemServiceImpl)ItemHandlerImpl{
	return &ItemHandler{
		serv:serv,
		valid: validator.New(),
	}
}

func(h *ItemHandler)CreateItem(w http.ResponseWriter, r *http.Request, p router.Params){
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
	var input model.CreateItemInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		res := helper.BadRequestErr("Bad request: validation failed", err)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	if err := h.valid.Struct(&input); err != nil {
		res := helper.BadRequestErr("Bad request", err)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	item, err := h.serv.CreateItemService(ctx, &input)
	if err != nil {
		res := helper.InternalErr("Failed to create item: ", err)
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
func(h *ItemHandler)GetItemByID(w http.ResponseWriter, r *http.Request, p router.Params){
	username := p.ByName("username")
	itemID, err := uuid.Parse(p.ByName("item_id"))
	if err != nil {
		res := helper.BadRequestErr("Bad request: Invalid item ID", nil)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	input := &model.GetItemInput{
		ItemID: itemID,
		BelongsTo: username,
	}
	item, err := h.serv.GetItemByIDService(r.Context(), input)
	if err != nil {
		res := helper.InternalErr("Failed to fetch item: ", err)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	res := helper.Response{
		Status: http.StatusOK,
		Message: "OK",
		Data: &item,
		Err: nil,
	}
	helper.JSONResponse(w, res.Status, res)
}
func(h *ItemHandler)GetAllItems(w http.ResponseWriter, r *http.Request, p router.Params){
	username := p.ByName("username")
	items, err := h.serv.GetAllItemsService(r.Context(), username)
	if err != nil {
		res := helper.BadRequestErr("Bad request: unable to fetch items", nil)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	res := helper.Response{
		Status: http.StatusOK,
		Message: "Items fetched",
		Data: &items,
		Err: nil,
	}
	helper.JSONResponse(w, res.Status, res)
}
func(h *ItemHandler)UpdateItem(w http.ResponseWriter, r *http.Request, p router.Params){
	ctx := r.Context()
	userCtx, ok := ctx.Value(middleware.UserContextKey).(*middleware.ContextKey)
	if !ok {
		res := helper.UnauthorizedErr("Unauthorized: ", nil)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	itemID, err := uuid.Parse(p.ByName("item_id"))
	if err != nil {
		res := helper.BadRequestErr("Invalid ID: item ID parsing failed", err)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	username := p.ByName("username")
	if username != userCtx.UsernameKey {
		res := helper.ForbiddenErr("Forbidden access: ", nil)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	var input model.UpdateItemInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		res := helper.BadRequestErr("Bad request: ", err)
		helper.JSONResponse(w, res.Status, res)
		return
	}
	getItem := model.GetItemInput{
		ItemID: itemID,
		BelongsTo: username,
	}
	updatedItem, err := h.serv.UpdateItemService(ctx, &input, &getItem)
	if err != nil {
		log.Println(itemID)
		res := helper.Response{
			Status: http.StatusInternalServerError,
			Message: "invalid id",
			Data: updatedItem,
			Err: err.Error(),
		}
		helper.JSONResponse(w, res.Status, res)
		return
	}
	res := helper.Response{
		Status: http.StatusOK,
		Message: "Item updated",
		Data: &updatedItem,
		Err: nil,
	}
	helper.JSONResponse(w, res.Status, res)
}
func(h *ItemHandler)DeleteItem(w http.ResponseWriter, r *http.Request, p router.Params) { 
    ctx := r.Context()
    userCtx, ok := ctx.Value(middleware.UserContextKey).(*middleware.ContextKey)
    if !ok {
        res := helper.UnauthorizedErr("Unauthorized: ", nil)
        helper.JSONResponse(w, res.Status, res)
        return
    }

    itemID, err := uuid.Parse(p.ByName("item_id"))
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
    err = h.serv.DeleteItemService(ctx, &itemID)
    if err != nil {
        res := helper.InternalErr("Failed to delete item: ", err)
        helper.JSONResponse(w, res.Status, res)
        return
    }

    res := helper.Response{
        Status:  http.StatusOK,
        Message: "Item deleted successfully",
        Data:    nil,
        Err:     nil,
    }
    helper.JSONResponse(w, res.Status, res)
}