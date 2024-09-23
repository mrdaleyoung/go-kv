package services

import "go-kv/internal/repository"

// KVServiceInterface defines the methods for the KV service
type KVServiceInterface interface {
	Put(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	ListKeys() []string
}

type KVService struct {
	repo *repository.KVRepository
}

func NewKVService(repo *repository.KVRepository) *KVService {
	return &KVService{repo: repo}
}

func (s *KVService) Put(key string, value interface{}) error {
	return s.repo.Put(key, value)
}

func (s *KVService) Get(key string) (interface{}, error) {
	return s.repo.Get(key)
}

func (s *KVService) Delete(key string) error {
	return s.repo.Delete(key)
}

func (s *KVService) ListKeys() []string {
	return s.repo.ListKeys()
}
