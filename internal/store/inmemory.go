package store

import (
	"errors"
	"math/rand"
	"quote-api/internal/model"
	"sync"
)

// В целом тут есть разные направления где оптимизировать
// Можно сделать только quotes map[int]model.Text, но тогда у нас будет рандомный элемент за линию
// Можно обойти если использовать что обход мапы сейчас внутри устроен через рандом, но это плохая практика т.к он на самом деле iteration order is unspecified
// Можно сделать как тут и за счет дополнительной памяти на хранение индексов сделать рандомизацию и удаление за O(1) по времени
// Выбирать лучше из контекста прода, тут просто пример для среднего случая

type inMemoryStore struct {
	sync.RWMutex
	quotes map[int]model.Quote // основное хранение цитат
	ids    []int               // массив ID, чтобы можно было за O(1) взять случайный элемент
	idIdx  map[int]int         // ID -> индекс в ids
	nextID int
}

func NewInMemoryStore() QuoteStore {
	return &inMemoryStore{
		quotes: make(map[int]model.Quote),
		ids:    make([]int, 0),
		idIdx:  make(map[int]int),
		nextID: 1,
	}
}

func (s *inMemoryStore) Add(q model.Quote) model.Quote {
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

func (s *inMemoryStore) GetAll() []model.Quote {
	s.RLock()
	defer s.RUnlock()
	result := make([]model.Quote, 0, len(s.quotes))
	for _, q := range s.quotes {
		result = append(result, q)
	}
	return result
}

func (s *inMemoryStore) GetByAuthor(author string) []model.Quote {
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

func (s *inMemoryStore) GetRandom() (model.Quote, error) {
	s.RLock()
	defer s.RUnlock()
	if len(s.ids) == 0 {
		return model.Quote{}, errors.New("no quotes available")
	}
	randomID := s.ids[rand.Intn(len(s.ids))]
	return s.quotes[randomID], nil
}

func (s *inMemoryStore) Delete(id int) bool {
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
