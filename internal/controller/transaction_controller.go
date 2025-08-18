package controller

import (
	"corebanking/internal/dto"
	"corebanking/internal/service"
	"corebanking/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TransactionController estrutura para lidar com transações
type TransactionController struct {
	Service      *service.TransactionService
	ErrorHandler utils.ErrorHandler
}

// NewTransactionController cria um novo TransactionController
func NewTransactionController(service *service.TransactionService, errHandler utils.ErrorHandler) *TransactionController {
	return &TransactionController{Service: service, ErrorHandler: errHandler}
}

// RegisterRoutes registra as rotas no mux
func (c *TransactionController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/transactions/", c.RouteTransaction)
	mux.HandleFunc("/api/transactions/today", c.GetTransactionsToday)
	mux.HandleFunc("/api/transactions/range", c.GetTransactionsInRange)
	mux.HandleFunc("/api/transactions/type/", c.GetTransactionsByType)
	mux.HandleFunc("/api/transactions/event", c.HandleTransactionEvent)
	mux.HandleFunc("/api/transactions", c.CreateTransaction)
}

// routeTransaction trata GET /api/transactions/{transactionId}
func (c *TransactionController) RouteTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleHTTPError(w, nil, "Failed to instance method RESTful.", c.ErrorHandler)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 { // ["api","transactions","{transactionId}"]
		utils.HandleHTTPError(w, nil, "Failed to split path.", c.ErrorHandler)
		return
	}

	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		utils.HandleHTTPError(w, nil, "Failed to parse string to int.", c.ErrorHandler)
		return
	}

	c.GetTransactionByID(w, r, id)
}

// CreateTransaction cria uma nova transação
func (c *TransactionController) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleHTTPError(w, nil, "Failed to instance method RESTful.", c.ErrorHandler)
		return
	}

	var req dto.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleHTTPError(w, nil, "Invalid request body.", c.ErrorHandler)
		return
	}

	transaction, err := c.Service.CreateTransaction(&req)
	if err != nil {
		utils.HandleHTTPError(w, nil, "Failed to create request body.", c.ErrorHandler)
		return
	}

	respondJSON(w, http.StatusCreated, transaction)
}

// HandleTransactionEvent processa eventos de depósito, saque ou transferência
func (c *TransactionController) HandleTransactionEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleHTTPError(w, nil, "Failed to instance method RESTful.", c.ErrorHandler)
		return
	}

	var req dto.EventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleHTTPError(w, nil, "Invalid request body.", c.ErrorHandler)
		return
	}

	result, err := c.Service.HandleTransaction(&req)
	if err != nil {
		utils.HandleHTTPError(w, nil, "Invalid request body.", c.ErrorHandler)
		return
	}

	respondJSON(w, http.StatusCreated, result)
}

// GetTransactionByID retorna transação pelo ID
func (c *TransactionController) GetTransactionByID(w http.ResponseWriter, r *http.Request, id int64) {
	transaction, err := c.Service.GetTransactionByID(id)
	if err != nil {
		utils.HandleHTTPError(w, nil, "Failed to recovery transactionByID.", c.ErrorHandler)
		return
	}

	respondJSON(w, http.StatusOK, transaction)
}

// GetTransactionsToday retorna transações do dia
func (c *TransactionController) GetTransactionsToday(w http.ResponseWriter, r *http.Request) {
	transactions := c.Service.GetTransactionsToday()
	respondJSON(w, http.StatusOK, transactions)
}

// GetTransactionsInRange retorna transações em um intervalo de datas
func (c *TransactionController) GetTransactionsInRange(w http.ResponseWriter, r *http.Request) {
	beginStr := r.URL.Query().Get("begin")
	endStr := r.URL.Query().Get("end")

	if beginStr == "" || endStr == "" {
		utils.HandleHTTPError(w, nil, "Failed to recovery paremeters begin or end date.", c.ErrorHandler)
		return
	}

	begin, err := time.Parse(time.RFC3339, beginStr)
	if err != nil {
		utils.HandleHTTPError(w, nil, "Failed in parse begin date time format data.", c.ErrorHandler)
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		utils.HandleHTTPError(w, nil, "Failed in parse end date time format data.", c.ErrorHandler)
		return
	}

	transactions := c.Service.GetTransactionsInRange(begin, end)
	respondJSON(w, http.StatusOK, transactions)
}

// GetTransactionsByType retorna transações por tipo
func (c *TransactionController) GetTransactionsByType(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 4 { // ["api","transactions","type","{operationTypeId}"]
		utils.HandleHTTPError(w, nil, "Failed to split path.", c.ErrorHandler)
		return
	}

	typeID, err := strconv.Atoi(parts[3])
	if err != nil {
		utils.HandleHTTPError(w, nil, "Failed parse data in transactionByType.", c.ErrorHandler)
		return
	}

	transactions := c.Service.GetTransactionsByType(typeID)
	respondJSON(w, http.StatusOK, transactions)
}
