FROM golang:1.23 AS build
WORKDIR /src
COPY ./cmd/main.go /src
RUN CGO_ENABLED=0 GOOS=linux go build -o kvs /src/cmd/main.go
# Мы используем образ "scratch", не содержащий распределяемых файлов.
# Получившийся образ будет содержать только двоичный файл службы.
FROM scratch
# Скопировать имеющийся двоичный файл с хоста.
COPY --from=build /src/kvs .
# Сообщить фреймворку Docker, что служба будет использовать порт 8080.
EXPOSE 8080
# Команда, которая должна быть выполнена при запуске контейнера.
CMD ["/kvs"]