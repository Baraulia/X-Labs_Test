package validation

import "testing"

func TestIsEmail(t *testing.T) {
	validator := New()
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"test.email@example.co.uk", true},
		{"test123", false},
		{"test@", false},
		{"@example.com", false},
	}

	for _, test := range tests {
		result := validator.IsEmail(test.email)
		if result != test.expected {
			t.Errorf("Expected IsEmail(%s) to be %v, but got %v", test.email, test.expected, result)
		}
	}
}
