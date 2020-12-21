package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/trojan-t/crud/cmd/app/middleware"
	"github.com/trojan-t/crud/pkg/customers"
	"github.com/trojan-t/crud/pkg/managers"
)

// Server is struct
type Server struct {
	mux          *mux.Router
	customersSvc *customers.Service
	managersSvc  *managers.Service
}

// NewServer is function
func NewServer(m *mux.Router, customersSvc *customers.Service, managersSvc *managers.Service) *Server {
	return &Server{
		mux:          m,
		customersSvc: customersSvc,
		managersSvc:  managersSvc,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// Init is function
func (s *Server) Init() {
	log.Println("start init method")
	customersAuthenticateMd := middleware.Authenticate(s.customersSvc.IDByToken)
	customersSubrouter := s.mux.PathPrefix("/api/customers").Subrouter()
	customersSubrouter.Use(customersAuthenticateMd)

	customersSubrouter.HandleFunc("", s.handleCustomerRegistration).Methods("POST")
	customersSubrouter.HandleFunc("/token", s.handleCustomerGetToken).Methods("POST")
	customersSubrouter.HandleFunc("/products", s.handleCustomerGetProducts).Methods("GET")

	managersAuthenticateMd := middleware.Authenticate(s.managersSvc.IDByToken)
	managersSubRouter := s.mux.PathPrefix("/api/managers").Subrouter()
	managersSubRouter.Use(managersAuthenticateMd)
	managersSubRouter.HandleFunc("", s.handleManagerRegistration).Methods("POST")
	managersSubRouter.HandleFunc("/token", s.handleManagerGetToken).Methods("POST")
	managersSubRouter.HandleFunc("/sales", s.handleManagerGetSales).Methods("GET")
	managersSubRouter.HandleFunc("/sales", s.handleManagerMakeSales).Methods("POST")
	managersSubRouter.HandleFunc("/products", s.handleManagerGetProducts).Methods("GET")
	managersSubRouter.HandleFunc("/products", s.handleManagerChangeProducts).Methods("POST")
	managersSubRouter.HandleFunc("/products/{id:[0-9]+}", s.handleManagerRemoveProductByID).Methods("DELETE")
	managersSubRouter.HandleFunc("/customers", s.handleManagerGetCustomers).Methods("GET")
	managersSubRouter.HandleFunc("/customers", s.handleManagerChangeCustomer).Methods("POST")
	managersSubRouter.HandleFunc("/customers/{id:[0-9]+}", s.handleManagerRemoveCustomerByID).Methods("DELETE")

}

func errorWriter(w http.ResponseWriter, httpSts int, err error) {
	log.Println("writeError: ", err)
	http.Error(w, http.StatusText(httpSts), httpSts)
}

// JSONResponse is function
func JSONResponse(w http.ResponseWriter, iData interface{}) {
	data, err := json.Marshal(iData)
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print("response write error: ", err)
	}
}
