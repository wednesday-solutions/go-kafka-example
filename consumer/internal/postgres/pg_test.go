package postgres_test

import (
	"consumer/internal/postgres"
	. "github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConnect(t *testing.T) {
	ApplyFunc(os.Getenv, func(key string) string {
		return `{
  "dbClusterIdentifier": "xxx",
  "password": "go_template_role456",
  "dbname": "go_template",
  "engine": "postgres",
  "port": 5432,
  "host": "localhost",
  "username": "go_template_role"
}`
	})
	db, err := postgres.Connect()
	if err != nil {
		assert.NotNil(t, db)
	}
	db.Close()
}
