package handlers

import (
	"encoding/json"
	"l0-wb-tech/internal/cache"
	"log/slog"
	"net/http"
	"strings"
)

type Handler struct {
	cache  *cache.Cache
	logger *slog.Logger
}

func New(c *cache.Cache, log *slog.Logger) *Handler {
	return &Handler{
		cache: c,
		logger: log,
	}
}

func (h *Handler) GetOrderByUID(w http.ResponseWriter, r *http.Request) {
	orderUID := strings.TrimPrefix(r.URL.Path, "/order/")

	if orderUID == "" {
		http.Error(w, "Order UID is required", http.StatusBadRequest)
	}

	order, is_found := h.cache.Get(orderUID)

	if !is_found {
		h.logger.Warn("Заказ не найден в кэше", slog.String("order_uid", orderUID))
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(order); err != nil {
		h.logger.Error("Ошибка при кодировании JSON-ответа", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Заказ успешно отдан из кэша", slog.String("order_uid", orderUID))
}