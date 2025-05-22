package store

import (
	"context"
	"math/rand"
	"quote-api/internal/model"
	"sync"
)

// В целом тут есть разные направления где оптимизировать
// Можно сделать только quotes map[int]model.Quote, но тогда у нас будет рандомный элемент за линию
// Можно обойти если использовать что обход мапы внутри устроен через рандом, но это плохая практика т.к он на самом деле iteration order is unspecified
// Можно сделать как тут и за счет дополнительной памяти на хранение индексов сделать поиск рандомной и удаление цитаты за O(1) по времени
// Выбирать лучше из контекста прода, тут просто пример для среднего случая

type inMemoryStore struct {
	sync.RWMutex
	quotes map[int]model.Quote
	ids    []int
	idIdx  map[int]int
	nextID int
}

func NewInMemoryStore() QuoteStore {
	return &inMemoryStore{
		quotes: make(map[int]model.Quote),
		ids:    []int{},
		idIdx:  make(map[int]int),
		nextID: 1,
	}
}

func (s *inMemoryStore) Add(ctx context.Context, q model.Quote) model.Quote {
	if q.Author == "" || q.Text == "" {
		return model.Quote{}
	}

	s.Lock()
	defer s.Unlock()
	q.ID = s.nextID
	s.quotes[q.ID] = q
	s.ids = append(s.ids, q.ID)
	s.idIdx[q.ID] = len(s.ids) - 1
	s.nextID++
	return q
}

func (s *inMemoryStore) GetAll(ctx context.Context) []model.Quote {
	s.RLock()
	defer s.RUnlock()
	result := make([]model.Quote, 0, len(s.quotes))
	for _, q := range s.quotes {
		result = append(result, q)
	}
	return result
}

func (s *inMemoryStore) GetByAuthor(ctx context.Context, author string) []model.Quote {
	s.RLock()
	defer s.RUnlock()
	var result []model.Quote
	for _, q := range s.quotes {
		if q.Author == author {
			result = append(result, q)
		}
	}
	return result
}

func (s *inMemoryStore) GetRandom(ctx context.Context) (model.Quote, error) {
	s.RLock()
	defer s.RUnlock()
	if len(s.ids) == 0 {
		return model.Quote{}, model.ErrQuoteNotFound
	}
	randomID := s.ids[rand.Intn(len(s.ids))]
	return s.quotes[randomID], nil
}

func (s *inMemoryStore) Delete(ctx context.Context, id int) bool {
	s.Lock()
	defer s.Unlock()
	idx, exists := s.idIdx[id]
	if !exists {
		return false
	}

	lastIdx := len(s.ids) - 1
	lastID := s.ids[lastIdx]

	s.ids[idx] = lastID
	s.idIdx[lastID] = idx

	s.ids = s.ids[:lastIdx]
	delete(s.idIdx, id)
	delete(s.quotes, id)
	return true
}
