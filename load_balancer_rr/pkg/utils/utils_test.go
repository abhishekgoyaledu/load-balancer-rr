package utils

import (
	"testing"
)

func TestValidatePort(t *testing.T) {
	tests := []struct {
		port       int
		expectErr  bool
		errMessage string
	}{
		{0, true, "port must be between 1 and 65535"},     // Invalid port - too low
		{65536, true, "port must be between 1 and 65535"}, // Invalid port - too high
		{8082, false, ""}, // Valid port number
	}

	for _, tt := range tests {
		err := ValidatePort(tt.port)

		if tt.expectErr {
			if err == nil {
				t.Errorf("Expected an error for port %d, got nil", tt.port)
			} else if err.Error() != tt.errMessage {
				t.Errorf("Expected error message '%s', got '%s'", tt.errMessage, err.Error())
			}
		} else {
			if err != nil {
				t.Errorf("Did not expect an error for port %d, got '%v'", tt.port, err)
			}
		}
	}
}
