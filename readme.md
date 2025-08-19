# Core Banking API (Go)

## Project Structure

### 1. Controller
Responsible for exposing REST endpoints for client and transaction operations.

- Receives HTTP requests and forwards them to the corresponding Service.
- Returns Response DTOs or encapsulated error messages.

**Endpoints:**

| Method | Path | Description |
|--------|------|-------------|
| POST   | /api/accounts | Create account |
| GET    | /api/accounts/{accountId} | Search account |
| GET    | /api/accounts/balance | Return balance |
| POST   | /api/accounts/overdraft | Set overdraft |
| POST   | /api/accounts/reset | Reset Data |
| POST   | /api/transactions | Create transaction |
| POST   | /api/transactions/event | Handle event to operate |
| GET    | /api/transactions/{transactionId} | Search transaction |
| GET    | /api/transactions/today | List transactions of the day |
| GET    | /api/transactions/range | List transactions in a date range |
| GET    | /api/transactions/type/{operationTypeId} | List transactions by type |

---

### 2. Service
Contains the core business logic.

- Receives requests from the controller, applies validations, and orchestrates repository calls.

**Examples:**

- `TransactionService` → manages creation and retrieval of transactions.
- `AccountService` → manages clients and accounts.

---

### 3. Repository
Responsible for data persistence (using Go data structures or a database).

- Abstracts queries and allows operations such as:
  - `FindById`
  - `FindByDate`
  - `FindByAccountId`
  - `Save(transaction)`

---

### 4. Request / Response DTOs
- **Request DTO:** Input data for a request.  
  Example: `TransactionRequest` contains `accountId`, `transactionType`, `amount`.
- **Response DTO:** Output data returned by the API.  
  Example: `TransactionResponse` returns `id`, `transactionType`, `amount`, `createdAt`.

---

## Business Rules

- **Single account per client:** Each client has a unique account linked to their data.
- **Transaction association:** Every operation performed creates a transaction linked to the respective account.

**Transaction Types:**

| Type | ID |
|------|----|
| Normal purchase | 1 |
| Installment purchase | 2 |
| Withdrawal | 3 |
| Credit voucher | 4 |

**Transaction Values:**

- Normal purchases and withdrawals → negative values
- Credit vouchers → positive values
- Installment purchases → treated as negative but linked to installments

**Daily Transaction Control:**

- The system allows querying transactions for the current day.

---

### Validations

- Prevent creating transactions with invalid types.
- Prevent operations on non-existent accounts.
- Ensure sufficient balance for withdrawals.

---

## Log Flow (Go Implementation)

- The log file is located at: log/transactions.log
- To read logs:

```bash
  cat log/transactions.log

# Explanation logic of logs: 

1. Controller receives the HTTP request.
2. If an error occurs (e.g., invalid method), the controller calls `utils.HandleHTTPError`, passing the `ErrorWorker` as logger.
3. `ErrorWorker` formats the log and sends it to `LogChannel`.
4. `LogChannel` worker writes the message to the log file asynchronously.

## **Flow:**

    [HTTP Request]
    ↓
    [Controller] → validation / service
    ↓
    error? → utils.HandleHTTPError
    ↓
    [ErrorWorker.Handle] → formats log
    ↓
    [LogChannel.Send] → sends to channel
    ↓
    [LogChannel.StartWorker] → writes to log file


## Go Commands to remember:

```bash
# Initialize Go module
go mod init

# Download dependencies and organize go.sum
go mod tidy

# Build executable binary
go build -o corebanking

# Run the project
go run main.go

# Run unit tests the project
go test ./test/<name_file_test>

## Prerequisites
- Go 1.20+
- Terminal / curl


## Quick Setup for use API run in Setup:

```bash
# Initialize module
go mod init corebanking

# Download dependencies
go mod tidy

# Build executable
go build -o corebanking

# Run the server
go run main.go


##  Routine to search 
 - 1. Reset (opcional)

Endpoint: POST /api/accounts/reset

description: clear previous state, to start with zero.

- 2. Create account

Endpoint: POST /api/accounts

input: { "documentNumber": "12345678900" }

output: accountId gerado.
👉 Este accountId será usado em todos os próximos passos.

- 3. Definir limite de cheque especial (opcional)

Endpoint: POST /api/accounts/overdraft

input: { "accountId": "123", "limit": 500.00 }

description: configura limite de crédito da conta.

- 4. Consultar dados da conta

Endpoint: GET /api/accounts/{accountId}

returned: accountId, documentNumber.
👉 Serve para validar que a conta foi criada corretamente.

- 5. Consultar saldo da conta

Endpoint: GET /api/accounts/balance?account_id={accountId}

returned: { "balance": 1500.00 } (por exemplo).

- 6. Criar uma transação

Existem dois jeitos diferentes na sua coleção:

 - 6.1 Transação direta

Endpoint: POST /api/transactions

input: { "accountId": "123", "operationTypeId": 1, "amount": 100.00 }

output: detalhes da transação.

 - 6.2 Evento de transação (estilo Pix/depósito)

Endpoint: POST /api/transactions/event

input: { "type": "deposit", "destination": "123", "amount": 50.00 }

output: nova versão do saldo da conta.


- 7. Consultar transações

description:  Depois de registrar transações, é possível buscá-las de várias formas:

Por ID: GET /api/transactions/{transactionId}

Do dia: GET /api/transactions/today

Por período: GET /api/transactions/range?begin=...&end=...

Por tipo: GET /api/transactions/type/{operationTypeId}


##Flow to search

> Reset Accounts (opcional)

> Create Account

> Set Overdraft (opcional)

> Get Account by ID

> Get Balance

> Create Transaction ou Handle Transaction Event

> Get Transaction(s) (por ID, day, interval, type)