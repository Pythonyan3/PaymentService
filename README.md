# 💰 PaymentService 💸

## 🤫 \*Some company* test task ✔️

[![Build Status](https://travis-ci.org/joemccann/dillinger.svg?branch=master)](https://travis-ci.org/joemccann/dillinger)

## Overview 🔍
🌐 REST API Web-Service for transactions management 🔒. Designed with Golang (net/http package) 🦦 and PostgreSQL (sqlx package)💽.

## Description ✨

Service allows to work with `Transaction` entity, create new, retrieve and update existed. `Transaction` includes following fields

```go
type Transaction struct {
	Id int
	UserId int
	Amount int64
	Status string
	Currency string
	UserEmail string
	CreatedAt time.Time
	UpdatedAt time.Time
}
```

🗃️ `Transaction` can has one of following statuses 🏷️:
1. `NEW`;
2. `ERROR`;
3. `SUCCESS`;
4. `FAILED`;
5. `CANCELED` (additional status).

Statuses `NEW`, `ERROR`, `SUCCESS`, `FAILED` were mentioned in the task. Status `CANCELED` added according to need to perform `Transaction` canceling (🔨removing data from DB is not the best approach I guess 🙃).

Statuses `NEW` and `ERROR` can be assigned to `Transaction` during creating (possibility to assign `NEW` status equals to `80%` and `ERROR` status - `20%`).

Statuses `SUCCESS` and `FAILED` can be assigned to a `Transaction` by requesting specific API endpoint.

As a previous two statuses, `CANCELED` status can be assigned by a specific API endpoint.

Statuses `SUCCESS` and `FAILED` are **terminal statuses** according to the task. Also `CANCELED` and `ERROR` statuses added to **terminal statuses**.

`Transaction` **CAN NOT** be updated with new status if it already has one of **terminal statuses**.

## 🏗️ Install & Run🏃

Download or copy repository:

```bash
git clone https://github.com/Pythonyan3/PaymentService.git
```

Before running application add ``.env`` file to the root of repository, similar to example below:

```bash
DB_PORT=5432
DB_USER=postgres
DB_HOST=localhost
SERVICE_PORT=8000
DB_PASSWORD=your_pass
DB_NAME=payments
DB_SSL_MODE=disable
```

Run service without docker:
```bash
# need to perform DB migrations
# for migrations management was used golang-migrate/migrate tool
migrate -path {path_to_migrations_folder} -database {postgres://{db_user}:{db_pass}@{db_host}:{db_port}/{db_name}?sslmode=disable up
# build up exectable file
go build -o ./cmd/payment/main ./cmd/payment/main.go
# run service
./cmd/payment/main
```

Using docker:

```bash
docker-compose build
docker-compose up
```

## API Doc 📚

Service allow to work with ``Transaction`` entity.

### List of API endpoints:

1. `/api/transactions/ (POST)` - creating new transaction;
2. `/api/transactions/{pk}/ (GET)` - retrieve transaction info;
3. `/api/transactions/{pk}/cancel/ (PUT/PATCH)` - update transaction status to `CANCELED`;
4. `/api/transactions/{pk}/proceed/ (PUT/PATCH)` - set transaction status to `CANCELED`;
5. `/api/users/{pk}/transactions/ (GET)` - retrieve list of user transactions;
6. `/api/users/{email}/transactions/ (GET)` - retrieve list of user transactions.

### Some examples of usage

#### `/api/transactions/` (POST) -  request body example:
```json
{
	"user_id": 1,
	"user_email": "email@mail.com",
	"amount": 100,
	"currency": "RUB"
}
```

Example of response:
```json
{
	"id": 1
	"user_id": 1,
	"user_email": "email@mail.com",
	"amount": 100,
	"currency": "RUB",
	"status": "NEW",
	"created_at": "2022-06-12T18:09:14.796895+03:00",
	"updated_at": "2022-06-12T18:09:14.796895+03:00"
}
```

#### `/api/transactions/{pk}/proceed/` (PUT/PATCH) - request body example:

```json
{
	"status": "SUCCESS"
}
```

Example of response:
```json
{
	"id": 1
	"user_id": 1,
	"user_email": "email@mail.com",
	"amount": 100,
	"currency": "RUB",
	"status": "SUCCESS",
	"created_at": "2022-06-12T18:09:14.796895+03:00",
	"updated_at": "2022-06-12T18:11:14.796895+03:00"
}
```

### Points to make service better 😎

1. 📄 Add pagination to responses of endpoints which possibly can return a lot of data (list of transactions endpoints);
2.  🤓 According to most of data retrieving operations from DB used PK or FK a good way to add some indexes to it;
3. 🧐 Usage of ORM can make work with entities in DB more simply.