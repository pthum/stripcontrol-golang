package model

import (
	"fmt"
	"testing"

	"github.com/pthum/stripcontrol-golang/internal/testutils"
	"github.com/stretchr/testify/assert"
)

type encodeTest[T any] struct {
	name  string
	input T
	want  string
}

type decodeTest[T any] struct {
	name  string
	input string
	want  T
}

func runEncodeTests[T any](t *testing.T, testCases []encodeTest[T]) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := testutils.JsonEncode(t, tc.input)
			fmt.Println(result)
			assert.Equal(t, tc.want, result)

		})
	}
}

func runDecodeTests[T any](t *testing.T, testCases []decodeTest[T]) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result T
			testutils.JsonDecode(t, tc.input, &result)
			fmt.Println(result)
			assert.Equal(t, tc.want, result)

		})
	}
}
