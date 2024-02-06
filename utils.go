package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
)

func isAlphanumeric(s string) bool {
	// Regular expression to match alphanumeric characters
	pattern := "^[a-zA-Z0-9]*$"

	// Compile the regular expression
	regex := regexp.MustCompile(pattern)

	// Test if the string matches the pattern
	return regex.MatchString(s)
}

// generateRandomString returns a random string of length n consisting of alphanumeric characters
func generateRandomString(n int) string {
	// Define the characters to choose from
	characters := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	// Create a slice to store the random bytes
	randomBytes := make([]byte, n)

	// Generate random bytes
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
		if err != nil {
			panic(err)
		}
		randomBytes[i] = characters[num.Int64()]
	}

	// Convert bytes to string and return
	return string(randomBytes)
}

// isUsernameValid returns error
func isUsernameValid(u string) error {
	if len(u) < 4 {
		return fmt.Errorf("username should have atleast 4 characters")
	}
	if !isAlphanumeric(u) {
		return fmt.Errorf("username should only have alphanumeric characters")
	}
	return nil
}
