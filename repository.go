package gogit

import (
	"sync"
)

type Repository struct {
	logger   Logger
	location string
	mu       sync.RWMutex
}

type LogCallback struct {
	RanCommand func()
}
