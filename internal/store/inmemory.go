package store

import (
	"fmt"
	"math/rand"
	"quote-api/internal/model"
	"sync"
)

type inMemoryStore struct {
	sync.Mutex
	quotes map[int]model.QuoteNote
	nextID int
}

func NewInMemoryStore() QuoteStore {
	return &inMemoryStore{
		quotes: make(map[int]model.QuoteNote),
		nextID: 1,
	}
}

func (s *inMemoryStore) Add(q model.QuoteNote) model.QuoteNote {
	s.Lock()
	defer s.Unlock()
	q.ID = s.nextID
	s.quotes[q.ID] = q
	s.nextID++
	return q
}

func (s *inMemoryStore) GetAll() []model.QuoteNote {
	s.Lock()
	defer s.Unlock()
	result := make([]model.QuoteNote, 0, len(s.quotes))
	for _, q := range s.quotes {
		result = append(result, q)
	}
	return result
}

func (s *inMemoryStore) GetByAuthor(author string) []model.QuoteNote {
	s.Lock()
	defer s.Unlock()
	var result []model.QuoteNote
	for _, q := range s.quotes {
		if q.Author == author {
			result = append(result, q)
		}
	}
	return result
}

func (s *inMemoryStore) GetRandom() (model.QuoteNote, error) {
	s.Lock()
	defer s.Unlock()
	if len(s.quotes) == 0 {
		return model.QuoteNote{}, fmt.Errorf("no quotes available")
	}
	index := rand.Intn(len(s.quotes))
	for _, q := range s.quotes {
		if index == 0 {
			return q, nil
		}
		index--
	}
	return model.QuoteNote{}, fmt.Errorf("unexpected error")
}

func (s *inMemoryStore) Delete(id int) bool {
	s.Lock()
	defer s.Unlock()
	if _, exists := s.quotes[id]; exists {
		delete(s.quotes, id)
		return true
	}
	return false
}
