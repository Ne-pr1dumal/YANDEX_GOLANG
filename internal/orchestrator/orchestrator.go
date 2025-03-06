package orchestrator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

type ASTNode struct {
	IsLeaf        bool
	Value         float64
	Operator      string
	Left, Right   *ASTNode
	TaskScheduled bool
}

type parser struct {
	input string
	pos   int
}

func ParseAST(expression string) (*ASTNode, error) {
	expr := strings.ReplaceAll(expression, " ", "")
	if expr == "" {
		return nil, fmt.Errorf("empty expression")
	}

	p := &parser{input: expr, pos: 0}
	node, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if p.pos < len(p.input) {
		return nil, fmt.Errorf("unexpected token at position %d", p.pos)
	}
	return node, nil
}

func (p *parser) peek() rune {
	if p.pos < len(p.input) {
		return rune(p.input[p.pos])
	}
	return 0
}

func (p *parser) get() rune {
	ch := p.peek()
	p.pos++
	return ch
}

func (p *parser) parseExpression() (*ASTNode, error) {
	node, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for {
		switch p.peek() {
		case '+', '-':
			op := string(p.get())
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			node = &ASTNode{Operator: op, Left: node, Right: right}
		default:
			return node, nil
		}
	}
}

func (p *parser) parseTerm() (*ASTNode, error) {
	node, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for {
		switch p.peek() {
		case '*', '/':
			op := string(p.get())
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			node = &ASTNode{Operator: op, Left: node, Right: right}
		default:
			return node, nil
		}
	}
}

func (p *parser) parseFactor() (*ASTNode, error) {
	if p.peek() == '(' {
		p.get()
		node, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.get() != ')' {
			return nil, fmt.Errorf("missing closing parenthesis")
		}
		return node, nil
	}

	start := p.pos
	if p.peek() == '+' || p.peek() == '-' {
		p.get()
	}

	for unicode.IsDigit(p.peek()) || p.peek() == '.' {
		p.get()
	}

	numStr := p.input[start:p.pos]
	if numStr == "" {
		return nil, fmt.Errorf("invalid number at position %d", start)
	}

	value, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number format: %s", numStr)
	}

	return &ASTNode{IsLeaf: true, Value: value}, nil
}

type Config struct {
	Addr                string
	TimeAddition        int
	TimeSubtraction     int
	TimeMultiplications int
	TimeDivisions       int
}

func Configuration() *Config {
	cfg := &Config{
		Addr:                "8080",
		TimeAddition:        100,
		TimeSubtraction:     100,
		TimeMultiplications: 100,
		TimeDivisions:       100,
	}

	if port := os.Getenv("PORT"); port != "" {
		cfg.Addr = port
	}

	if ta, err := strconv.Atoi(os.Getenv("TIME_ADDITION_MS")); err == nil && ta > 0 {
		cfg.TimeAddition = ta
	}

	if ts, err := strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS")); err == nil && ts > 0 {
		cfg.TimeSubtraction = ts
	}

	if tm, err := strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS")); err == nil && tm > 0 {
		cfg.TimeMultiplications = tm
	}

	if td, err := strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS")); err == nil && td > 0 {
		cfg.TimeDivisions = td
	}

	return cfg
}

type Orchestrator struct {
	Config      *Config
	exprStore   map[string]*Expression
	taskStore   map[string]*Task
	taskQueue   []*Task
	mu          sync.Mutex
	exprCounter int64
	taskCounter int64
}

type Expression struct {
	ID     string   `json:"id"`
	Expr   string   `json:"expression"`
	Status string   `json:"status"`
	Result *float64 `json:"result,omitempty"`
	AST    *ASTNode `json:"-"`
}

type Task struct {
	ID            string   `json:"id"`
	ExprID        string   `json:"-"`
	Arg1          float64  `json:"arg1"`
	Arg2          float64  `json:"arg2"`
	Operation     string   `json:"operation"`
	OperationTime int      `json:"operation_time"`
	Node          *ASTNode `json:"-"`
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		Config:    Configuration(),
		exprStore: make(map[string]*Expression),
		taskStore: make(map[string]*Task),
		taskQueue: make([]*Task, 0),
	}
}

func (o *Orchestrator) calculateHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req struct{ Expression string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Expression == "" {
		http.Error(w, `{"error":"Invalid Request"}`, http.StatusBadRequest)
		return
	}

	ast, err := ParseAST(req.Expression)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusUnprocessableEntity)
		return
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	o.exprCounter++
	exprID := fmt.Sprintf("%d", o.exprCounter)
	expr := &Expression{
		ID:     exprID,
		Expr:   req.Expression,
		Status: "pending",
		AST:    ast,
	}

	o.exprStore[exprID] = expr
	o.Tasks(expr)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": exprID})
}

func (o *Orchestrator) Tasks(expr *Expression) {
	var traverse func(*ASTNode)
	traverse = func(node *ASTNode) {
		if node == nil || node.IsLeaf {
			return
		}

		traverse(node.Left)
		traverse(node.Right)

		if node.Left.IsLeaf && node.Right.IsLeaf && !node.TaskScheduled {
			o.taskCounter++
			taskID := fmt.Sprintf("%d", o.taskCounter)

			opTime := 100
			switch node.Operator {
			case "+":
				opTime = o.Config.TimeAddition
			case "-":
				opTime = o.Config.TimeSubtraction
			case "*":
				opTime = o.Config.TimeMultiplications
			case "/":
				opTime = o.Config.TimeDivisions
			}

			task := &Task{
				ID:            taskID,
				ExprID:        expr.ID,
				Arg1:          node.Left.Value,
				Arg2:          node.Right.Value,
				Operation:     node.Operator,
				OperationTime: opTime,
				Node:          node,
			}

			node.TaskScheduled = true
			o.taskStore[taskID] = task
			o.taskQueue = append(o.taskQueue, task)
		}
	}
	traverse(expr.AST)
}

func (o *Orchestrator) RunServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", o.calculateHandler)
	mux.HandleFunc("/api/v1/expressions", o.expressionsHandler)
	mux.HandleFunc("/api/v1/expressions/", o.expressionIDHandler)
	mux.HandleFunc("/internal/task", o.taskHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"Not Found"}`, http.StatusNotFound)
	})

	go o.monitorTasks()
	return http.ListenAndServe(":"+o.Config.Addr, mux)
}

func (o *Orchestrator) monitorTasks() {
	for range time.Tick(2 * time.Second) {
		o.mu.Lock()
		if len(o.taskQueue) > 0 {
			log.Printf("Pending tasks: %d", len(o.taskQueue))
		}
		o.mu.Unlock()
	}
}

func (o *Orchestrator) taskHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	switch r.Method {
	case http.MethodGet:
		o.getTaskHandler(w, r)
	case http.MethodPost:
		o.postTaskHandler(w, r)
	default:
		http.Error(w, `{"error":"Method Not Allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (o *Orchestrator) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	if len(o.taskQueue) == 0 {
		http.Error(w, `{"error":"No Tasks"}`, http.StatusNotFound)
		return
	}

	task := o.taskQueue[0]
	o.taskQueue = o.taskQueue[1:]

	if expr, exists := o.exprStore[task.ExprID]; exists {
		expr.Status = "in_progress"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"task": task})
}

func (o *Orchestrator) postTaskHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	var req struct {
		ID     string  `json:"id"`
		Result float64 `json:"result"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ID == "" {
		http.Error(w, `{"error":"Invalid Request"}`, http.StatusBadRequest)
		return
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	task, exists := o.taskStore[req.ID]
	if !exists {
		http.Error(w, `{"error":"Task Not Found"}`, http.StatusNotFound)
		return
	}

	task.Node.IsLeaf = true
	task.Node.Value = req.Result
	delete(o.taskStore, req.ID)

	if expr, exists := o.exprStore[task.ExprID]; exists {
		o.Tasks(expr)
		if expr.AST.IsLeaf {
			expr.Status = "completed"
			expr.Result = &expr.AST.Value
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
}

// Обработчики выражений
func (o *Orchestrator) expressionsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	exprs := make([]*Expression, 0, len(o.exprStore))
	for _, expr := range o.exprStore {
		if expr.AST != nil && expr.AST.IsLeaf {
			expr.Status = "completed"
			expr.Result = &expr.AST.Value
		}
		exprs = append(exprs, expr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": exprs})
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func (o *Orchestrator) expressionIDHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/api/v1/expressions/"):]
	o.mu.Lock()
	expr, exists := o.exprStore[id]
	o.mu.Unlock()

	if !exists {
		http.Error(w, `{"error":"Not Found"}`, http.StatusNotFound)
		return
	}

	if expr.AST != nil && expr.AST.IsLeaf {
		expr.Status = "completed"
		expr.Result = &expr.AST.Value
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"expression": expr})
}
