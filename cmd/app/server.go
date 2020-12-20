package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/trojan-t/crud/pkg/customers"
	"github.com/trojan-t/crud/pkg/security"
	"golang.org/x/crypto/bcrypt"
)

// Predefined types
var (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
)

// Server is struct
type Server struct {
	mux          *mux.Router
	customersSvc *customers.Service
	securitySvc  *security.Service
}

// NewServer is function
func NewServer(mux *mux.Router, customersSvc *customers.Service, securitySvc *security.Service) *Server {
	return &Server{mux: mux, customersSvc: customersSvc, securitySvc: securitySvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

// Init is method
func (s *Server) Init() {
	s.mux.HandleFunc("/customers", s.handleGetAllCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/active", s.handleGetAllActiveCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomerByID).Methods(GET)
	s.mux.HandleFunc("/customers/{id}", s.handleRemoveCustomerByID).Methods(DELETE)
	s.mux.HandleFunc("/customers/{id}/block", s.handleUnblockCustomerByID).Methods(DELETE)
	s.mux.HandleFunc("/customers/{id}/block", s.handleBlockCustomerByID).Methods(POST)

	s.mux.HandleFunc("/api/customers", s.handleSaveCustomer).Methods(POST)
	s.mux.HandleFunc("/api/customers/token", s.handleGenerateToken).Methods(POST)
	s.mux.HandleFunc("/api/customers/token/validate", s.handleValidateToken).Methods(POST)
}

// handleGetCustomerByID is method
func (s *Server) handleGetCustomerByID(writer http.ResponseWriter, request *http.Request) {
	idParam := mux.Vars(request)["id"]
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	item, err := s.customersSvc.ByID(request.Context(), id)
	log.Println(item)
	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(writer, http.StatusNotFound, err)
		return
	}

	if err != nil {
		log.Println(err)
		errorWriter(writer, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(writer, item)
}

func (s *Server) handleGetAllCustomers(writer http.ResponseWriter, request *http.Request) {
	items, err := s.customersSvc.All(request.Context())

	if err != nil {
		errorWriter(writer, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(writer, items)
}

func (s *Server) handleGetAllActiveCustomers(writer http.ResponseWriter, request *http.Request) {
	items, err := s.customersSvc.AllActive(request.Context())

	if err != nil {
		errorWriter(writer, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(writer, items)
}

func (s *Server) handleSaveCustomer(writer http.ResponseWriter, request *http.Request) {
	var item *customers.Customer

	if err := json.NewDecoder(request.Body).Decode(&item); err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	if err != nil {
		errorWriter(writer, http.StatusInternalServerError, err)
		return
	}
	item.Password = string(hashed)

	customer, err := s.customersSvc.Save(request.Context(), item)
	if err != nil {
		errorWriter(writer, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(writer, customer)
}

func (s *Server) handleRemoveCustomerByID(writer http.ResponseWriter, request *http.Request) {
	idParam := mux.Vars(request)["id"]
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	item, err := s.customersSvc.Delete(request.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(writer, http.StatusNotFound, err)
		return
	}

	if err != nil {
		errorWriter(writer, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(writer, item)
}

func (s *Server) handleBlockCustomerByID(writer http.ResponseWriter, request *http.Request) {
	idP := mux.Vars(request)["id"]
	id, err := strconv.ParseInt(idP, 10, 64)

	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	item, err := s.customersSvc.ChangeActive(request.Context(), id, false)

	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(writer, http.StatusNotFound, err)
		return
	}

	if err != nil {
		errorWriter(writer, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(writer, item)
}

func (s *Server) handleUnblockCustomerByID(writer http.ResponseWriter, request *http.Request) {
	idParam := mux.Vars(request)["id"]
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	item, err := s.customersSvc.ChangeActive(request.Context(), id, true)
	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(writer, http.StatusNotFound, err)
		return
	}

	if err != nil {
		errorWriter(writer, http.StatusInternalServerError, err)
		return
	}
	jsonResponse(writer, item)
}

func (s *Server) handleGenerateToken(w http.ResponseWriter, r *http.Request) {
	var item *struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	token, err := s.securitySvc.TokenForCustomer(r.Context(), item.Login, item.Password)

	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	jsonResponse(w, map[string]interface{}{"status": http.StatusText(http.StatusOK), "token": token})
}

func (s *Server) handleValidateToken(w http.ResponseWriter, r *http.Request) {
	var item *struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	id, err := s.securitySvc.AuthenticateCustomer(r.Context(), item.Token)

	if err != nil {
		status := http.StatusInternalServerError
		text := http.StatusText(http.StatusInternalServerError)
		if err == security.ErrNoSuchUser {
			status = http.StatusNotFound
			text = "not found"
		}
		if err == security.ErrExpireToken {
			status = http.StatusBadRequest
			text = "expired"
		}

		data, err := json.Marshal(map[string]interface{}{"status": "fail", "reason": text})
		if err != nil {
			errorWriter(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, err = w.Write(data)
		if err != nil {
			log.Print(err)
		}
		return
	}

	result := make(map[string]interface{})
	result["status"] = "ok"
	result["customerId"] = id

	data, err := json.Marshal(result)
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func errorWriter(w http.ResponseWriter, httpSts int, err error) {
	log.Print(err)
	http.Error(w, http.StatusText(httpSts), httpSts)
}

func jsonResponse(writer http.ResponseWriter, data interface{}) {
	item, err := json.Marshal(data)
	if err != nil {
		errorWriter(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(item)
	if err != nil {
		log.Println("Error write response: ", err)
	}
}
