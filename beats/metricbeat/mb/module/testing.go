package module

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/useproject/elastic-beats/beats/libbeat/beat"
	"github.com/useproject/elastic-beats/beats/libbeat/common"
	"github.com/useproject/elastic-beats/beats/libbeat/testing"
)

// receiveOneEvent receives one event from the events channel then closes the
// returned done channel. If no events are received it will close the returned
// done channel after the timeout period elapses.
func receiveOneEvent(d testing.Driver, events <-chan beat.Event, timeout time.Duration) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		select {
		case <-time.Tick(timeout):
			d.Error("error", errors.New("timeout waiting for an event"))
		case event, ok := <-events:
			if !ok {
				return
			}

			// At this point in the pipeline the error has been converted to a
			// string and written to error.message.
			if v, err := event.Fields.GetValue("error.message"); err == nil {
				if errMsg, ok := v.(string); ok {
					d.Error("error", errors.New(errMsg))
					return
				}
			}

			outputJSON(d, &event)
		}
	}()

	return done
}

func outputJSON(d testing.Driver, event *beat.Event) {
	out := event.Fields.Clone()
	out.Put("@timestamp", common.Time(event.Timestamp))
	jsonData, err := json.MarshalIndent(out, "", " ")
	if err != nil {
		d.Error("convert error", err)
		return
	}

	d.Result(string(jsonData))
}
