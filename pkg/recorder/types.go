package recorder

import "github.com/google/uuid"

type StubRecord struct {
	ID    uuid.UUID
	Count int
}
