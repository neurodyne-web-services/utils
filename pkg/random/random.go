// Package random contains different random generators.
package rand

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"math/rand"
	"time"
)

// Random generates a random int between min and max, inclusive.
func Random(min int, max int) int {
	return newRand().Intn(max-min+1) + min
}

// RandomInt picks a random element in the slice of ints.
func RandomInt(elements []int) int {
	index := Random(0, len(elements)-1)
	return elements[index]
}

// RandomString picks a random element in the slice of string.
func RandomString(elements []string) string {
	index := Random(0, len(elements)-1)
	return elements[index]
}

const base62chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const uniqueIDLength = 6 // Should be good for 62^6 = 56+ billion combinations

// UniqueId returns a unique (ish) id we can attach to resources and tfstate files so they don't conflict with each other
// Uses base 62 to generate a 6 character string that's unlikely to collide with the handful of tests we run in
// parallel. Based on code here: http://stackoverflow.com/a/9543797/483528
func UniqueID() string {
	var out bytes.Buffer

	generator := newRand()
	for i := 0; i < uniqueIDLength; i++ {
		out.WriteByte(base62chars[generator.Intn(len(base62chars))])
	}

	return out.String()
}

// newRand creates a new random number generator, seeding it with the current system time.
func newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func GenRandomName(pref string) string {
	return pref + "-" + UniqueID()
}

// GetHash - calculates a SHA1 hash 20 bytes for the input string.
func GetHash(d string) string {
	ab20 := sha1.Sum([]byte(d))
	return fmt.Sprintf("%x", ab20)
}
