# GoBank

Production-style backend сервис для имитации банковских операций: регистрация, депозиты и переводы между пользователями с использованием транзакций PostgreSQL.


## Технологии
 - Go 1.22+
 - Gin (HTTP layer)
 - PostgreSQL
 - pgx (работа без ORM)
 - JWT (golang-jwt)
 - bcrypt
 - Docker
 - Clean Architecture approach


## Возможности

 - Регистрация и авторизация пользователей
 - JWT-аутентификация
 - Хэширование паролей (bcrypt)
 - Депозит средств
 - Переводы между пользователями
 - Атомарные транзакции БД
 - Проверка достаточности средств
 - Защита от self-transfer
 - Graceful shutdown
 - Dockerized deployment


## Архитектура

Проект построен по принципам разделения ответственности:
Client -> Gin -> Middleware (JWT) -> Handler -> Usecase -> Repository -> Database


## Структура проекта
```bash
internal/
├── DTO
├── app
├── delivery/http
├── middleware
├── repository
├── usecase
├── service
├── domain
├── config
└── server
```


## Аутентификация

Используется JWT (HS256):
 - user_id хранится в claims
 - roles хранятся в claims
 - проверяется:
    подпись
    срок действия (exp)
    issuer
 - токен передаётся через заголовок:
 Authorization: Bearer <token>


## Логика переводов

Перевод средств реализован через транзакцию PostgreSQL:
```bash
tx, err := repo.BeginTransaction(ctx)
defer tx.Rollback(ctx)

repo.MinusMoneyTx(...)
repo.AddMoneyTx(...)

tx.Commit(ctx)
```
Гарантируется:
 - атомарность операции
 - отсутствие отрицательного баланса
 - защита от перевода самому себе
 - проверка достаточности средств


## База данных

Таблица Users
```bash
| Поле     | Тип     |
| -------- | ------- |
| id       | bigint  |
| name     | text    |
| email    | text    |
| password | text    |
| roles    | text[]  |
| amount   | numeric |
```

Таблица Transactions
```bash
| Поле       | Тип       |
| ---------- | --------- |
| from_user  | bigint    |
| to_user    | bigint    |
| amount     | numeric   |
| created_at | timestamp |
| status     | text      |
```


## API

POST /register
Регистрация пользователя

Пример тела запроса:
```bash
{
  "name": "Dexter",
  "email": "dexter@example.com",
  "password": "bayharbour"
}
```
Пример ответа:
```bash
{
  "answer": "<jwt_token>"
}
```

POST /login
Авторизация

Пример тела запроса:
```bash
{
  "email": "dexter@example.com",
  "password": "bayharbour"
}
```

POST /user/deposit
Пополнение баланса
(требуется JWT)

Пример тела запроса:
```bash
{
  "amount": 1000
}
```

## POST /user/transaction
Перевод средств
(требуется JWT)

Пример тела запроса:
```bash
{
  "to": 2,
  "amount": 500
}
```


## GET /user/me
Получение информации о себе


## Инженерные особенности
 - Использование pgx.Tx для ACID-транзакций
 - Проверка RowsAffected для защиты от race conditions
 - Контексты с таймаутом для DB-операций
 - Разделение слоёв (Handler → Usecase → Repository)
 - Graceful shutdown через SIGINT/SIGTERM
 - Dependency injection вручную
