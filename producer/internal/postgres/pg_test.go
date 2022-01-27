package postgres_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wednesday-solutions/go-template-producer/internal/postgres"
)

func TestConnect(t *testing.T) {
	db, err := postgres.Connect()
	if err != nil {
		assert.NotNil(t, db)
	}
	db.Close()
}
