package servicediscovery

// KVPair contains a key and a string.
type KVPair struct {
	Key   string
	Value string
}

type ServiceDiscovery interface {
	GetServices() []*KVPair
	WatchService() chan []*KVPair
	RemoveWatcher(ch chan []*KVPair)
	Clone(servicePath string) ServiceDiscovery
	Close()
}
