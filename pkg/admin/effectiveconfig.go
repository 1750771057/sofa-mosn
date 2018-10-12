package admin

import (
	"sync"

	"github.com/alipay/sofa-mosn/pkg/api/v2"
)

type effectiveConfig struct {
	MOSNConfig interface{}            `json:"mosn_config,omitempty"`
	Listener   map[string]v2.Listener `json:"listener,omitempty"`
	Cluster    map[string]v2.Cluster  `json:"cluster,omitempty"`
}

var (
	conf = effectiveConfig{
		Listener: make(map[string]v2.Listener),
		Cluster:  make(map[string]v2.Cluster),
	}
	mutex = new(sync.RWMutex)
)

func Reset() {
	conf.MOSNConfig = nil
	conf.Listener = make(map[string]v2.Listener)
	conf.Cluster = make(map[string]v2.Cluster)
}

func SetMOSNConfig(msonConfig interface{}) {
	conf.MOSNConfig = msonConfig
}

// SetListenerConfig
// Set listener config when AddOrUpdateListener
func SetListenerConfig(listenerName string, listenerConfig v2.Listener) {
	if config, ok := conf.Listener[listenerName]; ok {
		config.ListenerConfig.FilterChains = listenerConfig.ListenerConfig.FilterChains
		config.ListenerConfig.StreamFilters = listenerConfig.ListenerConfig.StreamFilters
		conf.Listener[listenerName] = config
	} else {
		conf.Listener[listenerName] = listenerConfig
	}
}

func SetClusterConfig(clusterName string, cluster v2.Cluster) {
	conf.Cluster[clusterName] = cluster
}

func SetHosts(clusterName string, hostConfigs []v2.Host) {
	if cluster, ok := conf.Cluster[clusterName]; ok {
		cluster.Hosts = hostConfigs
	}
}

// Dump
// Dump all config
func Dump() ([]byte, error) {
	mutex.Lock()
	defer mutex.Unlock()
	return json.Marshal(conf)
}
