package uuid

import (
	"crypto/rand"
	"fmt"

	"lexbox/internal/domain/shared"
)

type IDGenerator struct{}

func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

func (g *IDGenerator) NewID() shared.ID {
	return shared.NewID(newUUID())
}

func newUUID() string {
	b := make([]byte, 16)

	_, err := rand.Read(b)
	if err != nil {
		panic(fmt.Errorf("cannot generate uuid: %w", err))
	}

	// Version 4
	b[6] = (b[6] & 0x0f) | 0x40
	// Variant RFC 4122
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		b[0:4],
		b[4:6],
		b[6:8],
		b[8:10],
		b[10:16],
	)
}
