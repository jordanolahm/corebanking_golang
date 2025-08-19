package main

import (
	"corebanking/config"
	"corebanking/internal/controller"
	"corebanking/internal/event"
	"corebanking/internal/repository"
	"corebanking/internal/service"
	"corebanking/internal/worker"
	"net/http"
	"os"
)

func main() {
	cfg := config.LoadConfig()

	logChannel, err := event.NewLogChannel("log/transactions.log", 100)
	if err != nil {
		panic("Failed to initialize log channel: " + err.Error())
	}
	defer logChannel.Close()

	logChannel.Send("[INFO] Apllication has been started.")

	if err := os.MkdirAll("log", os.ModePerm); err != nil {
		panic("Failed to create log directory: " + err.Error())
	}

	errorWorker := worker.NewErrorWorker(logChannel)
	logChannel.Send("[INFO] Log worker started")

	accountRepo := repository.NewAccountRepository()
	transactionRepo := repository.NewTransactionRepository()
	logChannel.Send("[INFO] Repositories initialized")

	// Inicializar serviços
	accountService := service.NewAccountService(accountRepo)
	transactionService := service.NewTransactionService(transactionRepo, accountRepo)
	logChannel.Send("[INFO] Services initialized")

	// Inicializar controllers
	accountController := controller.NewAccountController(accountService, errorWorker)
	transactionController := controller.NewTransactionController(transactionService, errorWorker)
	logChannel.Send("[INFO] Controllers initialized")

	apiPrefix := "/api/" + cfg.Version

	// Configurar roteador HTTP
	mux := http.NewServeMux()

	// Rotas de accounts
	mux.HandleFunc(apiPrefix+"/accounts", accountController.CreateAccount)
	mux.HandleFunc(apiPrefix+"/accounts/", accountController.RouteAccount)
	mux.HandleFunc(apiPrefix+"/accounts/balance", accountController.GetBalance)
	mux.HandleFunc(apiPrefix+"/accounts/overdraft", accountController.SetOverdraft)
	mux.HandleFunc(apiPrefix+"/accounts/reset", accountController.Reset)

	// Rotas de transactions
	mux.HandleFunc(apiPrefix+"/transactions", transactionController.CreateTransaction)
	mux.HandleFunc(apiPrefix+"/transactions/event", transactionController.HandleTransactionEvent)
	mux.HandleFunc(apiPrefix+"/transactions/today", transactionController.GetTransactionsToday)
	mux.HandleFunc(apiPrefix+"/transactions/range", transactionController.GetTransactionsInRange)
	mux.HandleFunc(apiPrefix+"/transactions/type/", transactionController.GetTransactionsByType)
	mux.HandleFunc(apiPrefix+"/transactions/", transactionController.RouteTransaction)

	// Iniciar servidor
	serverAddr := ":" + cfg.Port
	logChannel.Send("[INFO] Starting server on " + serverAddr)

	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		logChannel.Send("[ERROR] Server failed to start: " + err.Error())
	}
}
