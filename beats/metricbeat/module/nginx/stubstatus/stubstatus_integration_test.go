// +build integration

package stubstatus

import (
	"testing"

	"github.com/useproject/elastic-beats/beats/libbeat/tests/compose"
	mbtest "github.com/useproject/elastic-beats/beats/metricbeat/mb/testing"
	"github.com/useproject/elastic-beats/beats/metricbeat/module/nginx"

	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	compose.EnsureUp(t, "nginx")

	f := mbtest.NewEventFetcher(t, getConfig())
	event, err := f.Fetch()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	t.Logf("%s/%s event: %+v", f.Module().Name(), f.Name(), event)

	// Check number of fields.
	assert.Equal(t, 10, len(event))
}

func TestData(t *testing.T) {
	f := mbtest.NewEventFetcher(t, getConfig())

	err := mbtest.WriteEvent(f, t)
	if err != nil {
		t.Fatal("write", err)
	}
}

func getConfig() map[string]interface{} {
	return map[string]interface{}{
		"module":     "nginx",
		"metricsets": []string{"stubstatus"},
		"hosts":      []string{nginx.GetNginxEnvHost()},
	}
}
