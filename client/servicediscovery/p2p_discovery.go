package servicediscovery

type P2PDiscovery struct {
	server   string
	metadata string
}

// NewP2PDiscovery returns a new P2PDiscovery.
func NewP2PDiscovery(server, metadata string) ServiceDiscovery {
	return &P2PDiscovery{server: server, metadata: metadata}
}

// Clone clones this ServiceDiscovery with new servicePath.
func (d P2PDiscovery) Clone(servicePath string) ServiceDiscovery {
	return &d
}

// GetServices returns the static server
func (d P2PDiscovery) GetServices() []*KVPair {
	return []*KVPair{&KVPair{Key: d.server, Value: d.metadata}}
}

// WatchService returns a nil chan.
func (d P2PDiscovery) WatchService() chan []*KVPair {
	return nil
}

func (d *P2PDiscovery) RemoveWatcher(ch chan []*KVPair) {}

func (d *P2PDiscovery) Close() {

}
