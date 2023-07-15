package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
    "strconv"
	"github.com/gorilla/mux"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
    Error string  `json:"error"`
}

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {

	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}

}

func (s *APIServer) Run() {

	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandelFunc(s.handleAccount))

	router.HandleFunc("/account/{id}", makeHTTPHandelFunc(s.handleAccountByID))

	router.HandleFunc("/transfer", makeHTTPHandelFunc(s.handleTransfer))

	log.Println("JSON API serve runnign on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)

}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {

	switch r.Method {
	case "GET":
		return s.handleGetAccount(w, r)

	case "POST":
		return s.handleCreateAccount(w, r)

    default:
		return fmt.Errorf("method not allowed Bozzo %s\n", r.Method)
	}
}

func (s *APIServer) handleAccountByID(w http.ResponseWriter, r *http.Request) error {
    switch r.Method {
    case "GET":
        return s.handleGetAccountByID(w, r)
    case "DELETE": 
        return s.handleDeleteAccount(w, r)
    default: 
        return fmt.Errorf("method not allowed %s\n", r.Method)
    }

}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
        id, err := getID(r)
        if err != nil {
            return err
        }
        acc, err := s.store.GetAccountByID(id)
        if err != nil {
            return err
        }

        return WriteJSON(w, http.StatusOK, acc)
 
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
    id, err := getID(r) 

    if err != nil {
        return err
    }

    if err := s.store.DeleteAccount(id); err != nil {
        return err
    }
    return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := &CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
    transferReq := &TransferRequest{}

    if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
        return err
    }

    defer r.Body.Close()


	return WriteJSON(w, http.StatusOK, transferReq)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandelFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
			//Error Handling :/
		}
	}
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return 0, fmt.Errorf("invalid id given %s", idStr)
    }
    return id, nil
 
}
