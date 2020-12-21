package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/trojan-t/crud/cmd/app/middleware"
	"github.com/trojan-t/crud/pkg/types"
)

const admin = "ADMIN"

func (s *Server) handleManagerRegistration(w http.ResponseWriter, r *http.Request) {
	id, err := middleware.Authentication(r.Context())
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}
	if id == 0 {
		errorWriter(w, http.StatusForbidden, err)
		return
	}

	if !s.managersSvc.IsAdmin(r.Context(), id) {
		errorWriter(w, http.StatusForbidden, err)
		return
	}

	var regItem struct {
		ID    int64    `json:"id"`
		Name  string   `json:"name"`
		Phone string   `json:"phone"`
		Roles []string `json:"roles"`
	}

	err = json.NewDecoder(r.Body).Decode(&regItem)
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	item := &types.Manager{
		ID:    regItem.ID,
		Name:  regItem.Name,
		Phone: regItem.Phone,
	}

	for _, role := range regItem.Roles {
		if role == admin {
			item.IsAdmin = true
			break
		}
	}

	token, err := s.managersSvc.Create(r.Context(), item)
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	JSONResponse(w, map[string]interface{}{"token": token})
}

func (s *Server) handleManagerGetToken(w http.ResponseWriter, r *http.Request) {
	var manager *types.Manager
	err := json.NewDecoder(r.Body).Decode(&manager)
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	token, err := s.managersSvc.Token(r.Context(), manager.Phone, manager.Password)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	JSONResponse(w, map[string]interface{}{"token": token})
}

func (s *Server) handleManagerChangeProducts(w http.ResponseWriter, r *http.Request) {
	id, err := middleware.Authentication(r.Context())
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	if id == 0 {
		errorWriter(w, http.StatusForbidden, err)
		return
	}
	product := &types.Product{}
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	product, err = s.managersSvc.SaveProduct(r.Context(), product)
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	JSONResponse(w, product)
}

func (s *Server) handleManagerMakeSales(w http.ResponseWriter, r *http.Request) {
	id, err := middleware.Authentication(r.Context())
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	if id == 0 {
		errorWriter(w, http.StatusForbidden, err)
		return
	}

	sale := &types.Sale{}
	sale.ManagerID = id
	err = json.NewDecoder(r.Body).Decode(&sale)

	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	sale, err = s.managersSvc.MakeSale(r.Context(), sale)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	JSONResponse(w, sale)
}

func (s *Server) handleManagerGetSales(w http.ResponseWriter, r *http.Request) {
	id, err := middleware.Authentication(r.Context())
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	if id == 0 {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusForbidden, err)
		return
	}

	total, err := s.managersSvc.GetSales(r.Context(), id)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	JSONResponse(w, map[string]interface{}{"manager_id": id, "total": total})
}

func (s *Server) handleManagerGetProducts(w http.ResponseWriter, r *http.Request) {
	items, err := s.managersSvc.Products(r.Context())
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	JSONResponse(w, items)
}

func (s *Server) handleManagerRemoveProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := middleware.Authentication(r.Context())
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	if id == 0 {
		errorWriter(w, http.StatusForbidden, err)
		return
	}

	idParam, ok := mux.Vars(r)["id"]
	if !ok {
		errorWriter(w, http.StatusBadRequest, errors.New("id required"))
		return
	}
	productID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	err = s.managersSvc.RemoveProductByID(r.Context(), productID)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}
}

// handleRemoveCustomerByID is method
func (s *Server) handleManagerRemoveCustomerByID(w http.ResponseWriter, r *http.Request) {
	id, err := middleware.Authentication(r.Context())
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	if id == 0 {
		errorWriter(w, http.StatusForbidden, err)
		return
	}

	idParam, ok := mux.Vars(r)["id"]
	if !ok {
		errorWriter(w, http.StatusBadRequest, errors.New("Missing id"))
		return
	}

	customerID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	err = s.managersSvc.RemoveCustomerByID(r.Context(), customerID)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}
}

func (s *Server) handleManagerGetCustomers(w http.ResponseWriter, r *http.Request) {
	id, err := middleware.Authentication(r.Context())
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	if id == 0 {
		errorWriter(w, http.StatusForbidden, err)
		return
	}

	items, err := s.managersSvc.Customers(r.Context())
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	JSONResponse(w, items)
}

func (s *Server) handleManagerChangeCustomer(w http.ResponseWriter, r *http.Request) {
	id, err := middleware.Authentication(r.Context())
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	if id == 0 {
		errorWriter(w, http.StatusForbidden, err)
		return
	}

	customer := &types.Customer{}
	err = json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	customer, err = s.managersSvc.ChangeCustomer(r.Context(), customer)
	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	JSONResponse(w, customer)
}
