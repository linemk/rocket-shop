# Docker Compose для Rocket Shop

Эта директория содержит Docker Compose конфигурации для всех сервисов проекта Rocket Shop.

## Структура

- `core/` - базовая сеть для всех сервисов
- `order/` - PostgreSQL для OrderService
- `inventory/` - MongoDB для InventoryService

## Запуск

### 1. Создать сеть

```bash
cd core
docker-compose up -d
```

### 2. Запустить PostgreSQL для OrderService

```bash
cd ../order
docker-compose up -d
```

### 3. Запустить MongoDB для InventoryService

```bash
cd ../inventory
docker-compose up -d
```

### 4. Запустить все сервисы одной командой

```bash
# Из корневой директории проекта
docker-compose -f deploy/compose/core/docker-compose.yml up -d
docker-compose -f deploy/compose/order/docker-compose.yml up -d
docker-compose -f deploy/compose/inventory/docker-compose.yml up -d
```

## Остановка

```bash
docker-compose -f deploy/compose/inventory/docker-compose.yml down
docker-compose -f deploy/compose/order/docker-compose.yml down
docker-compose -f deploy/compose/core/docker-compose.yml down
```

## Переменные окружения

Все переменные окружения находятся в файле `.env` в этой директории:

- `ORDER_POSTGRES_USER` - пользователь PostgreSQL для OrderService
- `ORDER_POSTGRES_PASSWORD` - пароль PostgreSQL для OrderService
- `ORDER_POSTGRES_DB` - имя базы данных для OrderService
- `ORDER_POSTGRES_PORT` - порт PostgreSQL

- `INVENTORY_MONGO_USER` - пользователь MongoDB для InventoryService
- `INVENTORY_MONGO_PASSWORD` - пароль MongoDB для InventoryService
- `INVENTORY_MONGO_DB` - имя базы данных для InventoryService
- `INVENTORY_MONGO_PORT` - порт MongoDB

## Подключение к базам данных

### PostgreSQL (OrderService)

```bash
psql -h localhost -p 5432 -U order_user -d order_db
```

### MongoDB (InventoryService)

```bash
mongosh "mongodb://inventory_user:inventory_password@localhost:27017/inventory_db?authSource=admin"
```

## Проверка работы

### PostgreSQL

```bash
docker exec -it order-postgres psql -U order_user -d order_db -c "\dt"
```

### MongoDB

```bash
docker exec -it inventory-mongodb mongosh --eval "db.adminCommand('ping')"
```
