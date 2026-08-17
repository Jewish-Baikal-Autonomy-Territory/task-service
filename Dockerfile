FROM golang:1.26.6 AS build

LABEL org.opencontainers.image.source=https://github.com/jewish-baikal-autonomy-territory/task-service

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/main ./cmd/main.go

FROM scratch AS prod

COPY --from=build /etc/passwd /etc/passwd
USER 65534:65534

WORKDIR /app
COPY --from=build /app/main .
EXPOSE 8080

ENV GOGC=100
ENV GOMEMLIMIT=512MiB
ENTRYPOINT ["./main"]
