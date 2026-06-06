package repository

import (
	"sync"

	"rationalgo/internal/models"
	"rationalgo/internal/store"
)

// Store holds the dashboard state in memory.
type Store struct {
	mu    sync.RWMutex
	state models.AppState
}

// NewStore returns a store seeded with demo data.
func NewStore() *Store {
	return &Store{state: store.Seed()}
}

// State returns a copy of the current dashboard state.
func (s *Store) State() models.AppState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

// SetState replaces the dashboard state.
func (s *Store) SetState(state models.AppState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = state
}

// AddDecision prepends a decision to the feed.
func (s *Store) AddDecision(d models.Decision) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.Decisions = append([]models.Decision{d}, s.state.Decisions...)
}

// UpdateDecision applies a patch to a decision by ID.
func (s *Store) UpdateDecision(id string, patch func(*models.Decision)) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.state.Decisions {
		if s.state.Decisions[i].ID == id {
			patch(&s.state.Decisions[i])
			return true
		}
	}
	return false
}

// Reset restores seed data.
func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = store.Seed()
}
