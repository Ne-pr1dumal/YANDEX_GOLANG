FROM golang:1.23.4

WORKDIR /app
COPY go.mod ./
COPY internal/ ./internal/
COPY cmd/agent_start/main.go ./cmd/agent_start/

RUN go mod download
RUN go build -o /agent ./cmd/agent_start/main.go

CMD ["/agent"]