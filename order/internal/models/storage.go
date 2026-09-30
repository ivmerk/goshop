package models

import (
	"sync"
)

type OrderStorage struct {
	orders map[string]*Order // Карта для хранения заказов, где ключ - это уникальный идентификатор заказа (ID), а значение - указатель на структуру Order.
	mu     sync.RWMutex      // Мьютекс для обеспечения безопасного доступа к карте orders в многопоточной среде.
}

// NewOrderStorage создает и возвращает новый экземпляр OrderStorage с инициализированной картой orders.
func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*Order),
	}
}

func (s *OrderStorage) GetOrder(uuid string) (*Order, bool) {
	s.mu.RLock()         // Блокировка для чтения, чтобы предотвратить запись в карту orders во время чтения.
	defer s.mu.RUnlock() // Отложенная разблокировка после завершения функции.

	order, exists := s.orders[uuid] // Получение заказа по его уникальным идентификатору (UUID) из карты orders.
	return order, exists            // Возвращаем заказ и флаг существования заказа в карте.
}

func (s *OrderStorage) AddOrder(order *Order) {
	s.mu.Lock()         // Блокировка для записи, чтобы предотвратить одновременное чтение или запись в карту orders.
	defer s.mu.Unlock() // Отложенная разблокировка после завершения функции.

	s.orders[order.UUID] = order // Добавление нового заказа в карту orders по его уникальным идентификатору (UUID).
}

func (s *OrderStorage) UpdateOrder(order *Order) {
	s.mu.Lock()         // Блокировка для записи, чтобы предотвратить одновременное чтение или запись в карту orders.
	defer s.mu.Unlock() // Отложенная разблокировка после завершения функции.

	s.orders[order.UUID] = order // Обновление или добавление заказа в карту orders по его уникальным идентификатору (UUID).
}
