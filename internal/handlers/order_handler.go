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
	orderUID := strings.TrimPrefix(r.URL.Path, "/order/")

	if orderUID == "" {
		http.Error(w, "Order UID is required", http.StatusBadRequest)
		return
	}

	order, found := h.cache.Get(orderUID)

	if found {
		h.logger.Info("Заказ успешно отдан из кэша", slog.String("order_uid", orderUID))
		h.respondWithJSON(w, http.StatusOK, order)
		return
	}

	h.logger.Warn("Промах кэша, ищем заказ в БД", slog.String("order_uid", orderUID))
	orderFromDB, err := h.db.GetOrderByUID(r.Context(), orderUID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.Warn("Заказ не найден ни в кэше, ни в БД", slog.String("order_uid", orderUID))
			http.NotFound(w, r)
		} else {
			h.logger.Error("Ошибка при поиске заказа в БД", slog.Any("error", err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		return
	}

	h.logger.Info("Заказ найден в БД, кэш обновлен", slog.String("order_uid", orderUID))
	h.cache.Set(orderFromDB)

	h.respondWithJSON(w, http.StatusOK, orderFromDB)
}

func (h *Handler) respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		h.logger.Error("Ошибка при кодировании JSON-ответа", slog.Any("error", err))
	}
}
