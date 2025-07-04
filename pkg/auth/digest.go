// Package auth provides HTTP Digest Authentication implementation
// for the TiDB Cloud SDK. It supports RFC 2617 compliant digest authentication
// with MD5 hashing and quality of protection (qop) handling.
package auth

import (
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// DigestAuth implements HTTP Digest Authentication according to RFC 2617.
// It handles the challenge-response authentication flow required by TiDB Cloud API.
type DigestAuth struct {
	realm     string // Authentication realm from server
	nonce     string // Server-provided nonce value
	qop       string // Quality of protection (typically "auth")
	opaque    string // Opaque value from server
	algorithm string // Hash algorithm (typically "MD5")
	nc        int    // Nonce count for replay protection
	cnonce    string // Client-generated nonce
}

// NewDigestAuth creates a new DigestAuth instance.
// The instance is initialized with a nonce count of 1 and is ready
// to parse authentication challenges from the server.
func NewDigestAuth() *DigestAuth {
	return &DigestAuth{
		nc: 1,
	}
}

// ParseChallenge parses an HTTP Digest authentication challenge from the server.
// It extracts the realm, nonce, qop, opaque, and algorithm values from the
// WWW-Authenticate header and prepares the client for response generation.
func (d *DigestAuth) ParseChallenge(authHeader string) error {
	if authHeader == "" {
		return errors.New("empty auth header")
	}

	if !strings.HasPrefix(authHeader, "Digest ") {
		return errors.New("not a digest auth header")
	}

	challengeData := strings.TrimPrefix(authHeader, "Digest ")

	// Parse key-value pairs
	pairs := parseKeyValuePairs(challengeData)

	d.realm = pairs["realm"]
	d.nonce = pairs["nonce"]
	d.qop = pairs["qop"]
	d.opaque = pairs["opaque"]
	d.algorithm = pairs["algorithm"]

	if d.algorithm == "" {
		d.algorithm = "MD5"
	}

	if d.realm == "" {
		return errors.New("missing realm in digest challenge")
	}
	if d.nonce == "" {
		return errors.New("missing nonce in digest challenge")
	}

	// Generate cnonce for this auth
	d.cnonce = generateCnonce()

	return nil
}

// GenerateAuthHeader generates the Authorization header value for HTTP Digest authentication.
// It creates the digest response using the provided credentials and request details,
// following the RFC 2617 specification for digest calculation.
//
// Parameters:
//   - username: The API public key
//   - password: The API private key
//   - method: HTTP method (GET, POST, etc.)
//   - uri: Request URI path
//
// Returns:
//   - string: Complete Authorization header value, or empty string if not ready
func (d *DigestAuth) GenerateAuthHeader(username, password, method, uri string) string {
	if d.nonce == "" {
		return ""
	}

	ha1 := d.generateHA1(username, password)
	ha2 := d.generateHA2(method, uri)

	var response string
	if d.qop == "auth" {
		response = d.generateResponseWithQop(ha1, ha2)
	} else {
		response = d.generateResponseWithoutQop(ha1, ha2)
	}

	var authHeader strings.Builder
	authHeader.WriteString("Digest ")
	authHeader.WriteString(fmt.Sprintf(`username="%s"`, username))
	authHeader.WriteString(fmt.Sprintf(`, realm="%s"`, d.realm))
	authHeader.WriteString(fmt.Sprintf(`, nonce="%s"`, d.nonce))
	authHeader.WriteString(fmt.Sprintf(`, uri="%s"`, uri))
	authHeader.WriteString(fmt.Sprintf(`, response="%s"`, response))

	if d.qop != "" {
		authHeader.WriteString(fmt.Sprintf(`, qop=%s`, d.qop))
		authHeader.WriteString(fmt.Sprintf(`, nc=%08x`, d.nc))
		authHeader.WriteString(fmt.Sprintf(`, cnonce="%s"`, d.cnonce))
	}

	if d.opaque != "" {
		authHeader.WriteString(fmt.Sprintf(`, opaque="%s"`, d.opaque))
	}

	if d.algorithm != "" {
		authHeader.WriteString(fmt.Sprintf(`, algorithm=%s`, d.algorithm))
	}

	return authHeader.String()
}

func (d *DigestAuth) generateHA1(username, password string) string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s:%s:%s", username, d.realm, password)))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (d *DigestAuth) generateHA2(method, uri string) string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s:%s", method, uri)))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (d *DigestAuth) generateResponseWithQop(ha1, ha2 string) string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s:%s:%08x:%s:%s:%s",
		ha1, d.nonce, d.nc, d.cnonce, d.qop, ha2)))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (d *DigestAuth) generateResponseWithoutQop(ha1, ha2 string) string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s:%s:%s", ha1, d.nonce, ha2)))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func parseKeyValuePairs(data string) map[string]string {
	pairs := make(map[string]string)

	// Use regex to find key="value" pairs
	re := regexp.MustCompile(`(\w+)=(?:"([^"]*)"|([^,\s]+))`)
	matches := re.FindAllStringSubmatch(data, -1)

	for _, match := range matches {
		key := match[1]
		value := match[2]
		if value == "" {
			value = match[3]
		}
		pairs[key] = value
	}

	return pairs
}

func generateCnonce() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
