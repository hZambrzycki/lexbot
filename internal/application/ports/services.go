package ports

import (
	"time"

	"lexbox/internal/domain/shared"
)

type IDGenerator interface {
	NewID() shared.ID
}

type Clock interface {
	Now() time.Time
}
