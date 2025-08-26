package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"l0-wb-tech/internal/cache"
	"l0-wb-tech/internal/database"

	"gorm.io/gorm"
)

type Handler struct {
	cache  *cache.Cache
	db     *database.Storage
	logger *slog.Logger
}

func New(c *cache.Cache, db *database.Storage, log *slog.Logger) *Handler {
	return &Handler{
		cache:  c,
		db:     db,
		logger: log,
	}
}

func (h *Handler) GetOrderByUID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Warn("Неподдерживаемый метод HTTP", "method", r.Method)
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	orderUID := strings.TrimPrefix(r.URL.Path, "/order/")
	if orderUID == "" {
		h.logger.Warn("Пустой order_uid в запросе")
		http.Error(w, "Order UID обязателен", http.StatusBadRequest)
		return
	}

	const maxLen int = 100

	if len(orderUID) > maxLen {
		h.logger.Warn("Слишком длинный order_uid в запросе", "length", len(orderUID))
		http.Error(w, "Некорректный Order UID", http.StatusBadRequest)
		return
	}

	order, found := h.cache.Get(orderUID)

	if found {
		h.logger.Info("Заказ найден в кэше", "order_uid", orderUID)
		h.respondWithJSON(w, http.StatusOK, order)
		return
	}

	h.logger.Info("Промах кэша, поиск в БД", "order_uid", orderUID)

	orderFromDB, err := h.db.GetOrderByUID(r.Context(), orderUID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.Warn("Заказ не найден", "order_uid", orderUID)
			http.Error(w, "Заказ не найден", http.StatusNotFound)
		} else {
			h.logger.Error("Ошибка при поиске заказа в БД",
				"error", err,
				"order_uid", orderUID)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}

		return
	}

	h.logger.Info("Заказ найден в БД, обновляем кэш", "order_uid", orderUID)
	h.cache.Set(orderFromDB)
	h.respondWithJSON(w, http.StatusOK, orderFromDB)
}

func (h *Handler) respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		h.logger.Error("Ошибка при кодировании JSON ответа",
			"error", err,
			"status", status)

		if status == http.StatusOK {
			http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		}
	}
}
