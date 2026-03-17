package models

import "time"

type ItemID string

type InventoryEntry struct {
	PlayerID  string    `json:"player_id"`
	ItemID    string    `json:"item_id"`
	Qty       int64     `json:"qty"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Collection Mongodb collection
func (ie *InventoryEntry) Collection() string {
	return "inventory"
}

type ItemDef struct {
	CustomID    string    `bson:"customID" json:"id"`
	Type        string    `json:"type"`
	Rarity      string    `json:"rarity"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Stats   map[string]any    `json:"stats_json"`
	Tradable    bool      `json:"tradable"`
	BaseValue   int64     `json:"base_value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Collection Mongodb collection
func (i *ItemDef) Collection() string {
	return "item"
}

type InventoryResponse struct {
	PlayerID string             `json:"playerId"`
	Items    []InventoryItemDTO `json:"items"`
}

type InventoryItemDTO struct {
	ItemID string `json:"itemId"`
	Qty    int64  `json:"qty"`
}

type ItemDefResponse struct {
	ID          string         `json:"id"`
	Type        string         `json:"type"`
	Rarity      string         `json:"rarity"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Stats       map[string]any `json:"stats,omitempty"`
	Tradable    bool           `json:"tradable"`
	BaseValue   int64          `json:"baseValue,omitempty"`
}
