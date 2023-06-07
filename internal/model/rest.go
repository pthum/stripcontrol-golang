package model

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/pthum/null"
)

const (
	Table_ColorProfile = "color_profile"
	Table_LedStrip     = "ledstrip"
)

type IDer interface {
	GetID() int64
	TableName() string
}

type BaseModel struct {
	ID int64 `json:"id,omitempty" gorm:"primary_key" csv:"id"`
}

func (b *BaseModel) GetID() int64 {
	return b.ID
}
func (b *BaseModel) GetNullID() null.Int {
	return null.NewInt(b.ID, true)
}
func (b *BaseModel) GetStringID() string {
	return strconv.FormatInt(b.ID, 10)
}

// GenerateID generates an ID between 0 and 500
func (b *BaseModel) GenerateID() {
	// we have to generate an ID, as in contrast to spring/quarkus, gorm does not provide the generate-id on ORM side
	// (as we do not want to migrate the DB to keep compatibility between the different implementations)
	// we keep this logic simple and just generate a random int64. we do not expect too much strips, so the chance of a collision should be low
	s1 := rand.NewSource(time.Now().UnixNano())
	//#nosec G404 - false positive, math random is okay here
	r1 := rand.New(s1)
	b.ID = r1.Int63n(500)
}

// ColorProfile The ColorProfile which reflects color and brightness
type ColorProfile struct {
	BaseModel
	Blue       null.Int `json:"blue,omitempty" csv:"blue"`
	Brightness null.Int `json:"brightness,omitempty" csv:"brightness"`
	Green      null.Int `json:"green,omitempty" csv:"green"`
	Red        null.Int `json:"red,omitempty" csv:"red"`
}

// TableName sets the table name for the color profile
func (ColorProfile) TableName() string {
	return Table_ColorProfile
}

// LedStrip definition of a LED strip and its configuration
type LedStrip struct {
	BaseModel
	Name        string   `json:"name,omitempty" csv:"name"`
	Description string   `json:"description,omitempty" csv:"description"`
	Enabled     bool     `json:"enabled,omitempty" csv:"enabled"`
	MisoPin     null.Int `json:"misoPin,omitempty" gorm:"column:miso_pin" csv:"miso_pin"`
	NumLeds     null.Int `json:"numLeds,omitempty" gorm:"column:num_leds" csv:"num_leds"`
	SclkPin     null.Int `json:"sclkPin,omitempty" gorm:"column:sclk_pin" csv:"sclk_pin"`
	SpeedHz     null.Int `json:"speedHz,omitempty" gorm:"column:speed_hz" csv:"speed_hz"`
	ProfileID   null.Int `json:"profileId,omitempty" gorm:"column:profile_id" csv:"profile_id"`
}

// TableName sets the table name for the led strip
func (LedStrip) TableName() string {
	return Table_LedStrip
}
