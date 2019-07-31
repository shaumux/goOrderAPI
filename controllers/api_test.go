package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPlaceInvalidOrder(t *testing.T) {
	body := `{}`
	req, err := http.NewRequest("POST", "/orders", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateOrder)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Incorrect Status code Expected: %v, Got: %v", http.StatusBadRequest, rr.Code)
	}
}

func TestInvalidCordinates(t *testing.T) {
	body := `{"origin":["zxc", "sds"],"destination":["vv","sds"]}`
	req, err := http.NewRequest("POST", "/orders", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateOrder)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Incorrect Status code Expected: %v, Got: %v", http.StatusBadRequest, rr.Code)
	}
}

func TestInvalidCoordinateCombination(t *testing.T) {
	body := `{"origin":["22.372081","114.107877"],"destination":["14.5965788","120.9445402"]}`
	req, err := http.NewRequest("POST", "/orders", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateOrder)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Incorrect Status code Expected: %v, Got: %v", http.StatusBadRequest, rr.Code)
	}

}

func TestValidCreate(t *testing.T) {
	body := `{"origin":["12.9749391","77.6365496"],"destination":["12.9676997","77.6511029"]}`
	req, err := http.NewRequest("POST", "/orders", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateOrder)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Incorrect Status code Expected: %v, Got: %v", http.StatusOK, rr.Code)
	}
	resp := struct {
		Status   string    `json:"status"`
		ID       uuid.UUID `json:"id"`
		Distance int       `json:"distance"`
	}{}
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp.Status != "UNASSIGNED" {
		t.Errorf("Incorrect status Got:%v Expected:UNASSIGNED", resp.Status)
	}
}

func TestGetInvalidParamLimit(t *testing.T) {
	req, err := http.NewRequest("GET", "/orders?page=1&limit=abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListOrders)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Incorrect Status code Expected: %v, Got: %v", http.StatusBadRequest, rr.Code)
	}
}

func TestGetInvalidParamPage(t *testing.T) {
	req, err := http.NewRequest("GET", "/orders?page=a&limit=100", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListOrders)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Incorrect Status code Expected: %v, Got: %v", http.StatusBadRequest, rr.Code)
	}
}

func TestGetValidGet(t *testing.T) {
	req, err := http.NewRequest("GET", "/orders?page=1&limit=100", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListOrders)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Incorrect Status code Expected: %v, Got: %v", http.StatusOK, rr.Code)
	}
}

func TestTakeOrder(t *testing.T) {
	body := `{"origin":["12.9749391","77.6365496"],"destination":["12.9676997","77.6511029"]}`
	req, err := http.NewRequest("POST", "/orders", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateOrder)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Incorrect Status code Expected: %v, Got: %v", http.StatusOK, rr.Code)
	}
	resp := struct {
		Status   string    `json:"status"`
		ID       uuid.UUID `json:"id"`
		Distance int       `json:"distance"`
	}{}
	json.NewDecoder(rr.Body).Decode(&resp)
	takeCall := func(readChan <-chan string, resultChan chan<- string) {
		body := `{"status": "TAKEN"}`
		orderId := <-readChan
		req, err := http.NewRequest("PATCH", fmt.Sprintf("/orders/%s", orderId), strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", orderId)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler := http.HandlerFunc(TakeOrder)
		handler.ServeHTTP(rr, req)
		if rr.Code == http.StatusOK {
			resp := struct {
				Status string `json:"status"`
			}{}
			json.NewDecoder(rr.Body).Decode(&resp)
			resultChan <- resp.Status
		} else {
			resp := struct {
				Error string `json:"error"`
			}{}
			json.NewDecoder(rr.Body).Decode(&resp)
			resultChan <- resp.Error
		}
		return
	}

	readChan1 := make(chan string)
	readChan2 := make(chan string)
	readChan3 := make(chan string)
	resultChan1 := make(chan string)
	resultChan2 := make(chan string)
	resultChan3 := make(chan string)
	go takeCall(readChan1, resultChan1)
	go takeCall(readChan2, resultChan2)
	go takeCall(readChan3, resultChan3)
	readChan1 <- resp.ID.String()
	readChan2 <- resp.ID.String()
	readChan3 <- resp.ID.String()

	result1 := <-resultChan1
	result2 := <-resultChan2
	result3 := <-resultChan3

	if result1 == result2 && result2 == result3 {
		t.Errorf("Same order taken multiple times")
	}
}
