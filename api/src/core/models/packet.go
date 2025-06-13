package models

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ------------------------------------------------------------
// : Packet
// ------------------------------------------------------------
type Packet struct {
	Version string      `json:"version" msgpack:"version"`
	From    string      `json:"from"    msgpack:"from"`
	To      string      `json:"to"      msgpack:"to"`
	Action  string      `json:"action"  msgpack:"action"`
	Data    interface{} `json:"data"    msgpack:"data"`
}

// ------------------------------------------------------------
// : Constructor
// ------------------------------------------------------------
func NewPacket(version, from, to, action string, data interface{}) (*Packet, error) {
	if version   == "" { return nil, fmt.Errorf("Version is required") }
	if from      == "" { return nil, fmt.Errorf("From is required") }
	if to        == "" { return nil, fmt.Errorf("To is required") }
	if action    == "" { return nil, fmt.Errorf("Action is required") }

	p := &Packet{
		Version: version,
		From   : from,
		To     : to,
		Action : action,
		Data   : data,
	}

	return p, nil
}

// ------------------------------------------------------------
// : Setters
// ------------------------------------------------------------
func (p *Packet) SetVersion(version string) {
	p.Version = version
}

func (p *Packet) SetFrom(from string) {
	p.From = from
}

func (p *Packet) SetTo(to string) {
	p.To = to
}

func (p *Packet) SetAction(action string) {
	p.Action = action
}

func (p *Packet) SetData(data interface{}) {
	p.Data = data
}

// ------------------------------------------------------------
// : Serialize
// ------------------------------------------------------------
func (p *Packet) ToJSON() string {
    if p.Data == nil { p.Data = map[string]string{} }

    buffer  := &bytes.Buffer{}
    encoder := json.NewEncoder(buffer)
    encoder.SetEscapeHTML(false) // Disable HTML escaping

    err := encoder.Encode(p)
    if err != nil {
        return ""
    }
    return buffer.String()
}
