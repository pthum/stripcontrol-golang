package model

import (
	"testing"

	"github.com/pthum/null"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/schema"
)

func TestStripJsonEncode(t *testing.T) {
	tests := []encodeTest[LedStrip]{
		{
			name: "test filled",
			input: LedStrip{
				BaseModel:   BaseModel{ID: 185},
				Description: "Test",
				Enabled:     false,
				MisoPin:     null.IntFrom(12),
				Name:        "Test",
				NumLeds:     null.IntFrom(5),
				SclkPin:     null.IntFrom(13),
				SpeedHz:     null.IntFrom(80000),
			},
			want: `{"id":185,"name":"Test","description":"Test","misoPin":12,"numLeds":5,"sclkPin":13,"speedHz":80000,"profileId":null}`,
		},
		{
			name:  "test empty",
			input: LedStrip{},
			want:  `{"misoPin":null,"numLeds":null,"sclkPin":null,"speedHz":null,"profileId":null}`,
		},
	}

	runEncodeTests(t, tests)
}

func TestStripJsonDecode(t *testing.T) {
	tests := []decodeTest[LedStrip]{
		{
			name:  "test filled",
			input: `{"id":185,"name":"Test","description":"Test","misoPin":12,"numLeds":5,"sclkPin":13,"speedHz":80000,"profileId":null}`,
			want: LedStrip{
				BaseModel:   BaseModel{ID: 185},
				Description: "Test",
				Enabled:     false,
				MisoPin:     null.IntFrom(12),
				Name:        "Test",
				NumLeds:     null.IntFrom(5),
				SclkPin:     null.IntFrom(13),
				SpeedHz:     null.IntFrom(80000),
			},
		},
		{
			name:  "test numbers as string",
			input: `{"id":185,"name":"Test","description":"Test","misoPin":"12","numLeds":"5","sclkPin":"13","speedHz":"80000","profileId":null}`,
			want: LedStrip{
				BaseModel:   BaseModel{ID: 185},
				Description: "Test",
				Enabled:     false,
				MisoPin:     null.IntFrom(12),
				Name:        "Test",
				NumLeds:     null.IntFrom(5),
				SclkPin:     null.IntFrom(13),
				SpeedHz:     null.IntFrom(80000),
			},
		},
		{
			name:  "test empty",
			input: `{"misoPin":null,"numLeds":null,"sclkPin":null,"speedHz":null,"profileId":null}`,
			want:  LedStrip{},
		},
	}

	runDecodeTests(t, tests)
}
func TestProfileJsonEncode(t *testing.T) {
	dummyProfile := ColorProfile{
		BaseModel:  BaseModel{ID: 185},
		Red:        null.IntFrom(123),
		Green:      null.IntFrom(234),
		Blue:       null.IntFrom(12),
		Brightness: null.IntFrom(1),
	}
	tests := []encodeTest[ColorProfile]{
		{
			name:  "test filled",
			input: dummyProfile,
			want:  `{"id":185,"blue":12,"brightness":1,"green":234,"red":123}`,
		},
		{
			name:  "test empty",
			input: ColorProfile{},
			want:  `{"blue":null,"brightness":null,"green":null,"red":null}`,
		},
	}

	runEncodeTests(t, tests)
}

func TestProfileJsonDecode(t *testing.T) {
	dummyProfile := ColorProfile{
		BaseModel:  BaseModel{ID: 185},
		Red:        null.IntFrom(123),
		Green:      null.IntFrom(234),
		Blue:       null.IntFrom(12),
		Brightness: null.IntFrom(1),
	}
	tests := []decodeTest[ColorProfile]{
		{
			name:  "test filled",
			input: `{"id":185,"blue":12,"brightness":1,"green":234,"red":123}`,
			want:  dummyProfile,
		},
		{
			name:  "test numbers as string",
			input: `{"id":185,"blue":"12","brightness":"1","green":"234","red":"123"}`,
			want:  dummyProfile,
		},
		{
			name:  "test empty",
			input: `{"blue":null,"brightness":null,"green":null,"red":null}`,
			want:  ColorProfile{},
		},
	}
	runDecodeTests(t, tests)
}

func TestTableName(t *testing.T) {

	tests := []struct {
		name  string
		input schema.Tabler
		want  string
	}{
		{
			name:  "test profile",
			input: LedStrip{},
			want:  "ledstrip",
		},
		{
			name:  "test profile",
			input: ColorProfile{},
			want:  "color_profile",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.input.TableName())
		})
	}
}

func TestGenerateID(t *testing.T) {
	const maxID int64 = 500
	const minID int64 = 0
	bm := BaseModel{}
	for i := 1; i <= 2000; i++ {
		bm.GenerateID()
		assert.GreaterOrEqual(t, bm.GetID(), minID)
		assert.LessOrEqual(t, bm.GetID(), maxID)
	}
}
