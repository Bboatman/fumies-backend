package models

import (
	"time"

	"github.com/google/uuid"
)

type Metric struct {
	Id       uuid.UUID
	Label    string
	category string
	G        float64
	T        float64
	F        float64
	M        float64
}

type PerfumeMetric struct {
	Id        uuid.UUID `json:"id,omitempty"`
	NoteId    uuid.UUID `json:"note_id" binding:"required"`
	PerfumeId uuid.UUID `json:"perfume_id" binding:"required"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
