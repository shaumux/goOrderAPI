package controllers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/gofrs/uuid"
	"goOrderAPI/models"
	"net/http"
	"strconv"
)

type Order struct {
	*models.Order
}

type OrderList struct {
	ORDERS []*models.Order
}

func (o *OrderList) MarshalJSON() ([]byte, error) {
	orders := []*Order{}
	for _, order := range o.ORDERS {
		ordr := Order{order}
		orders = append(orders, &ordr)
	}
	return json.Marshal(orders)
}

func (o *OrderList) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	return nil
}

type readOrder struct {
	ID       uuid.UUID
	Status   string
	Distance int
}

func newReadOrder(order *Order) *readOrder {
	return &readOrder{
		order.ID,
		order.Status,
		order.Distance,
	}
}

func (o *Order) MarshalJSON() ([]byte, error) {
	return json.Marshal(newReadOrder(o))
}

func (o *Order) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	return nil
}

func (o *Order) Bind(r *http.Request) error {
	if err := o.Validate(); err != nil {
		return err
	}
	return nil
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	order := &Order{}
	if err := render.Bind(r, order); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	if _, err := order.Create(); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	render.Render(w, r, order)
}

func TakeOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	lookUpID, err := uuid.FromString(id)
	if err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	lookupOrder := &Order{}
	lookupOrder.Order = &models.Order{}
	lookupOrder.ID = lookUpID
	orderInterface, err := lookupOrder.Get(true)
	if err != nil {
		render.Render(w, r, ErrNotFound)
		return
	}

	//order := orderInterface.(models.Order)
	lookupOrder.Order = orderInterface.(*models.Order)
	//if err:=render.Bind(r, updateOrder); err!=nil{
	//	render.Render(w, r, ErrBadRequest(err))
	//	return
	//}
	type Status struct {
		Status string
	}
	status := &Status{}
	json.NewDecoder(r.Body).Decode(&status)
	if status.Status != "TAKEN" {
		render.Render(w, r, ErrBadRequest(errors.New("Invalid status")))
		return
	}
	if lookupOrder.Status != "TAKEN" {
		if _, err := lookupOrder.Update(map[string]interface{}{"status": status.Status}); err != nil {
			render.Render(w, r, ErrBadRequest(err))
			return
		}
	} else {
		render.Render(w, r, ErrBadRequest(errors.New("Order Already Taken")))
		return
	}

	render.Render(w, r, RequestSuccessfull)
}

func ListOrders(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		render.Render(w, r, ErrBadRequest(errors.New("Invalid Limit")))
		return
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		render.Render(w, r, ErrBadRequest(errors.New("Invalid page")))
		return
	}
	ordersInterface, err := (&Order{}).All(page, limit)
	if err != nil {
		render.Render(w, r, ErrBadRequest(err))
	}
	orders := &OrderList{}
	orders.ORDERS = ordersInterface.([]*models.Order)
	render.Render(w, r, orders)
}
