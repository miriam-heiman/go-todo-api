package auth

import (
	"crypto/hmac"     // HMAC signature creation for JWTs
	"crypto/sha256"   // SHA-256 hashing algorithm that works with HMAC to create secure signatures
	"encoding/base64" // Base64 encoding/decoding, converts binary data to/from JWT base64 text format
	"encoding/json"   // JSON encoding/decoding, converts Go structs to/from header and payload JSON objects
	"errors"          // Error creation - returns errors for invalid tokens, expired tokens etc.
	"os"              // Operating System Functions - to get the JWT secret key from environment variables
	"strings"         // String operations like split, trim, contains etc. It is needed to parse JWT tokens
	"time"            // Time/Date operations to set token expiration time and check if token has expired

	"go-todo-api/internal/models"
)

// ============================================================================
// JWT TOKEN GENERATION AND VALIDATION
// ============================================================================

// TokenExpiration is how long tokens are valid (1 hour)
const TokenExpiration = time.Hour

// getSecretKey returns the secret key used to sign JWTs
// This key should be kept secret and stored in environment variables
func getSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-this-in-production"
	}
	return []byte(secret)
}

// GenerateToken creates a new JWT token for a user
// Steps:
// 1. Create header (algorithm info)
// 2. Create payload (user data + expiration)
// 3. Create signature (proves token  hasn't been tampered with)
// 4. Combine: header.payload.signature
func GenerateToken(userID, username string) (string, error) {
	// Step 1: CREATE HEADER
	// Header tells what algorithm we're using
	header := map[string]string{
		"alg": "HS256", // HMAC with SHA-256
		"typ": "JWT",   // Token type
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	// Base64 encode the header (makes it URL-safe)
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)

	// Step 2: CREATE PAYLOAD (CLAIMS)
	// Payload contains the actual data about the user
	now := time.Now()
	claims := models.JWTClaims{
		UserID:   userID,
		Username: username,
		Iat:      now.Unix(),                      // Issued at (now)
		Exp:      now.Add(TokenExpiration).Unix(), // Expires in 1 hour
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	// Base64 encode the payload
	claimsEncoded := base64.RawURLEncoding.EncodeToString(claimsJSON)

	// Step 3: CREATE SIGNATURE
	// Signature proves the token hasn't been tampered with
	// It's created by: HMAC-SHA256(header.payload, secret_key)
	message := headerEncoded + "." + claimsEncoded
	signature := createSignature(message, getSecretKey())
	signatureEncoded := base64.RawURLEncoding.EncodeToString(signature)

	// Step 4: COMBINE ALL PARTS
	// Final JWT: header.payload.signature
	token := message + "." + signatureEncoded

	return token, nil
}

// ValidateToken checks if a JWT token is valid
// Checks:
// 1. Is token format correct? (3 parts separated by dots)
// 2. Is signature valid? (not tampered with)
// 3. Is token expired? (check exp field)
func ValidateToken(tokenString string) (*models.JWTClaims, error) {
	// Step 1: SPLIT TOKENS INTO PARTS
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	headerEncoded := parts[0]
	claimsEncoded := parts[1]
	signatureEncoded := parts[2]

	// Step 2: VERIFY SIGNATURE
	// Recreate the signature and compare with the one in the token
	message := headerEncoded + "." + claimsEncoded
	expectedSignature := createSignature(message, getSecretKey())
	expectedSignatureEncoded := base64.RawURLEncoding.EncodeToString(expectedSignature)

	// Compare signatures (constant time to prevent timing attacks)
	if !hmac.Equal([]byte(signatureEncoded), []byte(expectedSignatureEncoded)) {
		return nil, errors.New("invalid token signature")
	}

	// Step 3: DECODE PAYLOAD
	// Base64 decode the claims
	claimsJSON, err := base64.RawURLEncoding.DecodeString(claimsEncoded)
	if err != nil {
		return nil, errors.New("invalid token encoding")
	}

	// Parse JSON into claims struct
	var claims models.JWTClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, errors.New("invalid token claims")
	}

	// Step 4: CHECK EXPIRATION
	// Is the token expired?
	now := time.Now().Unix()
	if claims.Exp < now {
		return nil, errors.New("token expired")
	}

	// Token is valid!
	return &claims, nil
}

// createSignature creates HMAC-SHA256 signature
// This is the cryptographic function that proves the token hasn't been tampered with
func createSignature(message string, secret []byte) []byte {
	// Create HMAC hasher with SHA256
	h := hmac.New(sha256.New, secret)

	// Write the message to the hasher
	h.Write([]byte(message))

	// Get the final hash (this is our signature)
	return h.Sum(nil)
}

// ValidateCredentials checks if username and password are correct
func ValidateCredentials(username, password string) (string, error) {
	validUsers := map[string]string{
		"john":  "password123", // username: password
		"admin": "admin123",
	}

	// Check if username exists and password matches
	if validPassword, exists := validUsers[username]; exists {
		if validPassword == password {
			// Password correct! Return user ID
			// In production, this would be the MongoDB ObjectID
			return "user_" + username, nil
		}
	}

	// Invalid credentials
	return "", errors.New("invalid username or password")
}
