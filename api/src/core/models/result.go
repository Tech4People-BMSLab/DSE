package models

// ------------------------------------------------------------
// : Result
// ------------------------------------------------------------
type Result struct {
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`

	Url          string `json:"url"`
	Localization string `json:"localization"`

	Keyword string `json:"keyword"`
	Website string `json:"website"`

	Html string `json:"html"`
}
