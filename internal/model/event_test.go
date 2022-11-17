package model

import (
	"testing"

	"github.com/pthum/null"
)

func TestStripEventJson(t *testing.T) {
	stripSave := *NewStripEvent(null.IntFrom(234), Save).With(LedStrip{
		BaseModel: BaseModel{ID: 234},
		Name:      "test",
	})
	stripSave.Strip.With(ColorProfile{
		BaseModel:  BaseModel{ID: 185},
		Red:        null.IntFrom(123),
		Green:      null.IntFrom(234),
		Blue:       null.IntFrom(12),
		Brightness: null.IntFrom(1),
	})
	tests := []encodeTest[StripEvent]{
		{
			name:  "save stripevent",
			input: stripSave,
			want:  `{"type":"SAVE","id":234,"state":{"id":234,"name":"test","profile":{"id":185,"blue":12,"brightness":1,"green":234,"red":123}}}`,
		},
		{
			name: "save stripevent without profile",
			input: *NewStripEvent(null.IntFrom(234), Save).With(LedStrip{
				BaseModel: BaseModel{ID: 234},
				Name:      "test",
			}),
			want: `{"type":"SAVE","id":234,"state":{"id":234,"name":"test","profile":null}}`,
		},
		{
			name:  "delete stripevent",
			input: *NewStripEvent(null.IntFrom(234), Delete),
			want:  `{"type":"DELETE","id":234,"state":null}`,
		},
	}

	runEncodeTests(t, tests)
}

func TestProfileEventJson(t *testing.T) {
	dummyProfile := ColorProfile{
		BaseModel:  BaseModel{ID: 185},
		Red:        null.IntFrom(123),
		Green:      null.IntFrom(234),
		Blue:       null.IntFrom(12),
		Brightness: null.IntFrom(1),
	}
	tests := []encodeTest[ProfileEvent]{
		{
			name:  "save profileevent",
			input: *NewProfileEvent(null.IntFrom(123), Save).With(dummyProfile),
			want:  `{"type":"SAVE","id":123,"state":{"id":185,"blue":12,"brightness":1,"green":234,"red":123}}`,
		},
		{
			name:  "delete profileevent",
			input: *NewProfileEvent(null.IntFrom(123), Delete),
			want:  `{"type":"DELETE","id":123,"state":null}`,
		},
	}

	runEncodeTests(t, tests)
}
