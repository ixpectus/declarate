package kv

import (
	"fmt"
	"log"

	"github.com/recoilme/pudge"
)

type KV struct {
	pudge *pudge.Db
}

func New(path string, clear bool) *KV {
	cfg := pudge.DefaultConfig
	cfg.SyncInterval = 0 // disable every second fsync
	if clear {
		pudge.DeleteFile(path)
	}
	db, err := pudge.Open(path, cfg)
	if err != nil {
		log.Panic(err)
	}

	return &KV{
		pudge: db,
	}
}

func (k *KV) Set(key string, value string) error {
	return k.pudge.Set(key, value)
}

func (k *KV) Get(key string) (string, error) {
	var value string

	err := k.pudge.Get(key, &value)
	if err != nil {
		return "", fmt.Errorf("kv get: %w", err)
	}

	return value, nil
}

func (k *KV) Reset() error {
	if err := k.pudge.DeleteFile(); err != nil {
		return fmt.Errorf("kv delete: %w", err)
	}
	return nil
}
