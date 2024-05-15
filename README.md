# KPI_Drive_test


## Структура проекта

- `cmd/dataBufer/main.go`: Основной API сервис для обработки запросов.
- `cmd/natsPublisher/main.go`: Сервис для публикации сообщений в NATS Streaming.
- `internal`: Внутренние пакеты проекта.
- `config/config.yaml`: Конфигурационный файл.

## Требования

- Docker
- Docker Compose

## Установка

1. Клонируйте репозиторий:

    ```sh
    git clone https://github.com/yourusername/KPI_Drive_test.git
    cd KPI_Drive_test
    ```

2. Создайте файл конфигурации `config/config.yaml` (если его нет):

    ```yaml
    env: "dev"
    http_server:
      address: "0.0.0.0:8001" # адрес сервера
      timeout: 4s
      idle_timeout: 60s
    nats:
      cluster_id: "test-cluster" # имя кластера
      client_id: "client-123" # имя клиента
      url: "nats://nats-streaming:4223" # ссылка для подключения к nats
    ```

## Запуск

Используйте Docker Compose для запуска всех сервисов:

```sh
docker-compose up --build
```

Это создаст и запустит три контейнера:

- `kpi_drive_app`: Основной API сервис, доступный по адресу http://localhost:8001.
- `nats_publisher`: Сервис для публикации сообщений в NATS Streaming.
- `nats_streaming`: NATS Streaming сервер.
Использование

## Автоматическая отправка фактов
При запуске контейнера `nats_publisher`, автоматически будет отправлено 10 фактов в NATS Streaming. 
Этот контейнер публикует факты, используя данные, заданные в `cmd/natsPublisher/main.go`.

## API Ручка
## Отправка данных через API:
```
POST /api/fact
Host: localhost:8001
Content-Type: multipart/form-data

period_start=2024-05-01
period_end=2024-05-31
period_key=month
indicator_to_mo_id=227373
indicator_to_mo_fact_id=0
value=1
fact_time=2024-05-31
is_plan=0
auth_user_id=40
comment=buffer KVSH-user
```

## Остановка

Для остановки и удаления всех контейнеров используйте команду ctrl+c или:
```
docker-compose down
```
