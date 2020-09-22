package utiles

import (
	"github.com/google/uuid"
	"testing"
)

func TestIsValidUUID(t *testing.T) {

	tests := []struct {
		name string
		uuid string
		isValid bool
	}{ {
		name: "test valid uuid4",
		uuid: uuid.New().String(),
		isValid: true,
	   },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if IsValidUUID(tt.uuid) != tt.isValid {
				t.Errorf("IsValidUUID(%s) returnd %v, we want %v", tt.uuid, !tt.isValid, tt.isValid)
			}
		})
	}
}
