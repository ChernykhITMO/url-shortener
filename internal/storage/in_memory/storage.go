package in_memory

import (
	"sync"
)

type Storage struct {
	aliasToURL map[string]string
	urlToAlias map[string]string
	mux        sync.RWMutex
}

func New() *Storage {
	return &Storage{
		aliasToURL: make(map[string]string),
		urlToAlias: make(map[string]string),
	}
}
