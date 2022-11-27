package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/UserNaMEeman/shops/app"
)

func (h *Handler) uploadOrder(w http.ResponseWriter, r *http.Request) {
	// var order
	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	guid := fmt.Sprintf("%s", ctx.Value("guid"))
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	order := h.services.Orders
	orderGUID, free := order.CheckOrder(guid, string(body))
	if free {
		if err = order.UploadOrderNumber(guid, string(body)); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		fmt.Println("guid: ", guid, "order: ", string(body))
		w.WriteHeader(http.StatusAccepted)
		return
	}
	if orderGUID == guid {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.WriteHeader(http.StatusConflict)
		return
	}
}

func (h *Handler) getOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	guid := fmt.Sprintf("%s", ctx.Value(app.TypeGUID))
	newOrder := h.services.Orders
	orders, err := newOrder.GetOrders(guid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		fmt.Println("guid for no content: ", guid)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	data, err := json.Marshal(orders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("send to user: ", string(data))
	w.Header().Add("Content-Type", "application/json")

	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
