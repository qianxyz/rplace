FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/rplace .

FROM alpine:3
WORKDIR /app
COPY --from=build /out/rplace ./rplace
COPY static static
EXPOSE 8080
CMD ["./rplace"]
