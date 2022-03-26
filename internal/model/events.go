package model

import (
	"encoding/json"

	"github.com/pthum/null"
)

// StripEvent a StripEvent
type StripEvent struct {
	Type  string   `json:"type,omitempty"`
	ID    null.Int `json:"id,omitempty"`
	Strip OptStrip `json:"state,omitempty"`
}

// OptStrip Optional Strip
type OptStrip struct {
	Valid bool
	Strip struct {
		ID      int64      `json:"id,omitempty"`
		Name    string     `json:"name,omitempty"`
		Enabled bool       `json:"enabled,omitempty"`
		MisoPin int64      `json:"misoPin,omitempty"`
		SclkPin int64      `json:"sclkPin,omitempty"`
		NumLeds int64      `json:"numLeds,omitempty"`
		SpeedHz int64      `json:"speedHz,omitempty"`
		Profile OptProfile `json:"profile,omitempty"`
	}
}

//OptProfile Optional Profile
type OptProfile struct {
	Profile ColorProfile
	Valid   bool
}

// ProfileEvent a StripEvent
type ProfileEvent struct {
	Type  string     `json:"type,omitempty"`
	ID    null.Int   `json:"id,omitempty"`
	State OptProfile `json:"state,omitempty"`
}

// MarshalJSON marshals json for OptProfile
func (profile OptProfile) MarshalJSON() (data []byte, err error) {
	if profile.Valid {
		data, err = json.Marshal(profile.Profile)
		return
	}
	err = nil
	data = []byte("null")
	return
}

// MarshalJSON marshals json for OptProfile
func (strip OptStrip) MarshalJSON() (data []byte, err error) {
	if strip.Valid {
		data, err = json.Marshal(strip.Strip)
		return
	}
	err = nil
	data = []byte("null")
	return
}
