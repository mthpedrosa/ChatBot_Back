package requests

import (
	"errors"
	"time"
)

type CostParams struct {
	MetaId    string `json:"meta_id" validate:"required"`
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required"`
}

// Convertendo as datas no serviço
func (c *CostParams) ParseDates() (time.Time, time.Time, error) {
	start, err := time.Parse(time.RFC3339, c.StartDate)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("invalid start date format")
	}

	end, err := time.Parse(time.RFC3339, c.EndDate)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("invalid end date format")
	}

	return start, end, nil
}
