package auth

import "testing"

func TestHashingUtils(t *testing.T) {
	testCases := []struct {
		passwordOne    string
		passwordTwo    string
		passwordsMatch bool
	}{
		{"hello", "hello", true},
		{"hello", "bye", false},
	}
	for _, tc := range testCases {
		hashedOne, err := HashPassword(tc.passwordOne)
		if err != nil {
			t.Errorf("Not expecting any error while hashing, got %s", err.Error())
		} else {
			comparison := CompareHashToPassword(tc.passwordTwo, hashedOne)
			if tc.passwordsMatch != comparison {
				t.Errorf("Expecting password-to-hash comparison to yield %v, got %v", tc.passwordsMatch, comparison)
			}
		}
	}
}

func TestGenerateToken(t *testing.T) {
	testCases := []struct {
		tokenLength          int
		expectedStringLength int
	}{
		{32, 44},
		{16, 24},
		{48, 64},
	}
	for _, tc := range testCases {
		token, err := GenerateToken(tc.tokenLength)
		if err != nil {
			t.Errorf("Not expecting an error when generating a new token, got %s", err.Error())
		}
		if len(token) != tc.expectedStringLength {
			t.Errorf("Expecting base64-encoded token to be of length %d, got %d", tc.expectedStringLength, len(token))
		}
	}
}
