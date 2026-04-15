package memory

import (
	"fmt"
	"sync/atomic"

	"lexbox/internal/domain/shared"
)

type IDGenerator struct {
	counter uint64
}

func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

func (g *IDGenerator) NewID() shared.ID {
	n := atomic.AddUint64(&g.counter, 1)
	return shared.NewID(fmt.Sprintf("id-%d", n))
}
