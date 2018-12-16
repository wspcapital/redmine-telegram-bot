package converse_state

import (
	"errors"
	"log"
	"time"
)

// StateStore will find or create new state by ID.
func FindOrCreateState(stateStore StateStore, ID int64) (*State, error) {

	state := stateStore.GetStateById(ID)
	if state == nil {
		state = stateStore.CreateState(ID)
	}

	if state == nil {
		return nil, errors.New("state can not be found and created")
	}

	return state, nil
}

// The blank for the case of adding default values.
func New(stateStore StateStore) *State {
	state := new(State)
	state.stateStore = stateStore
	return state
}

// State of conversation with user.
type State struct {
	ID              int64
	UserName        string
	IsLogged        bool
	CurrentQuestion int64
	isCreatedNow    bool
	UpdatedAt       time.Time
	CreatedAt       time.Time

	stateStore StateStore
}

// State store is an intermediary between the application and a data store.
type StateStore interface {
	UpdateState(state *State)
	GetStateById(id int64) *State
	CreateState(id int64) *State
}

// Check is this is first time chat created (getter).
func (s *State) IsJustCreated() bool {
	return s.isCreatedNow
}

// Setter of field
func (s *State) SetJustCreated(status bool) *State {
	s.isCreatedNow = status

	return s
}

// Save current state to data store
func (s *State) UpdateState() {
	if s.stateStore != nil {
		s.stateStore.UpdateState(s)
		return
	}

	log.Fatal("Can't update state because state store not defines")
}
