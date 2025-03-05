package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Easter Egg
//             _                   
//  __ _  ___ | | __ _ _ __   __ _ 
// / _` |/ _ \| |/ _` | '_ \ / _` |
//| (_| | (_) | | (_| | | | | (_| |
// \__, |\___/|_|\__,_|_| |_|\__, |
// |___/                     |___/
var (
	ErrDivisionByZero  = errors.New("division by zero")
	ErrInvalidOperator = errors.New("invalid operator")
)

type Agent struct {
	ComputingPower  int
	OrchestratorURL string
}

func NewAgent() *Agent {
	cp, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil || cp < 1 {
		cp = 1
	}

	orchestratorURL := os.Getenv("ORCHESTRATOR_URL")
	if orchestratorURL == "" {
		orchestratorURL = "http://localhost:8080"
	}

	return &Agent{
		ComputingPower:  cp,
		OrchestratorURL: orchestratorURL,
	}
}

func (a *Agent) Start() {
	for i := 0; i < a.ComputingPower; i++ {
		log.Printf("Starting worker %d", i)
		go a.worker(i)
	}
	select {}
}

func (a *Agent) worker(id int) {
	for {
		task, ok := a.fetchTask(id)
		if !ok {
			time.Sleep(2 * time.Second)
			continue
		}

		result, err := a.processTask(id, task)
		if err != nil {
			log.Printf("Worker %d: error processing task %s: %v", id, task.ID, err)
			continue
		}

		a.sendResult(id, task.ID, result)
	}
}

func (a *Agent) fetchTask(id int) (*struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}, bool) {
	resp, err := http.Get(a.OrchestratorURL + "/internal/task")
	if err != nil {
		log.Printf("Worker %d: error fetching task: %v", id, err)
		return nil, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, false
		}
		log.Printf("Worker %d: unexpected status code: %d", id, resp.StatusCode)
		return nil, false
	}

	var taskResp struct {
		Task struct {
			ID            string  `json:"id"`
			Arg1          float64 `json:"arg1"`
			Arg2          float64 `json:"arg2"`
			Operation     string  `json:"operation"`
			OperationTime int     `json:"operation_time"`
		} `json:"task"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		log.Printf("Worker %d: error decoding task: %v", id, err)
		return nil, false
	}

	return &taskResp.Task, true
}

func (a *Agent) processTask(id int, task *struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}) (float64, error) {
	log.Printf("Worker %d: processing task %s: %f %s %f (%d ms)",
		id, task.ID, task.Arg1, task.Operation, task.Arg2, task.OperationTime)

	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

	result, err := calculate(task.Operation, task.Arg1, task.Arg2)
	if err != nil {
		return 0, fmt.Errorf("calculation error: %w", err)
	}

	return result, nil
}

func (a *Agent) sendResult(id int, taskID string, result float64) {
	resultPayload := map[string]interface{}{
		"id":     taskID,
		"result": result,
	}

	payloadBytes, _ := json.Marshal(resultPayload)
	resp, err := http.Post(a.OrchestratorURL+"/internal/task", "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		log.Printf("Worker %d: error sending result for task %s: %v", id, taskID, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Worker %d: error response for task %s: %s", id, taskID, string(body))
		return
	}

	log.Printf("Worker %d: successfully completed task %s with result %f", id, taskID, result)
}

func calculate(operation string, a, b float64) (float64, error) {
	switch operation {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		return a / b, nil
	default:
		return 0, ErrInvalidOperator
	}
}

func CalculateExpression(expression string) (float64, error) {
	return 0, fmt.Errorf("not implemented")
}
