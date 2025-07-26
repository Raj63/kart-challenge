package models

// Migration represents a database migration record.
type Migration struct {
	ID          string `bson:"id" json:"id"`                   // Unique identifier for the migration
	Version     string `bson:"version" json:"version"`         // Migration version (e.g., "0001")
	Description string `bson:"description" json:"description"` // Description of what the migration does
	AppliedAt   int64  `bson:"applied_at" json:"applied_at"`   // Unix timestamp when migration was applied
	Status      string `bson:"status" json:"status"`           // Status of the migration (e.g., "completed", "failed")
}

// MigrationData represents the structure of a migration JSON file.
type MigrationData struct {
	Version     string    `json:"version"`     // Migration version
	Description string    `json:"description"` // Migration description
	Products    []Product `json:"products"`    // Products to seed
}
