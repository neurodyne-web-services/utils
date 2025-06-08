package id_test

import (
	"fmt"
	"testing"

	"github.com/neurodyne-web-services/utils/pkg/id"
	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
)

func Test_id(t *testing.T) {
	ident := uuid.NewString()
	fmt.Printf("uuid: %s\n", ident)

	ulid, err := id.ToULID(ident)
	assert.NoError(t, err)

	fmt.Printf("ulid: %s\n", ulid)

	newID, err := id.ToUUID(ulid)
	assert.NoError(t, err)

	assert.Equal(t, ident, newID)
}
