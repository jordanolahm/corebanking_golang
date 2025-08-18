package controller

import (
	"corebanking/internal/dto"
	"corebanking/internal/service"
	"corebanking/internal/utils"
	"encoding/json"
	"net/http"
	"strings"
)

type AccountController struct {
	Service      *service.AccountService
	ErrorHandler utils.ErrorHandler
}

func NewAccountController(service *service.AccountService, errHandler utils.ErrorHandler) *AccountController {
	return &AccountController{
		Service:      service,
		ErrorHandler: errHandler,
	}
}

// RegisterRoutes registra as rotas no mux
func (c *AccountController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/accounts/", c.RouteAccount)
	mux.HandleFunc("/api/accounts/balance", c.GetBalance)
	mux.HandleFunc("/api/accounts", c.CreateAccount)
	mux.HandleFunc("/api/accounts/overdraft", c.SetOverdraft)
	mux.HandleFunc("/api/accounts/reset", c.Reset)
}

func (c *AccountController) RouteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleHTTPError(w, nil, "Failed to instance method RESTful.", c.ErrorHandler)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 {
		utils.HandleHTTPError(w, nil, "Failed to split path.", c.ErrorHandler)
		return
	}

	accountID := parts[2]
	c.GetAccount(w, r, accountID)
}

func (c *AccountController) GetAccount(w http.ResponseWriter, r *http.Request, accountID string) {
	account, err := c.Service.GetAccount(accountID)
	if err != nil {
		utils.HandleHTTPError(w, err, "Failed to get account.", c.ErrorHandler)
		return
	}
	respondJSON(w, http.StatusOK, account)
}

func (c *AccountController) GetBalance(w http.ResponseWriter, r *http.Request) {
	accountID := r.URL.Query().Get("account_id")
	if accountID == "" {
		utils.HandleHTTPError(w, nil, "Failed to recovery account_id.", c.ErrorHandler)
		return
	}

	balance, err := c.Service.GetBalance(accountID)
	if err != nil {
		utils.HandleHTTPError(w, err, "Failed to get balance.", c.ErrorHandler)
		return
	}

	respondJSON(w, http.StatusOK, balance)
}

func (c *AccountController) CreateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleHTTPError(w, nil, "Failed to instance method RESTful.", c.ErrorHandler)
		return
	}

	var req dto.AccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleHTTPError(w, err, "Failed to decode request.", c.ErrorHandler)
		return
	}

	account, err := c.Service.CreateAccount(req.DocumentNumber)
	if err != nil {
		utils.HandleHTTPError(w, err, "Failed to create account.", c.ErrorHandler)
		return
	}

	respondJSON(w, http.StatusCreated, account)
}

func (c *AccountController) SetOverdraft(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleHTTPError(w, nil, "Failed to instance method RESTful.", c.ErrorHandler)
		return
	}

	var req dto.OverdraftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleHTTPError(w, err, "Failed to decode request.", c.ErrorHandler)
		return
	}

	if err := c.Service.ConfigOverdraft(req.AccountID, req.Limit); err != nil {
		utils.HandleHTTPError(w, err, "Failed in set new overdraft limit.", c.ErrorHandler)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *AccountController) Reset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleHTTPError(w, nil, "Failed to instance method RESTful.", c.ErrorHandler)
		return
	}

	c.Service.Reset()
	w.WriteHeader(http.StatusOK)
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
