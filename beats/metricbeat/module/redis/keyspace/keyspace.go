package keyspace

import (
	"time"

	"github.com/useproject/elastic-beats/beats/libbeat/common"
	"github.com/useproject/elastic-beats/beats/libbeat/logp"
	"github.com/useproject/elastic-beats/beats/metricbeat/mb"
	"github.com/useproject/elastic-beats/beats/metricbeat/mb/parse"
	"github.com/useproject/elastic-beats/beats/metricbeat/module/redis"

	rd "github.com/garyburd/redigo/redis"
)

var (
	debugf = logp.MakeDebug("redis-keyspace")
)

func init() {
	mb.Registry.MustAddMetricSet("redis", "keyspace", New,
		mb.WithHostParser(parse.PassThruHostParser),
		mb.DefaultMetricSet(),
	)
}

// MetricSet for fetching Redis server information and statistics.
type MetricSet struct {
	mb.BaseMetricSet
	pool *rd.Pool
}

// New creates new instance of MetricSet
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	// Unpack additional configuration options.
	config := struct {
		IdleTimeout time.Duration `config:"idle_timeout"`
		Network     string        `config:"network"`
		MaxConn     int           `config:"maxconn" validate:"min=1"`
		Password    string        `config:"password"`
	}{
		Network:  "tcp",
		MaxConn:  10,
		Password: "",
	}
	err := base.Module().UnpackConfig(&config)
	if err != nil {
		return nil, err
	}

	return &MetricSet{
		BaseMetricSet: base,
		pool: redis.CreatePool(base.Host(), config.Password, config.Network,
			config.MaxConn, config.IdleTimeout, base.Module().Config().Timeout),
	}, nil
}

// Fetch fetches metrics from Redis by issuing the INFO command.
func (m *MetricSet) Fetch() ([]common.MapStr, error) {
	// Fetch default INFO.
	info, err := redis.FetchRedisInfo("keyspace", m.pool.Get())
	if err != nil {
		return nil, err
	}

	debugf("Redis INFO from %s: %+v", m.Host(), info)
	return eventsMapping(info), nil
}
