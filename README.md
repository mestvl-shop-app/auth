# Shop auth service

Auth service for Shop

### Инструкция по локальному развертыванию приложения:
1. Поднять БД с помощью - `make compose`
2. Установка goose для накатывания миграций - `make install-goose` 
3. Установка миграций - `make migrations-up`
4. Запуск приложения - `make run`

### Дополнительно
1. Rollback миграций - `make migrations-down`
2. Статус миграций - `make migrations-status`
3. Генерация сваггер документации - `make swag`
4. Установка линтера - `make install-golangci-lint`
5. Запуск линтера - `make lint` 