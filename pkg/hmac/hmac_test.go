package hmac_test

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/neurodyne-web-services/utils/pkg/hmac"
	"github.com/stretchr/testify/assert"
)

const (
	digest = "sign this message"
	key    = "my key"
	hash   = "6791a762f7568f945c2e1e396cea243e944100a6"
)

func Test_GenerateInvalid_GivesError(t *testing.T) {
	input := []byte("test")
	signature := "ab"
	secretKey := "key"
	err := hmac.Validate(input, signature, secretKey)
	if err == nil {
		t.Errorf("expected error when signature didn't have at least 5 characters in length")
		t.Fail()
		return
	}

	wantErr := "valid hash prefixes: [sha1=, sha256=], got: ab"
	if err.Error() != wantErr {
		t.Errorf("want: %s, got: %s", wantErr, err.Error())
		t.Fail()
	}
}

func Test_hmac(t *testing.T) {
	t.Run("ValidateWithoutSha1PrefixFails", func(t *testing.T) {
		valid := hmac.Validate([]byte(digest), hash, key)
		if valid == nil {
			t.Errorf("Expected error due to missing prefix")
			t.Fail()
		}
	})

	t.Run("ValidateWithSha1Prefix", func(t *testing.T) {
		encodedHash := "sha1=" + hash

		err := hmac.Validate([]byte(digest), encodedHash, key)
		assert.NoError(t, err)
	})

	t.Run("SignWithKey", func(t *testing.T) {
		wantHash := hash

		hash := hmac.Sign([]byte(digest), []byte(key), sha1.New)
		encodedHash := hex.EncodeToString(hash)
		assert.Equal(t, wantHash, encodedHash)
	})

	t.Run("SignWithKey_SHA256", func(t *testing.T) {
		wantHash := "41f8b7712c58dc25be8d30cf25e57739a65f5f2f449b59a42e04da1f191512e7"

		hash := hmac.Sign([]byte(digest), []byte(key), sha256.New)
		encodedHash := hex.EncodeToString(hash)
		assert.Equal(t, wantHash, encodedHash)
	})

	t.Run("ValidateWithSha256Prefix", func(t *testing.T) {
		encodedHash := "sha256=" + "41f8b7712c58dc25be8d30cf25e57739a65f5f2f449b59a42e04da1f191512e7"
		err := hmac.Validate([]byte(digest), encodedHash, key)
		assert.NoError(t, err)
	})
}
