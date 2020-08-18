package repository

import (
	"context"
	"sync"

	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/dmitryt/otus-golang-hw/hw12_13_14_15_calendar/internal/utils"
	"github.com/imdario/mergo"
)

type MemoryRepo struct {
	storage map[int64]Event
	mx      sync.RWMutex
}

func (r *MemoryRepo) Connect(ctx context.Context, c *config.Config) error {
	return nil
}

func (r *MemoryRepo) Close() error {
	r.storage = nil
	return nil
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{storage: make(map[int64]Event)}
}

// Probably, need to display events only for particular user.
func (r *MemoryRepo) GetEvents() ([]Event, error) {
	r.mx.RLock()
	defer r.mx.RLock()
	result := []Event{}
	for _, v := range r.storage {
		result = append(result, v)
	}
	return result, nil
}

func (r *MemoryRepo) GetEvent(id int64) (ev Event, err error) {
	r.mx.RLock()
	defer r.mx.RLock()
	tmp, ok := r.storage[id]
	if !ok {
		return Event{}, ErrEventNotFound
	}
	ev = tmp
	return
}

func (r *MemoryRepo) CreateEvent(data Event) (Event, error) {
	id, err := utils.GenerateUID()
	if err != nil {
		return Event{}, ErrEventCreate
	}
	data.ID = id
	r.mx.Lock()
	defer r.mx.Unlock()
	r.storage[id] = data
	return data, nil
}

func (r *MemoryRepo) UpdateEvent(data Event) (event Event, err error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	event, ok := r.storage[data.ID]
	if !ok {
		return Event{}, ErrEventNotFound
	}
	err = mergo.Merge(&event, data, mergo.WithOverride)
	r.storage[data.ID] = event
	return
}

func (r *MemoryRepo) DeleteEvent(id int64) (err error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	_, ok := r.storage[id]
	if !ok {
		return ErrEventNotFound
	}
	delete(r.storage, id)
	return
}
