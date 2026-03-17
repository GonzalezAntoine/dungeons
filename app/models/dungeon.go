package models

import "time"

type Dungeon struct {
	CustomID    string    `bson:"customID" json:"id"`
	Title	   string    `json:"title"`
	Description string    `json:"description"`
	CreatedBy  string    `json:"createdBy"`
	Area	   string    `json:"area"`
	Boss  []BossStep `json:"boss"`
	Status	 string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BossStep struct {
	CustomID    string    `bson:"customID" json:"id"`
	DungeonID   string    `bson:"dungeonID" json:"dungeonId"`
	Order	   int64       `json:"order"`
	Name 	 string    `json:"name"`
	Location Location    `json:"location"`
	ZoneDescription string    `json:"zoneDescription"`
	Difficulty int64       `json:"difficulty"`
	Rewards  string    `json:"rewards"`
	BossState string    `json:"bossState"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	RadiusMeters float64 `json:"radiusMeters"`
}

type Run struct {
	CustomID    string    `bson:"customID" json:"id"`
	DungeonID   string    `bson:"dungeonID" json:"dungeonId"`
	PlayerID      string    `json:"playerId"`
	State string `json:"state"`
	CurrentStep int64 `json:"currentStep"`
	KilledSteps []KilledStep `json:"killedSteps"`
	StartedAt time.Time `json:"startedAt"`
	EndedAt time.Time `json:"endedAt"`
}

type KilledStep struct {
	BossStepID string `json:"bossStepId"`
	KilledAt time.Time `json:"killedAt"`
	Proof string `json:"proof"`
}

func (d *Dungeon) Collection() string {
	return "dungeon"
}