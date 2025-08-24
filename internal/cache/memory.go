package cache

import (
	"sync"

	"l0-wb-tech/internal/models"
)

type Cache struct {
	mu      sync.RWMutex
	storage map[string]models.Order
}

func New() *Cache {
	return &Cache{storage: make(map[string]models.Order)}
}

func (c *Cache) Get(orderUID string) (models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	order, is_found := c.storage[orderUID]

	return order, is_found
}

func (c *Cache) Load(orders []models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, order := range orders {
		c.storage[order.OrderUID] = order
	}
}
