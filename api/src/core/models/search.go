package models

// ------------------------------------------------------------
// : Search
// ------------------------------------------------------------
type Search struct {
	ID        uint64    `json:"id"`
	Token     string    `json:"token"`
	Timestamp string    `json:"timestamp"`
	Metadata  JSONBMap  `json:"metadata"`
}


