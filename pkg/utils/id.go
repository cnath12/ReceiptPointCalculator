package utils

import (
    "github.com/google/uuid"
    "sync"
)

var (
    idGenerator *IDGenerator
    once        sync.Once
)

type IDGenerator struct {
    mu sync.Mutex
}

func GetIDGenerator() *IDGenerator {
    once.Do(func() {
        idGenerator = &IDGenerator{}
    })
    return idGenerator
}

func (g *IDGenerator) GenerateID() string {
    g.mu.Lock()
    defer g.mu.Unlock()
    return uuid.New().String()
}

// Helper function for direct usage
func GenerateID() string {
    return GetIDGenerator().GenerateID()
}