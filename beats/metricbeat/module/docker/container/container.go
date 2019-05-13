package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"github.com/useproject/elastic-beats/beats/libbeat/common"
	"github.com/useproject/elastic-beats/beats/metricbeat/mb"
	"github.com/useproject/elastic-beats/beats/metricbeat/module/docker"
)

func init() {
	mb.Registry.MustAddMetricSet("docker", "container", New,
		mb.WithHostParser(docker.HostParser),
		mb.DefaultMetricSet(),
	)
}

type MetricSet struct {
	mb.BaseMetricSet
	dockerClient *client.Client
	dedot        bool
}

// New creates a new instance of the docker container MetricSet.
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	config := docker.DefaultConfig()
	if err := base.Module().UnpackConfig(&config); err != nil {
		return nil, err
	}

	client, err := docker.NewDockerClient(base.HostData().URI, config)
	if err != nil {
		return nil, err
	}

	return &MetricSet{
		BaseMetricSet: base,
		dockerClient:  client,
		dedot:         config.DeDot,
	}, nil
}

// Fetch returns a list of all containers as events.
// This is based on https://docs.docker.com/engine/reference/api/docker_remote_api_v1.24/#/list-containers.
func (m *MetricSet) Fetch() ([]common.MapStr, error) {
	// Fetch a list of all containers.
	containers, err := m.dockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	return eventsMapping(containers, m.dedot), nil
}
