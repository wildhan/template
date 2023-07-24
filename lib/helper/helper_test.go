package helper_test

import (
	"template/lib/helper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyStringToNull(t *testing.T) {
	testCases := []struct {
		name         string
		mockData     string
		expectResult interface{}
	}{
		{
			name:         "EmptyString",
			mockData:     "",
			expectResult: nil,
		},
		{
			name:         "NotEmptyString",
			mockData:     "HelloWorld",
			expectResult: "HelloWorld",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := helper.EmptyStringToNull(tc.mockData)

			assert.Equal(t, tc.expectResult, result)
		})
	}
}
