# Rocket Shop - Микросервисная платформа для продажи космических кораблей 🚀

Проект представляет собой распределенную систему на Go с event-driven архитектурой на базе Apache Kafka.

## Оглавление

- [Установка](#установка)
- [Запуск системы](#запуск-системы)
- [Тестовые сценарии](#тестовые-сценарии)
- [API Reference](#api-reference)
- [Архитектура](#архитектура)

---

## Установка

### Требования

- Go 1.21+
- Docker & Docker Compose
- Task CLI
- PostgreSQL
- MongoDB
- Apache Kafka

### Установка Task CLI

```bash
brew install go-task
```

---

## Запуск системы

### 1. Запуск инфраструктуры

```bash
# Создать Docker сеть
docker network create rocket-shop-network

# Запустить Kafka
docker-compose -f deploy/compose/core/docker-compose.yml up -d

# Запустить базы данных
task db:up

# Заполнить БД тестовыми данными
task db:seed
```

### 2. Запуск микросервисов

```bash
# Запустить все сервисы
task services:start:inventory
task services:start:payment
task services:start:order
task services:start:assembly
task services:start:notification
```

### 3. Остановка системы

```bash
# Остановить сервисы
task services:stop

# Остановить базы данных
task db:down

# Остановить Kafka
docker-compose -f deploy/compose/core/docker-compose.yml down
```

---

## Тестовые сценарии

### Сценарий 1: Создание заказа

#### HTTP Request (curl)

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "user-123",
    "partUUIDs": [
      "00000000-0000-0000-0000-000000000001",
      "00000000-0000-0000-0000-000000000002"
    ]
  }'
```

#### HTTP Request (Postman)

```
POST http://localhost:8080/api/v1/orders
Headers:
  Content-Type: application/json

Body (JSON):
{
  "userId": "user-123",
  "partUUIDs": [
    "00000000-0000-0000-0000-000000000001",
    "00000000-0000-0000-0000-000000000002"
  ]
}
```

#### Ожидаемый Response

```json
{
  "orderUuid": "851bc3b0-a4c7-43d5-a557-33473b33747b"
}
```

**Статус:** `201 Created`

#### Что произойдет

**В базе данных (PostgreSQL):**
- Создается запись в таблице `orders` со статусом `PENDING_PAYMENT`
- Сохраняется информация: `user_id`, `total_price`, `created_at`

**В Kafka:**
- Пока ничего (события отправляются только при оплате)

**В Telegram:**
- Пока ничего

---

### Сценарий 2: Оплата заказа

#### HTTP Request (curl)

```bash
ORDER_UUID="851bc3b0-a4c7-43d5-a557-33473b33747b"  # UUID из предыдущего шага

curl -X POST "http://localhost:8080/api/v1/orders/${ORDER_UUID}/pay" \
  -H "Content-Type: application/json" \
  -d '{
    "paymentMethod": "PAYMENT_METHOD_CARD"
  }'
```

#### HTTP Request (Postman)

```
POST http://localhost:8080/api/v1/orders/{{ORDER_UUID}}/pay
Headers:
  Content-Type: application/json

Body (JSON):
{
  "paymentMethod": "PAYMENT_METHOD_CARD"
}
```

**Доступные методы оплаты:**
- `PAYMENT_METHOD_CARD` - банковская карта
- `PAYMENT_METHOD_CASH` - наличные
- `PAYMENT_METHOD_CRYPTO` - криптовалюта

#### Ожидаемый Response

```json
{
  "transactionUuid": "47d0b01e-ca98-432d-b4c1-9e1c1bdc3614"
}
```

**Статус:** `200 OK`

#### Что произойдет

**В базе данных (PostgreSQL):**
- Статус заказа обновляется: `PENDING_PAYMENT` → `PAID`
- Сохраняется `transaction_uuid` и `payment_method`
- Обновляется `updated_at`

**В Kafka:**
1. **Order Service** отправляет событие `OrderPaid` в топик `order-paid`:
   ```json
   {
     "eventUuid": "e1db47c2-3f35-4abf-83d9-d199f531c309",
     "orderUuid": "851bc3b0-a4c7-43d5-a557-33473b33747b",
     "userUuid": "user-123",
     "paymentMethod": "PAYMENT_METHOD_CARD",
     "transactionUuid": "47d0b01e-ca98-432d-b4c1-9e1c1bdc3614"
   }
   ```

2. **Assembly Service** читает событие из `order-paid` и начинает сборку корабля (2-10 секунд)

3. **Assembly Service** отправляет событие `ShipAssembled` в топик `ship-assembled`:
   ```json
   {
     "eventUuid": "0bf809b7-35c6-4d7f-95ca-b85249cfd6bd",
     "orderUuid": "851bc3b0-a4c7-43d5-a557-33473b33747b",
     "userUuid": "user-123",
     "buildTimeSec": "5"
   }
   ```

4. **Order Service** читает событие из `ship-assembled` и обновляет статус заказа: `PAID` → `ASSEMBLED`

**В Telegram (приходят 2 сообщения):**

1. **Сообщение об оплате** (сразу после оплаты):
   ```
   💳 Платеж успешно обработан

   Информация о платеже:
   • Заказ: 851bc3b0-a4c7-43d5-a557-33473b33747b
   • Пользователь: user-123
   • Метод оплаты: PAYMENT_METHOD_CARD
   • Транзакция: 47d0b01e-ca98-432d-b4c1-9e1c1bdc3614

   Спасибо за вашу покупку!
   ```

2. **Сообщение о сборке** (через 2-10 секунд):
   ```
   🚀 Ваш заказ собран!

   Информация о доставке:
   • Заказ: 851bc3b0-a4c7-43d5-a557-33473b33747b
   • Пользователь: user-123
   • Время сборки: 5 сек

   Ваш заказ готов к доставке!
   ```

---

### Сценарий 3: Получение информации о заказе

#### HTTP Request (curl)

```bash
ORDER_UUID="851bc3b0-a4c7-43d5-a557-33473b33747b"

curl -X GET "http://localhost:8080/api/v1/orders/${ORDER_UUID}"
```

#### HTTP Request (Postman)

```
GET http://localhost:8080/api/v1/orders/{{ORDER_UUID}}
```

#### Ожидаемый Response (после оплаты и сборки)

```json
{
  "uuid": "851bc3b0-a4c7-43d5-a557-33473b33747b",
  "userId": "user-123",
  "totalPrice": 150000.00,
  "status": "ASSEMBLED",
  "paymentMethod": "PAYMENT_METHOD_CARD",
  "transactionUuid": "47d0b01e-ca98-432d-b4c1-9e1c1bdc3614",
  "createdAt": "2025-11-18T17:58:10Z",
  "updatedAt": "2025-11-18T17:58:20Z"
}
```

**Статус:** `200 OK`

**Возможные статусы заказа:**
- `PENDING_PAYMENT` - ожидает оплаты
- `PAID` - оплачен, но еще не собран
- `ASSEMBLED` - собран и готов к доставке
- `CANCELLED` - отменен

---

### Сценарий 4: Отмена заказа

#### HTTP Request (curl)

```bash
ORDER_UUID="851bc3b0-a4c7-43d5-a557-33473b33747b"

curl -X DELETE "http://localhost:8080/api/v1/orders/${ORDER_UUID}"
```

#### HTTP Request (Postman)

```
DELETE http://localhost:8080/api/v1/orders/{{ORDER_UUID}}
```

#### Ожидаемый Response

**Статус:** `204 No Content`

#### Что произойдет

**В базе данных (PostgreSQL):**
- Статус заказа обновляется на `CANCELLED`
- Обновляется `updated_at`

**В Kafka:**
- Пока ничего (в будущем можно добавить событие `OrderCancelled`)

**В Telegram:**
- Пока ничего

**Примечание:** Отменить можно только заказ в статусе `PENDING_PAYMENT`. Оплаченные заказы отменить нельзя.

---

## API Reference

### Orders API

| Метод | Endpoint | Описание |
|-------|----------|----------|
| POST | `/api/v1/orders` | Создать новый заказ |
| GET | `/api/v1/orders/{uuid}` | Получить информацию о заказе |
| POST | `/api/v1/orders/{uuid}/pay` | Оплатить заказ |
| DELETE | `/api/v1/orders/{uuid}` | Отменить заказ |

### Коды ответов

| Код | Описание |
|-----|----------|
| 200 | OK - Запрос выполнен успешно |
| 201 | Created - Ресурс создан |
| 204 | No Content - Запрос выполнен, тело ответа пустое |
| 400 | Bad Request - Некорректный запрос |
| 404 | Not Found - Ресурс не найден |
| 409 | Conflict - Конфликт (например, заказ уже оплачен) |
| 500 | Internal Server Error - Внутренняя ошибка сервера |

---

## Архитектура

### Микросервисы

1. **Order Service** (HTTP API + Kafka Producer + Kafka Consumer)
   - Управление заказами (CRUD)
   - Отправка событий `OrderPaid` в Kafka
   - Прием событий `ShipAssembled` из Kafka
   - База данных: PostgreSQL

2. **Assembly Service** (Kafka Consumer + Kafka Producer)
   - Симуляция сборки кораблей
   - Прием событий `OrderPaid` из Kafka
   - Отправка событий `ShipAssembled` в Kafka

3. **Notification Service** (Kafka Consumer)
   - Отправка уведомлений в Telegram
   - Прием событий `OrderPaid` и `ShipAssembled` из Kafka

4. **Payment Service** (gRPC Server)
   - Обработка платежей
   - Генерация transaction UUID

5. **Inventory Service** (gRPC Server)
   - Управление складом запчастей
   - База данных: MongoDB

### Event Flow

```
HTTP API → Order Service → order-paid → Assembly Service
    → ship-assembled → [Notification Service, Order Service]
    → Telegram
```

### Kafka Topics

- `order-paid` - события оплаты заказов (3 партиции)
- `ship-assembled` - события сборки кораблей (3 партиции)

### Consumer Groups

- `assembly-service` - читает из `order-paid`
- `notification-service-paid` - читает из `order-paid`
- `notification-service-assembled` - читает из `ship-assembled`
- `order-service` - читает из `ship-assembled`

---

## CI/CD

Проект использует GitHub Actions для непрерывной интеграции и доставки. Основные workflow:

- **CI** (`.github/workflows/ci.yml`) - проверяет код при каждом push и pull request
  - Линтинг кода
  - Проверка безопасности
  - Выполняется автоматическое извлечение версий из Taskfile.yml

---

## Troubleshooting

### Kafka не запускается

```bash
# Проверить логи Kafka
docker logs kafka

# Пересоздать контейнер
docker-compose -f deploy/compose/core/docker-compose.yml down
docker-compose -f deploy/compose/core/docker-compose.yml up -d
```

### Сервисы не могут подключиться к Kafka

```bash
# Проверить что Kafka доступен
docker ps | grep kafka

# Проверить что network создана
docker network ls | grep rocket-shop-network

# Если нет - создать
docker network create rocket-shop-network
```

### Telegram уведомления не приходят

1. Проверить логи notification service: `tail -f /tmp/notification.log`
2. Проверить что TELEGRAM_BOT_TOKEN и TELEGRAM_BOT_CHAT_ID правильные
3. Проверить что бот добавлен в чат
