# Build
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /out/app ./main.go

# Run
FROM alpine:3.20
RUN adduser -D -g '' appuser
USER appuser
WORKDIR /app
COPY --from=build /out/app /app/app
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/app"]
