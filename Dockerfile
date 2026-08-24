FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/server

FROM alpine:3.20
WORKDIR /app
RUN adduser -D -H app
COPY --from=build /out/server /app/server
COPY web/static /app/web/static
ENV PORT=8080
EXPOSE 8080
USER app
CMD ["/app/server"]
