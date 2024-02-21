FROM golang:1.21
#Создается рабочая директория
WORKDIR /app

#Копирование go.mod и go.sum
COPY go.mod go.sum ./

#Загрузка зависимостей
RUN go mod download

#Копирование остальных файлов проекта
COPY . .

#Сборка приложения проекта
RUN go build -o app ./cmd

#Запуск
CMD ["./app"]


  # Собирает докер образ на основе этого файла DockerFile и присваивает тег go:auth
  #sudo docker build -f Dockerfile -t go:auth ./


  #Запускает контейнер на основе созданного образа go:auth
  #sudo docker run go:auth