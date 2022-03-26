package model

import "github.com/pthum/null"

// ColorProfile The ColorProfile which reflects color and brightness
type ColorProfile struct {
	ID int64 `json:"id,omitempty" gorm:"primary_key"`

	Blue null.Int `json:"blue,string,omitempty"`

	Brightness null.Int `json:"brightness,string,omitempty"`

	Green null.Int `json:"green,string,omitempty"`

	Red null.Int `json:"red,string,omitempty"`
}

// TableName sets the table name for the color profile
func (ColorProfile) TableName() string {
	return "color_profile"
}
