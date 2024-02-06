package main

import (
	"crypto/rand"
	"math/big"
)

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
