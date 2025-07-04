package client

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		publicKey   string
		privateKey  string
		expectedErr bool
	}{
		{
			name:        "valid keys",
			publicKey:   "test_public_key",
			privateKey:  "test_private_key",
			expectedErr: false,
		},
		{
			name:        "empty public key",
			publicKey:   "",
			privateKey:  "test_private_key",
			expectedErr: true,
		},
		{
			name:        "empty private key",
			publicKey:   "test_public_key",
			privateKey:  "",
			expectedErr: true,
		},
		{
			name:        "both keys empty",
			publicKey:   "",
			privateKey:  "",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.publicKey, tt.privateKey)

			if tt.expectedErr {
				if err == nil {
					t.Errorf("NewClient() expected error but got none")
				}
				if client != nil {
					t.Errorf("NewClient() expected nil client but got %v", client)
				}
				return
			}

			if err != nil {
				t.Errorf("NewClient() unexpected error: %v", err)
				return
			}

			if client == nil {
				t.Error("NewClient() returned nil client")
				return
			}

			if client.baseURL != DefaultBaseURL {
				t.Errorf("NewClient() baseURL = %v, want %v", client.baseURL, DefaultBaseURL)
			}
		})
	}
}