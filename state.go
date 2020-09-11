package main

// GameState is menu / playing / paused / game over
type GameState int

// Current state
const (
	StateMenu GameState = iota
	StatePlaying
	StatePaused
	StateGameOver
)
