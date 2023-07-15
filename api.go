package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
    "strconv"
    "os"
    
	"github.com/gorilla/mux"
	jwt "github.com/golang-jwt/jwt/v4"


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

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))

	router.HandleFunc("/account/{id}", s.withJWTAuth(makeHTTPHandleFunc(s.handleAccountByID)))

	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer))
    router.HandleFunc("/login",    makeHTTPHandleFunc(s.handleLogin))
	log.Println("JSON API serve runnign on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)

}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error { 
     
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        return err
    }

    encPasswd := CreateHash(req.Password)
    acc, err := s.store.GetAccountByNumber(req.Number)
    
    if err != nil {
        return err
    }

    if acc.Password != encPasswd {
        return err
    }

    token, err := createJWT(acc)
    
    if err != nil {
        return err
    }
    return WriteJSON(w, http.StatusOK, map[string]string {"x-jwl-token":token})
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

	acc := NewAccount(createAccountReq.FirstName, createAccountReq.LastName, createAccountReq.Password)

	if err := s.store.CreateAccount(acc); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, acc)
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

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
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

func validateJWT(tokenString string) (*jwt.Token, error){
    secret := os.Getenv("JWT_SECRET")
    fmt.Println(secret)


    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secret), nil
    })

}

func deezNuts(w http.ResponseWriter) {
    WriteJSON(w, http.StatusForbidden, APIError{Error: "deez nuts"})
    return 
}

func (s *APIServer) withJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {

    return func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("x-jwt-token")
        token, err := validateJWT(tokenString)
        if err != nil {
            deezNuts(w)
            return 
        }
        if !token.Valid {
            deezNuts(w)
            return 
        }
        
        userID, err := getID(r)
        if err != nil {
            return 
        }
        
        acc, err := s.store.GetAccountByID(userID)
         
        if err != nil {
            deezNuts(w)
            return 
        }
        

        claims := token.Claims.(jwt.MapClaims)
        if acc.Number != int64(claims["accountNumber"].(float64)) {
            deezNuts(w)
            return

        }
        handlerFunc(w, r)
    }
}

func createJWT(acc *Account) (string, error) {

    claims := &jwt.MapClaims {
        "expiresAt": 15000,
        "accountNumber": acc.Number,
    }
     
    secret := os.Getenv("JWT_SECRET")
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    return token.SignedString([]byte(secret))
     
}
