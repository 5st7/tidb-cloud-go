package auth

import (
	"testing"
)

func TestDigestAuth_ParseChallenge(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		expectedRealm  string
		expectedNonce  string
		expectedQop    string
		expectedOpaque string
		expectedAlg    string
		expectedErr    bool
	}{
		{
			name:           "valid digest challenge",
			authHeader:     `Digest realm="tidbcloud", nonce="dcd98b7102dd2f0e8b11d0f600bfb0c093", qop="auth", opaque="5ccc069c403ebaf9f0171e9517f40e41", algorithm="MD5"`,
			expectedRealm:  "tidbcloud",
			expectedNonce:  "dcd98b7102dd2f0e8b11d0f600bfb0c093",
			expectedQop:    "auth",
			expectedOpaque: "5ccc069c403ebaf9f0171e9517f40e41",
			expectedAlg:    "MD5",
			expectedErr:    false,
		},
		{
			name:        "empty auth header",
			authHeader:  "",
			expectedErr: true,
		},
		{
			name:        "invalid auth header",
			authHeader:  "Basic realm=test",
			expectedErr: true,
		},
		{
			name:        "missing realm",
			authHeader:  `Digest nonce="test"`,
			expectedErr: true,
		},
		{
			name:        "missing nonce",
			authHeader:  `Digest realm="test"`,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := &DigestAuth{}
			err := auth.ParseChallenge(tt.authHeader)

			if tt.expectedErr {
				if err == nil {
					t.Errorf("ParseChallenge() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseChallenge() unexpected error: %v", err)
				return
			}

			if auth.realm != tt.expectedRealm {
				t.Errorf("ParseChallenge() realm = %v, want %v", auth.realm, tt.expectedRealm)
			}
			if auth.nonce != tt.expectedNonce {
				t.Errorf("ParseChallenge() nonce = %v, want %v", auth.nonce, tt.expectedNonce)
			}
			if auth.qop != tt.expectedQop {
				t.Errorf("ParseChallenge() qop = %v, want %v", auth.qop, tt.expectedQop)
			}
			if auth.opaque != tt.expectedOpaque {
				t.Errorf("ParseChallenge() opaque = %v, want %v", auth.opaque, tt.expectedOpaque)
			}
			if auth.algorithm != tt.expectedAlg {
				t.Errorf("ParseChallenge() algorithm = %v, want %v", auth.algorithm, tt.expectedAlg)
			}
		})
	}
}

func TestDigestAuth_GenerateAuthHeader(t *testing.T) {
	tests := []struct {
		name       string
		username   string
		password   string
		method     string
		uri        string
		setupAuth  func(*DigestAuth)
		expectAuth bool
	}{
		{
			name:     "valid auth generation",
			username: "testuser",
			password: "testpass",
			method:   "GET",
			uri:      "/api/v1beta/projects",
			setupAuth: func(auth *DigestAuth) {
				auth.realm = "tidbcloud"
				auth.nonce = "dcd98b7102dd2f0e8b11d0f600bfb0c093"
				auth.qop = "auth"
				auth.opaque = "5ccc069c403ebaf9f0171e9517f40e41"
				auth.algorithm = "MD5"
			},
			expectAuth: true,
		},
		{
			name:     "missing nonce",
			username: "testuser",
			password: "testpass",
			method:   "GET",
			uri:      "/api/v1beta/projects",
			setupAuth: func(auth *DigestAuth) {
				auth.realm = "tidbcloud"
				auth.qop = "auth"
				auth.opaque = "5ccc069c403ebaf9f0171e9517f40e41"
				auth.algorithm = "MD5"
			},
			expectAuth: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := &DigestAuth{}
			if tt.setupAuth != nil {
				tt.setupAuth(auth)
			}

			authHeader := auth.GenerateAuthHeader(tt.username, tt.password, tt.method, tt.uri)

			if tt.expectAuth {
				if authHeader == "" {
					t.Error("GenerateAuthHeader() expected non-empty header but got empty")
				}
				if !containsDigestFields(authHeader) {
					t.Errorf("GenerateAuthHeader() header doesn't contain expected digest fields: %s", authHeader)
				}
			} else {
				if authHeader != "" {
					t.Errorf("GenerateAuthHeader() expected empty header but got: %s", authHeader)
				}
			}
		})
	}
}

func containsDigestFields(header string) bool {
	return len(header) > 0 &&
		header[:6] == "Digest" &&
		contains(header, "username=") &&
		contains(header, "realm=") &&
		contains(header, "nonce=") &&
		contains(header, "uri=") &&
		contains(header, "response=")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && containsAt(s, substr, 1)
}

func containsAt(s, substr string, start int) bool {
	if start >= len(s) {
		return false
	}
	if len(s[start:]) < len(substr) {
		return false
	}
	if s[start:start+len(substr)] == substr {
		return true
	}
	return containsAt(s, substr, start+1)
}
