package models

import "github.com/google/uuid"

type Player struct {
	ID              string
	DisplayName     string
	PositionID      string
	ProgressPercent int
	HumanResource   int
	NumPower        int
	TotalCost       int
	IsGoaled        bool
}

func NewPlayer(displayName string, initPositionID string, numPower int) Player {

	if displayName == "" {
		panic("displayName must not be empty")
	} else if numPower < 2 {
		panic("numPower must not be smaller than 2")
	}

	return Player{
		ID:            uuid.Must(uuid.NewV7()).String(),
		DisplayName:   displayName,
		PositionID:    initPositionID,
		HumanResource: 1,
		NumPower:      numPower,
	}
}
