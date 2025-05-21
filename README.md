# BrandScout-TestTask
Тестовое задание для компании BrandScout.

## 📝 Quote API

Мини-сервис для хранения и управления цитатами.

### 🚀 Запуск

```bash
go run ./cmd/server
```

### 📦 Зависимости

- Go 1.24+
- `gorilla/mux`
- Только стандартная библиотека + `testify` для тестов

### 📌 Конечные точки

#### Добавление цитаты
```bash
curl -X POST http://localhost:8080/quotes \
  -H "Content-Type: application/json" \
  -d '{"author":"Confucius", "quote":"Life is simple, but we insist on making it complicated."}'
```

#### Получение всех цитат
```bash
curl http://localhost:8080/quotes
```

#### Фильтрация по автору
```bash
curl http://localhost:8080/quotes?author=Confucius
```

#### Получение случайной цитаты
```bash
curl http://localhost:8080/quotes/random
```

#### Удаление цитаты
```bash
curl -X DELETE http://localhost:8080/quotes/1
```

### 🧪 Тесты

```bash
go test ./...
```

### 🗂 Структура проекта

```
cmd/server           — запуск сервера
internal/handler     — HTTP-хендлеры
internal/store       — in-memory хранилище
internal/model       — модель QuoteNote
internal/middleware  — логгирование запросов
test/api_test.go     — интеграционные тесты
```

---