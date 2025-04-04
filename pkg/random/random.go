// Package random contains different random generators.
package random

import (
	"bytes"
	"crypto/md5"
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"time"
)

const base62chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const specialChars = "_,.?-~@#$%^+-&*()=\\/<>`"

var AllChars = fmt.Sprintf("%s%s", base62chars, specialChars)

// Random generates a random int between min and max, inclusive.
func Random(min int, max int) int {
	return newRand().Intn(max-min+1) + min
}

// RandomString picks a random element in the slice of string.
func RandomItem[T any](elements []T) T {
	index := Random(0, len(elements)-1)
	return elements[index]
}

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

// GenRandomName - generates a random name for an input key, e.g foo-sEncsH.
func GenRandomName(pref string) string {
	return pref + "-" + UniqueID()
}

// GenRandomName - generates a random name for an input key, e.g foo-sencsh, all lower case.
func GenRandomNameLower(pref string) string {
	return strings.ToLower(pref + "-" + UniqueID())
}

// GenerateRandomByteString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomByteString(n int, letters string) ([]byte, error) {
	ret := make([]byte, n)
	for i := range n {
		num, err := crand.Int(crand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return nil, err
		}
		ret[i] = letters[num.Int64()]
	}

	return ret, nil
}

// GenerateRandomString - same as byte string, but returns a pretty string.
func GenerateRandomString(n int) (string, error) {
	out, err := GenerateRandomByteString(n, base62chars)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// GenerateRandomString - same as byte string, but returns a pretty string.
func GenerateRandomPassword(n int) (string, error) {
	out, err := GenerateRandomByteString(n, AllChars)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Gethash - calculates string hash with a trimmed length.
func GetHash(s string, length int) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:length])
}
