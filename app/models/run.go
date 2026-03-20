package models

import "time"

type Run struct {
	CustomID    string       `bson:"customID" json:"id"`
	DungeonID   string       `bson:"dungeonID" json:"dungeonId"`
	PlayerID    string       `json:"playerId"`
	State       string       `json:"state"`
	CurrentStep int64        `json:"currentStep"`
	KilledSteps []KilledStep `json:"killedSteps"`
	StartedAt   time.Time    `json:"startedAt"`
	EndedAt     time.Time    `json:"endedAt"`
}

type KilledStep struct {
	BossStepID string    `json:"bossStepId"`
	KilledAt   time.Time `json:"killedAt"`
	Proof      string    `json:"proof"`
}

func (r *Run) Collection() string {
	return "run"
}