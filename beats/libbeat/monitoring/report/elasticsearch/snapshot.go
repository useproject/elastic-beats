package elasticsearch

import (
	"github.com/useproject/elastic-beats/beats/libbeat/common"
	"github.com/useproject/elastic-beats/beats/libbeat/monitoring"
)

func makeSnapshot(R *monitoring.Registry) common.MapStr {
	mode := monitoring.Full
	return common.MapStr(monitoring.CollectStructSnapshot(R, mode, false))
}
