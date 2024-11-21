package helper

import (
	"encoding/json"
	"net/http"
)


type Response struct {
	Status  int				`json:"status"`
	Message string			`json:"message"`
	Data    interface{}		`json:"data"`
	Err		interface{}		`json:"err"`
}

func JSONResponse(w http.ResponseWriter, status int, data interface{}){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}