FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /out/agent ./main.go

FROM alpine:3.20
RUN adduser -D -g '' appuser
USER appuser
WORKDIR /app
COPY --from=build /out/agent /app/agent
ENTRYPOINT ["/app/agent"]
