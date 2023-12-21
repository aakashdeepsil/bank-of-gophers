package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

// This function is used to start the API server
func (s *APIServer) Start() error {
	router := mux.NewRouter()

	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAuthentication(makeHTTPHandleFunc(s.handleAccountByID), s.store))
	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer))

	log.Println("Starting JSON API server on", s.listenAddress)

	return http.ListenAndServe(s.listenAddress, router)
}

// This function is used to create a new API server
func NewAPIServer(listenAddress string, store Storage) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		store:         store,
	}
}

func permissionDenied(w http.ResponseWriter) {
	log.Println("permission denied")
	WriteJSON(w, http.StatusUnauthorized, APIError{Message: "access denied"})
}

func withJWTAuthentication(handlerFunc http.HandlerFunc, store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("calling JWT authentication middleware")

		tokenString := r.Header.Get("x-jwt-token")

		log.Printf("token: %s", tokenString)

		token, err := validateJWTToken(tokenString)
		if err != nil {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		log.Println("token validated successfully")

		id := convertToInt(mux.Vars(r)["id"])
		account, err := store.GetAccountByID(id)
		if err != nil {
			log.Println("error getting account by ID")
			permissionDenied(w)
			return
		}

		log.Println("account retrieved successfully")

		claims := token.Claims.(jwt.MapClaims)

		log.Printf("claims retrieved successfully")

		if int64(claims["account_number"].(float64)) != account.Number {
			log.Println("account number mismatch")
			permissionDenied(w)
			return
		}

		log.Println("account retrieved successfully")

		handlerFunc(w, r)
	}
}

func validateJWTToken(token string) (*jwt.Token, error) {
	log.Println("validating JWT token")

	jwt_secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwt_secret), nil
	})
}

func createJWTToken(account *Account) (string, error) {
	log.Println("creating JWT token")

	jwt_secret := os.Getenv("JWT_SECRET")

	claims := jwt.MapClaims{
		"expires_at":     15000,
		"account_number": account.Number,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwt_secret))
}

// This function is used to create a custom HTTP handler function
func makeHTTPHandleFunc(fn APIHandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Message: err.Error()})
		}
	}
}

// This function is used to handle the /login endpoint
func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	log.Println("POST /login")

	if r.Method != http.MethodPost {
		return fmt.Errorf("unsupported method %s", r.Method)
	}

	req := new(LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("error decoding login request: %v", err)
	}

	account, err := s.store.GetAccountByNumber(req.Number)
	if err != nil {
		return fmt.Errorf("error getting account: %v", err)
	}

	if err := account.ValidatePassword(req.Password); err != nil {
		return fmt.Errorf("error validating password: %v", err)
	}

	log.Printf("account retrieved successfully: %v", account)

	token, err := createJWTToken(account)
	if err != nil {
		return fmt.Errorf("error creating JWT token: %v", err)
	}

	log.Printf("JWT token created successfully: %v", token)

	res := &LoginResponse{
		AccessToken: token,
		Number:      account.Number,
	}

	return WriteJSON(w, http.StatusOK, res)
}

// This function is used to handle the API endpoints
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGETAccount(w, r)
	case http.MethodPost:
		return s.handleCREATEAccount(w, r)
	case http.MethodDelete:
		return s.handleDELETEAccount(w, r)
	default:
		return fmt.Errorf("unsupported method %s", r.Method)
	}
}

// This function is used to handle the POST /account endpoint
func (s *APIServer) handleCREATEAccount(w http.ResponseWriter, r *http.Request) error {
	log.Println("POST /account")

	req := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	account, err := NewAccount(req.FirstName, req.LastName, req.Password)
	if err != nil {
		return err
	}

	log.Printf("account created successfully: %v", account)
	log.Printf("password encrypted successfully: %v", account.EncryptedPassword)

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

// This function is used to handle the GET /account endpoint
func (s *APIServer) handleGETAccount(w http.ResponseWriter, r *http.Request) error {
	log.Println("GET /account")

	accounts, err := s.store.GetAccount()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

// This function is used to handle the DELETE /account endpoint
func (s *APIServer) handleDELETEAccount(w http.ResponseWriter, r *http.Request) error {
	log.Println("DELETE /account")

	if err := s.store.DeleteAccount(); err != nil {
		return fmt.Errorf("error deleting account: %v", err)
	}

	return WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "All accounts deleted"})
}

// This function is used to handle the GET /account/{id} and DELETE /account/{id} endpoints
func (s *APIServer) handleAccountByID(w http.ResponseWriter, r *http.Request) error {
	log.Printf("Handling %s /account/{id}", r.Method)

	switch r.Method {
	case http.MethodGet:
		return s.handleGETAccountByID(w, r)
	case http.MethodDelete:
		return s.handleDELETEAccountByID(w, r)
	default:
		return fmt.Errorf("unsupported method %s", r.Method)

	}
}

// This function is used to handle the GET /account/{id} endpoint
func (s *APIServer) handleGETAccountByID(w http.ResponseWriter, r *http.Request) error {
	log.Println("GET /account/{id}")

	vars := mux.Vars(r)
	id := convertToInt(vars["id"])

	account, err := s.store.GetAccountByID(id)
	if err != nil {
		return fmt.Errorf("error getting account: %v", err)
	}

	return WriteJSON(w, http.StatusOK, account)
}

// This function is used to handle the DELETE /account/{id} endpoint
func (s *APIServer) handleDELETEAccountByID(w http.ResponseWriter, r *http.Request) error {
	log.Println("DELETE /account/{id}")

	vars := mux.Vars(r)
	id := convertToInt(vars["id"])

	if err := s.store.DeleteAccountByID(id); err != nil {
		return fmt.Errorf("error deleting account: %v", err)
	}

	return WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "Account deleted"})
}

// This function is used to handle the POST /transfer endpoint
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	log.Println("POST /transfer")

	return WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "Transfer successful"})
}

// This function is used to write JSON to the response writer
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// Function to convert string to int
func convertToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}
	return i
}
