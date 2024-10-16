FROM golang:1.22.3-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify

ARG APP_NAME
COPY ./cmd/${APP_NAME}/main.go ./cmd/${APP_NAME}/

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/${APP_NAME} ./cmd/${APP_NAME}/main.go

FROM alpine:3.20

WORKDIR /home/appuser/

ARG APP_NAME
COPY --from=builder /app/${APP_NAME} ./

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown appuser:appgroup ./${APP_NAME}

USER appuser

EXPOSE 8080
EXPOSE 9000
