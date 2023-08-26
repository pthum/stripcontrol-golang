package model

import (
	"encoding/json"

	"github.com/pthum/null"
)

// StripEvent a StripEvent
type StripEvent struct {
	Type  EventType `json:"type,omitempty"`
	ID    null.Int  `json:"id,omitempty"`
	Strip OptStrip  `json:"state,omitempty"`
}

func NewStripEvent(id null.Int, typ EventType) *StripEvent {
	evnt := &StripEvent{
		Type: typ,
		ID:   id,
	}
	return evnt
}

func (pe *StripEvent) With(strip *LedStrip) *StripEvent {
	if strip != nil {
		pe.Strip.Valid = true
		pe.Strip.Strip.ID = strip.ID
		pe.Strip.Strip.Name = strip.Name
		pe.Strip.Strip.Enabled = strip.Enabled
		pe.Strip.Strip.MisoPin = strip.MisoPin.Int64
		pe.Strip.Strip.SclkPin = strip.SclkPin.Int64
		pe.Strip.Strip.NumLeds = strip.NumLeds.Int64
		pe.Strip.Strip.SpeedHz = strip.SpeedHz.Int64
	}
	return pe
}
func (pe *OptStrip) With(profile ColorProfile) *OptStrip {
	pe.Strip.Profile.Valid = true
	pe.Strip.Profile.Profile = profile
	return pe
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

// OptProfile Optional Profile
type OptProfile struct {
	Profile ColorProfile
	Valid   bool
}

// ProfileEvent a StripEvent
type ProfileEvent struct {
	Type  EventType  `json:"type,omitempty"`
	ID    null.Int   `json:"id,omitempty"`
	State OptProfile `json:"state,omitempty"`
}

func NewProfileEvent(id null.Int, typ EventType) *ProfileEvent {
	evnt := &ProfileEvent{
		Type: typ,
		ID:   id,
	}
	evnt.State.Valid = false
	return evnt
}

func (pe *ProfileEvent) With(profile ColorProfile) *ProfileEvent {
	pe.State.Valid = true
	pe.State.Profile = profile
	return pe
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

//go:generate enumer -type=EventType -json -text -transform=upper
type EventType int

const (
	Unknown EventType = iota
	Save
	Delete
)
