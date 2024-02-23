package main

import (
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Knetic/govaluate"
)

var (
	idCounter        = 1 // Идентификатор выражения упрощен для демонстрации
	expressionsStore = NewExpressionsStore()
)

type Expression struct { // БД слушает на :5432
	ID           int       `json:"id"`
	Content      string    `json:"content"`
	Status       string    `json:"status"`
	Result       string    `json:"result"`
	CreatedAt    time.Time `json:"created_at"`
	CalculatedAt time.Time `json:"calculated_at"`
}

type ComputationalCapability struct {
	Resource  string `json:"resource"`
	Operation string `json:"operation"`
}

type ExpressionsStore struct {
	Expressions              map[int]*Expression `json:"expressions"`
	OperationsTime           map[string]int      `json:"operations"`
	MachineNums              int                 `json:"machine_nums"`
	ComputationalCapabilitys map[string]string   `json:"computational_capabilities"`
	CurrentWorkers           int                 `json:"current_workers"`
	mu                       sync.Mutex          `json:"-"`
	wg                       sync.WaitGroup      `json:"-"`
}

func NewExpressionsStore() *ExpressionsStore {
	// Инициализируем хранилище выражений дефолтными значениями
	return &ExpressionsStore{
		Expressions:              make(map[int]*Expression),
		OperationsTime:           map[string]int{"add": 60, "sub": 60, "mul": 120, "div": 120},
		MachineNums:              4,
		ComputationalCapabilitys: make(map[string]string),
	}
}

func (e *ExpressionsStore) AddExpression(expression *Expression) bool {
	// Проверяем наличие выражения в хранилище по ID
	_, ok := e.Expressions[expression.ID]
	if !ok {
		// Если выражение отсутствует - добавляем
		e.Expressions[expression.ID] = expression
		return true
	}

	return false
}

func (e *ExpressionsStore) SetCompCapability(cpuName string, operations []string) {
	e.ComputationalCapabilitys[cpuName] = strings.Join(operations, ", ")
}

func SubmitExpression(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметры из запроса
	queryParams := r.URL.Query()

	if queryParams.Get("content") == "" {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	content := queryParams.Get("content")

	// Проверяем выражение на наличие недопустимых символов
	validExpression := regexp.MustCompile(`^[0-9+\-*/\s()]+$`).MatchString(content)
	if !validExpression {
		http.Error(w, "Expression parsing error", http.StatusBadRequest)
		return
	}

	// Создаем новое выражение
	expression := &Expression{
		ID:        idCounter,
		Content:   content,
		Status:    "The expression will be calculated soon",
		Result:    "?",
		CreatedAt: time.Now(),
	}

	// Увеличиваем счетчик идентификатора выражения
	idCounter++

	// Добавляем выражение в хранилище
	expressionsStore.AddExpression(expression)
	w.Write([]byte("The expression has been successfully accepted. Current ID: " + strconv.Itoa(expression.ID)))

	// Выполнение дополнительных операций (если таковые имеются)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Отправляем GET-запрос на /get-task для запуска демона
				response, err := http.Get("http://localhost:8080/get-task")
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer response.Body.Close()
			}
		}
	}()
}

func GetExpressionsList(w http.ResponseWriter, r *http.Request) {
	// Отправляем GET-запрос на /get-result для получения списка выражений с результатами обработки данных
	response, err := http.Get("http://localhost:8080/get-result")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем список выражений с результатами обработки данных
	_, err = w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetExpressionsId(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметры из запроса
	queryParams := r.URL.Query()

	if queryParams.Get("id") == "" {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	expressionId, err := strconv.Atoi(queryParams.Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем запрос к хранилищу и получаем данные о выражении по его ID
	value, ok := expressionsStore.Expressions[expressionId]
	if !ok {
		http.Error(w, "Expression not found: ID"+strconv.Itoa(expressionId), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(value.Content + " - " + value.Status + " - " + value.CreatedAt.Format("2006-01-02 15:04:05") + " - " + value.CalculatedAt.Format("2006-01-02 15:04:05")))
}

func OperationsList(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметры из запроса
	queryParams := r.URL.Query()

	addTime := queryParams.Get("add")
	if addTime != "" {
		addTimeConv, err := strconv.Atoi(addTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		expressionsStore.OperationsTime["add"] = addTimeConv
	}

	subTime := queryParams.Get("sub")
	if subTime != "" {
		subTimeConv, err := strconv.Atoi(subTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		expressionsStore.OperationsTime["sub"] = subTimeConv
	}

	mulTime := queryParams.Get("mul")
	if mulTime != "" {
		mulTimeConv, err := strconv.Atoi(mulTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		expressionsStore.OperationsTime["mul"] = mulTimeConv
	}

	divTime := queryParams.Get("div")
	if divTime != "" {
		divTimeConv, err := strconv.Atoi(divTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		expressionsStore.OperationsTime["div"] = divTimeConv
	}

	// Операции в виде пар: имя операции + время его выполнения
	for key, value := range expressionsStore.OperationsTime {
		w.Write([]byte("Operation execution time: " + key + " - " + strconv.Itoa(value) + " minutes\n"))
	}
}

func GetComputationalCapabilitysList(w http.ResponseWriter, r *http.Request) {
	if expressionsStore.CurrentWorkers == 0 {
		w.Write([]byte("Computing resources is free. Goruntines available: " + strconv.Itoa(expressionsStore.MachineNums) + "\n"))
		return
	}

	for key, value := range expressionsStore.ComputationalCapabilitys {
		w.Write([]byte("Computing resource name: " + key + "\n" + "Running operation: " + value + "\n"))
	}
}

func CalcMachine(expression *Expression, cpuNum int) bool {
	defer expressionsStore.wg.Done()
	expressionsStore.mu.Lock()

	operationsList := make([]string, 0, len(expression.Content))
	totalWaitTime := 0
	for _, char := range expression.Content {
		switch char {
		case '+':
			totalWaitTime += expressionsStore.OperationsTime["add"]
			operationsList = append(operationsList, "addition")
		case '-':
			totalWaitTime += expressionsStore.OperationsTime["sub"]
			operationsList = append(operationsList, "subtraction")
		case '*':
			totalWaitTime += expressionsStore.OperationsTime["mul"]
			operationsList = append(operationsList, "multiply")
		case '/':
			totalWaitTime += expressionsStore.OperationsTime["div"]
			operationsList = append(operationsList, "division")
		}
	}

	expressionsStore.SetCompCapability("Goruntine "+strconv.Itoa(cpuNum), operationsList)

	expressionsStore.mu.Unlock()

	expressionEval, err := govaluate.NewEvaluableExpression(expression.Content)
	if err != nil {
		expression.Status = "Expression parsing error"
		return false
	}

	result, err := expressionEval.Evaluate(nil)
	if err != nil {
		expression.Status = "Expression parsing error"
		return false
	}

	timer := time.NewTimer(time.Duration(totalWaitTime) * time.Minute)

	<-timer.C

	expressionsStore.mu.Lock()
	expression.Status = "Done"
	expression.CalculatedAt = time.Now()
	expression.Result = strconv.Itoa(int(result.(float64)))

	expressionsStore.SetCompCapability("Goruntine "+strconv.Itoa(cpuNum), []string{})

	expressionsStore.CurrentWorkers--

	expressionsStore.mu.Unlock()

	return true
}

func SetCalcTask(w http.ResponseWriter, r *http.Request) {

	availableWorkers := expressionsStore.MachineNums - expressionsStore.CurrentWorkers

	for _, expression := range expressionsStore.Expressions {
		expressionsStore.wg.Add(1)
		if expression.Status == "The expression will be calculated soon" && availableWorkers > 0 {
			expressionsStore.CurrentWorkers++
			go CalcMachine(expression, expressionsStore.CurrentWorkers)
			availableWorkers--
		}
	}
	go func() {
		expressionsStore.wg.Wait()
	}()
}

func GetCalcTask(w http.ResponseWriter, r *http.Request) {
	for _, expression := range expressionsStore.Expressions {
		w.Write([]byte(expression.Content + "=" + expression.Result + " - " + expression.Status + " - " + expression.CreatedAt.Format("2006-01-02 15:04:05") + " - " + expression.CalculatedAt.Format("2006-01-02 15:04:05") + "\n"))
	}
}

func CORS(next http.HandlerFunc) http.HandlerFunc {
	// Добавляем middleware для доступа к другому источнику
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET")
		// defer w.WriteHeader(http.StatusOK)

		next(w, r)
	}
}

func main() {
	// TODO: Реализовать загрузку состояния из СУБД
	// Создаем мультиплексор для управления маршрутами
	mux := http.NewServeMux()

	// Добавление вычисления арифметического выражения
	// TODO: Реализовать отсутвие возможности добавлять одинаковые выражения
	mux.HandleFunc("/submit", CORS(SubmitExpression))

	// Получение списка выражений со статусами
	mux.HandleFunc("/expressions/list", CORS(GetExpressionsList))

	// Получение значения выражения по его идентификатору
	mux.HandleFunc("/expressions/get-by-id", CORS(GetExpressionsId))

	// Получение списка доступных операций со временем их выполения
	mux.HandleFunc("/operations", CORS(OperationsList))

	// Получение задачи для выполения
	mux.HandleFunc("/get-task", CORS(SetCalcTask))

	// Приём результата обработки данных
	mux.HandleFunc("/get-result", CORS(GetCalcTask))

	// Получение списка вычислительных ресурсов
	mux.HandleFunc("/status", CORS(GetComputationalCapabilitysList))

	http.ListenAndServe(":8080", mux)
}
