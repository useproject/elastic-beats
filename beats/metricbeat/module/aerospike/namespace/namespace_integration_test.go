// +build integration

package namespace

import (
	"testing"

	"github.com/useproject/elastic-beats/beats/libbeat/tests/compose"
	mbtest "github.com/useproject/elastic-beats/beats/metricbeat/mb/testing"
	"github.com/useproject/elastic-beats/beats/metricbeat/module/aerospike"
)

func TestData(t *testing.T) {
	compose.EnsureUp(t, "aerospike")

	f := mbtest.NewEventsFetcher(t, getConfig())
	err := mbtest.WriteEvents(f, t)
	if err != nil {
		t.Fatal("write", err)
	}
}

func getConfig() map[string]interface{} {
	return map[string]interface{}{
		"module":     "aerospike",
		"metricsets": []string{"namespace"},
		"hosts":      []string{aerospike.GetAerospikeEnvHost() + ":" + aerospike.GetAerospikeEnvPort()},
	}
}
