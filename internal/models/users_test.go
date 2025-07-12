package models

import (
	"testing"

	"github.com/embracexyz/snippetbox/internal/assert"
)

func TestUserExists(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := UserModel{DB: db}
			exist, err := m.Exists(tt.userID)

			assert.Assert(t, exist, tt.want)
			assert.NilError(t, err)
		})
	}
}
