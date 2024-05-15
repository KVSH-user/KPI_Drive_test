# Базовый образ для сборки
FROM golang:1.21 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum файлы
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы проекта в рабочую директорию
COPY . .

# Собираем приложение для двух команд
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o dataBuf ./cmd/dataBufer/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o natsPub ./cmd/natsPublisher/main.go

# Финальный образ
FROM alpine:latest

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы конфигурации
COPY config/config.yaml ./config/

# Копируем исполняемые файлы из образа сборки
COPY --from=builder /app/dataBuf .
COPY --from=builder /app/natsPub .

# Даем права на выполнение файлов
RUN chmod +x dataBuf natsPub

# Открываем порты для сервиса
EXPOSE 8001

# Указываем команду по умолчанию
CMD ["./dataBuf"]
