package system

import (
	"encoding/json"
	"fmt"

	"github.com/useproject/elastic-beats/beats/libbeat/common"
	"github.com/useproject/elastic-beats/beats/metricbeat/mb"
	"github.com/useproject/elastic-beats/beats/metricbeat/module/kubernetes"
)

func eventMapping(content []byte) ([]common.MapStr, error) {
	events := []common.MapStr{}

	var summary kubernetes.Summary
	err := json.Unmarshal(content, &summary)
	if err != nil {
		return nil, fmt.Errorf("Cannot unmarshal json response: %s", err)
	}

	node := summary.Node

	for _, syscontainer := range node.SystemContainers {
		containerEvent := common.MapStr{
			mb.ModuleDataKey: common.MapStr{
				"node": common.MapStr{
					"name": node.NodeName,
				},
			},
			"container":  syscontainer.Name,
			"start_time": syscontainer.StartTime,
			"cpu": common.MapStr{
				"usage": common.MapStr{
					"nanocores": syscontainer.CPU.UsageNanoCores,
					"core": common.MapStr{
						"ns": syscontainer.CPU.UsageCoreNanoSeconds,
					},
				},
			},
			"memory": common.MapStr{
				"usage": common.MapStr{
					"bytes": syscontainer.Memory.UsageBytes,
				},
				"workingset": common.MapStr{
					"bytes": syscontainer.Memory.WorkingSetBytes,
				},
				"rss": common.MapStr{
					"bytes": syscontainer.Memory.RssBytes,
				},
				"pagefaults":      syscontainer.Memory.PageFaults,
				"majorpagefaults": syscontainer.Memory.MajorPageFaults,
			},
		}
		events = append(events, containerEvent)
	}

	return events, nil
}
