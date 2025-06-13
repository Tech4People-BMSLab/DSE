package models

// ------------------------------------------------------------
// : Website
// ------------------------------------------------------------
type Website struct {
	Name  string `json:"name"`
	Query string `json:"query"`
	Url   string `json:"url"`
}
