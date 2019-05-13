package system

import (
	"github.com/useproject/elastic-beats/beats/libbeat/logp"
	"github.com/useproject/elastic-beats/beats/metricbeat/helper"
)

func initModule() {
	if err := helper.CheckAndEnableSeDebugPrivilege(); err != nil {
		logp.Warn("%v", err)
	}
}
