package id

import (
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

func ToULID(uuidString string) (string, error) {
	var uuidID uuid.UUID
	var err error
	var uuidBinary []byte
	var ulidID ulid.ULID

	if uuidID, err = uuid.Parse(uuidString); err != nil {
		return "", err
	}
	uuidBinary, _ = uuidID.MarshalBinary()
	_ = ulidID.UnmarshalBinary(uuidBinary)

	return ulidID.String(), nil
}

func ToUUID(ulidString string) (string, error) {
	var ulidID ulid.ULID
	var uuidID uuid.UUID
	var err error
	var ulidBinary []byte

	if ulidID, err = ulid.Parse(ulidString); err != nil {
		return "", err
	}

	ulidBinary, _ = ulidID.MarshalBinary()
	_ = uuidID.UnmarshalBinary(ulidBinary)

	return uuidID.String(), nil
}
