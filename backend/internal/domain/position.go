package domain

type PositionSide string

const (
	PositionSideYes PositionSide = "yes"
	PositionSideNo  PositionSide = "no"
)

type PositionStatus string

const (
	PositionStatusOpen      PositionStatus = "open"
	PositionStatusWon       PositionStatus = "won"
	PositionStatusLost      PositionStatus = "lost"
	PositionStatusCancelled PositionStatus = "cancelled"
)
