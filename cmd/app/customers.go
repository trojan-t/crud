package app

import (
	"encoding/json"
	"net/http"

	"github.com/trojan-t/crud/pkg/customers"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) handleCustomerRegistration(w http.ResponseWriter, r *http.Request) {

	var item *customers.Customer

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	item.Password = string(hashed)

	customer, err := s.customersSvc.Save(r.Context(), item)

	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	JSONResponse(w, customer)
}

func (s *Server) handleCustomerGetToken(writer http.ResponseWriter, request *http.Request) {
	//обявляем структуру для запроса
	var item *struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	//извелекаем данные из запраса
	if err := json.NewDecoder(request.Body).Decode(&item); err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}
	//взываем из сервиса  securitySvc метод AuthenticateCustomer
	token, err := s.customersSvc.Token(request.Context(), item.Login, item.Password)

	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	//вызываем функцию для ответа в формате JSON
	JSONResponse(writer, map[string]interface{}{"status": "ok", "token": token})

}

func (s *Server) handleCustomerGetProducts(writer http.ResponseWriter, request *http.Request) {

	items, err := s.customersSvc.Products(request.Context())
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	JSONResponse(writer, items)

}
