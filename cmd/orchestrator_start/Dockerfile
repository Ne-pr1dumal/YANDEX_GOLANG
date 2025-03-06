FROM golang:1.23.4

WORKDIR /app
COPY go.mod ./
COPY internal/ ./internal/
COPY cmd/orchestrator_start/main.go ./cmd/orchestrator_start/

RUN go mod download
RUN go build -o /orchestrator ./cmd/orchestrator_start/main.go

EXPOSE 8080
CMD ["/orchestrator"]