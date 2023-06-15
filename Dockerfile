FROM golang:1.19-alpine as builder

LABEL maintainer="Marvin Collins <marvincollins14@gmail.com>"
WORKDIR /app

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -ldflags="-w -s" -o axisapi .

FROM scratch
COPY --from=builder /app/axisapi /
COPY --from=builder /app/.env /

EXPOSE 80
ENTRYPOINT ["./axisapi"]