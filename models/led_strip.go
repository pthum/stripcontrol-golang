package models

import "github.com/pthum/null"

// LedStrip definition of a LED strip and its configuration
type LedStrip struct {
	ID int64 `json:"id,omitempty" gorm:"primary_key"`

	Description string `json:"description,omitempty"`

	Enabled bool `json:"enabled,omitempty"`

	MisoPin null.Int `json:"misoPin,string,omitempty" gorm:"column:miso_pin"`

	Name string `json:"name,omitempty"`

	NumLeds null.Int `json:"numLeds,string,omitempty" gorm:"column:num_leds"`

	SclkPin null.Int `json:"sclkPin,string,omitempty" gorm:"column:sclk_pin"`

	SpeedHz null.Int `json:"speedHz,omitempty" gorm:"column:speed_hz"`

	ProfileID null.Int `json:"profileId,omitempty" gorm:"column:profile_id"`
}

// TableName sets the table name for the led strip
func (LedStrip) TableName() string {
	return "ledstrip"
}
