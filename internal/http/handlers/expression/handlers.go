package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/1minepowminx/distributed_calculator/internal/storage"
	"github.com/1minepowminx/distributed_calculator/internal/utils/orchestrator/jwts"
)

type Request struct {
	Expression string `json:"expression"`
}

type ResponseData struct {
	ID         int64  `json:"id"`
	Expression string `json:"expression"`
	Answer     string `json:"answer"`
	Date       string `json:"date"`
	Status     string `json:"status"`
}

type ExpressionInteractor interface {
	InsertExpression(ctx context.Context, expr *storage.Expression) (int64, error)
	SelectExpressionsByID(ctx context.Context, userID int64) ([]storage.Expression, error)
	DeleteExpression(ctx context.Context, id int64) error
}

type contextKey string

const UserIDKey contextKey = "userid"

// Обработчик для сохранения выражения в БД
func CreateExpressionHandler(ctx context.Context, expressionSaver ExpressionInteractor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		date := time.Now()

		jsonDec := json.NewDecoder(r.Body)

		var req Request
		if err := jsonDec.Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		cookie, err := r.Cookie("auth_token")
		if err != nil {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			log.Printf("No cookie found")
			return
		}

		tokenString := cookie.Value

		tokenValue, err := jwts.VerifyJWTToken(tokenString)
		if err != nil {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			log.Printf("Error: %v", err)
			return
		}
		userID, err := strconv.ParseInt(tokenValue, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error: %v", err)
			return
		}

		var expressionStruct = storage.Expression{
			UserID:     userID,
			Expression: req.Expression,
			Answer:     "null",
			Date:       date.Format("2006/01/02 15:04:05"),
			Status:     "stored",
		}

		id, err := expressionSaver.InsertExpression(ctx, &expressionStruct)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		log.Printf("Successful CreateExpressionHandler operation; id = %d", id)
	}
}

// Обработчик для получения списка всех выражений из БД
func GetExpressionsHandler(ctx context.Context, expressionSaver ExpressionInteractor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		cookie, err := r.Cookie("auth_token")
		if err != nil {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			log.Printf("No cookie found")
			return
		}

		tokenString := cookie.Value

		tokenValue, err := jwts.VerifyJWTToken(tokenString)
		if err != nil {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			log.Printf("Error: %v", err)
			return
		}
		userID, err := strconv.ParseInt(tokenValue, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error: %v", err)
			return
		}

		allExpressions, err := expressionSaver.SelectExpressionsByID(ctx, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var respData []ResponseData

		for _, expr := range allExpressions {
			resp := ResponseData{
				ID:         expr.ID,
				Expression: expr.Expression,
				Answer:     expr.Answer,
				Date:       expr.Date,
				Status:     expr.Status,
			}

			respData = append(respData, resp)
		}

		json.NewEncoder(w).Encode(respData)
		log.Print("Successful GetExpressionsHandler operation")
	}
}

// Обработчик удаления выражения из БД
func DeleteExpressionHandler(ctx context.Context, expressionSaver ExpressionInteractor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = expressionSaver.DeleteExpression(ctx, int64(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		log.Print("Successful DeleteExpressionHandler operation")
	}
}
