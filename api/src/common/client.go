package common

import "time"

// ------------------------------------------------------------
// : Client
// ------------------------------------------------------------
type Client struct {
    ID        string `json:"id"        bson:"_id"`

    Token     string `json:"token"     bson:"token"`
    Connected bool   `json:"connected" bson:"connected"`

    Form     map[string]any `json:"form"     bson:"form"`
    Manifest map[string]any `json:"manifest" bson:"manifest"`
    Storage  map[string]any `json:"storage"  bson:"storage"`

    LastPing   time.Time `json:"last_ping"   bson:"last_ping"`
    LastSearch time.Time `json:"last_search" bson:"last_search"`
    
    Version  string `json:"version"  bson:"version"`
}
