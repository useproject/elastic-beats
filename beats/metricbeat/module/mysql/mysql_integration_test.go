// +build integration

package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/useproject/elastic-beats/beats/libbeat/tests/compose"
	_ "github.com/useproject/elastic-beats/beats/metricbeat/mb/testing"
)

func TestNewDB(t *testing.T) {
	compose.EnsureUp(t, "mysql")

	db, err := NewDB(GetMySQLEnvDSN())
	assert.NoError(t, err)

	err = db.Ping()
	assert.NoError(t, err)
}
