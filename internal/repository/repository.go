package repository

import (
	"encoding/json"
	"errors"
	"sync"
)

var ErrInexistent = errors.New("key does not exist")

type KVRepository struct {
	data  sync.Map
	mutex sync.Mutex // Add a mutex for synchronization
}

func NewKVRepository() *KVRepository {
	return &KVRepository{}
}

func (r *KVRepository) Put(key string, value interface{}) error {
	r.mutex.Lock()         // Lock before writing
	defer r.mutex.Unlock() // Ensure unlocking after the operation

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	r.data.Store(key, data)
	return nil
}

func (r *KVRepository) Get(key string) (interface{}, error) {
	r.mutex.Lock()         // Lock for reading
	defer r.mutex.Unlock() // Ensure unlocking after the operation

	value, ok := r.data.Load(key)
	if !ok {
		return nil, ErrInexistent
	}

	var result interface{}
	if err := json.Unmarshal(value.([]byte), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *KVRepository) Delete(key string) error {
	r.mutex.Lock()         // Lock before deleting
	defer r.mutex.Unlock() // Ensure unlocking after the operation

	if _, ok := r.data.Load(key); !ok {
		return ErrInexistent
	}
	r.data.Delete(key)
	return nil
}

func (r *KVRepository) ListKeys() []string {
	r.mutex.Lock()         // Lock while reading all keys
	defer r.mutex.Unlock() // Ensure unlocking after the operation

	keys := []string{}
	r.data.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys
}
