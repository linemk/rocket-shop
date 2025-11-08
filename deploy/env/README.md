# Управление переменными окружения

Эта директория содержит шаблоны и скрипты для управления переменными окружения всех микросервисов.

## Структура

```
deploy/env/
├── .env.template            # Шаблон основного файла (коммитится в git)
├── .env                     # Основной файл со всеми переменными (генерируется локально)
├── generate-env.sh          # Скрипт для генерации .env файлов
├── order.env.template       # Шаблон для Order сервиса
├── inventory.env.template   # Шаблон для Inventory сервиса
└── payment.env.template     # Шаблон для Payment сервиса
```

## Как это работает

### 1. Основной файл `.env`

Содержит все переменные окружения для всех сервисов с префиксами для избежания конфликтов:
- `ORDER_*` - переменные для Order сервиса
- `INVENTORY_*` - переменные для Inventory сервиса
- `PAYMENT_*` - переменные для Payment сервиса

### 2. Шаблоны `.env.template`

Каждый сервис имеет свой шаблон, который использует переменные из основного `.env`:

```bash
# Пример из order.env.template
ORDER_HTTP_HOST=${ORDER_HTTP_HOST}
ORDER_HTTP_PORT=${ORDER_HTTP_PORT}
```

### 3. Генерация `.env` файлов

Скрипт `generate-env.sh` обрабатывает шаблоны через `envsubst` и создает индивидуальные `.env` файлы для каждого сервиса в директориях `deploy/compose/*/`.

## Использование

### Автоматическая генерация (рекомендуется)

```bash
task env:generate
```

Эта команда:
1. Проверяет наличие `envsubst`
2. Загружает переменные из `deploy/env/.env`
3. Обрабатывает все шаблоны
4. Создает `.env` файлы в `deploy/compose/{service}/.env`

### Ручная генерация

```bash
cd deploy/env
./generate-env.sh
```

### Генерация для конкретного сервиса

```bash
export SERVICES="order"
./generate-env.sh
```

## Добавление новых переменных

1. Добавьте переменную с префиксом в `deploy/env/.env`:
   ```bash
   ORDER_NEW_VARIABLE=value
   ```

2. Добавьте использование в соответствующий шаблон:
   ```bash
   # В order.env.template
   NEW_VARIABLE=${ORDER_NEW_VARIABLE}
   ```

3. Перегенерируйте файлы:
   ```bash
   task env:generate
   ```

## Примечания

- Сгенерированные `.env` файлы в `deploy/compose/*/` добавлены в `.gitignore`
- Основной файл `deploy/env/.env` добавлен в `.gitignore` (не коммитится)
- Шаблон `deploy/env/.env.template` **ДОЛЖЕН** быть закоммичен в репозиторий
- Для локальной разработки скопируйте: `cp deploy/env/.env.template deploy/env/.env`
- В CI окружении копирование из template происходит автоматически
- Для production окружения создайте отдельный `.env.prod` файл с production значениями
- При запуске `task db:up` генерация происходит автоматически