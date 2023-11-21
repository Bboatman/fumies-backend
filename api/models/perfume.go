package models

import (
	"time"

	"github.com/google/uuid"
)

type Perfume struct {
	Id          uuid.UUID `json:"id"`
	Name        *string   `json:"name,omitempty"`
	House       *string   `json:"house,omitempty"`
	Url         *string   `json:"url,omitempty"`
	IsEmpty     *bool     `json:"is_empty,omitempty"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type ModifyPerfumeBody struct {
	Name        string       `json:"name" binding:"required"`
	House       string       `json:"house" binding:"required"`
	Url         *string      `json:"url"`
	IsEmpty     bool         `json:"is_empty"`
	Description *string      `json:"description"`
	Notes       *[]uuid.UUID `json:"notes", db:"notes"`
}

type PerfumeResponse struct {
	Id          uuid.UUID    `json:"id"`
	Name        *string      `json:"name,omitempty"`
	House       *string      `json:"house,omitempty"`
	Url         *string      `json:"url,omitempty"`
	IsEmpty     *bool        `json:"is_empty,omitempty"`
	Description *string      `json:"description,omitempty"`
	Notes       *[]uuid.UUID `json:"notes,omitempty" db:"notes"`
}

type RecommendationRequestBody struct {
	Notes []string `json:"notes"`
}

type PerfumeVector struct {
	PerfumeId uuid.UUID `json:"perfume_id" db:"perfume_id"`
	Vector    []float64 `json:"vector"`
	Epoch     *float64  `json:"epoch,omitempty" db:"epoch"`
	CosineSim *float64  `json:"cosine_sim"`
	Include   *bool     `json:"include"`
}
