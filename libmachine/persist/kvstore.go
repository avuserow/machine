package persist

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
	"github.com/docker/machine/libmachine/host"
)

const MachinePrefix = "machine/v0"

func init() {
	etcd.Register()
	//consul.Register()
	//zookeeper.Register()
	//boltdb.Register()
}

type Kvstore struct {
	store  store.Store
	prefix string
}

func NewKvstore(path string, certsDir string) *Kvstore {
	fmt.Printf(`XXX NewKvstore("%s", "%s")`, path, certsDir)
	var kvStore store.Store
	kvurl, err := url.Parse(path)
	if err != nil {
		panic(fmt.Sprintf("Malformed store path: %s %s", path, err))
	}
	switch kvurl.Scheme {
	case "etcd":
		// TODO - figure out how to get TLS support in here...
		kvStore, err = libkv.NewStore(
			store.ETCD,
			[]string{kvurl.Host},
			&store.Config{
				ConnectionTimeout: 10 * time.Second,
			},
		)
		// TODO other KV store types
	default:
		panic(fmt.Sprintf("Unsupporetd KV store type: %s", kvurl.Scheme))
	}

	return &Kvstore{
		store:  kvStore,
		prefix: kvurl.Path,
	}
}

func (s Kvstore) Save(host *host.Host) error {
	data, err := json.Marshal(host)
	if err != nil {
		return err
	}

	hostPath := filepath.Join(s.prefix, MachinePrefix, "machines", host.Name)
	err = s.store.Put(hostPath, data, nil)
	return err
}

func (s Kvstore) Exists(name string) (bool, error) {
	hostPath := filepath.Join(s.prefix, MachinePrefix, "machines", name)
	return s.store.Exists(hostPath)
}

func (s Kvstore) Load(name string) (*host.Host, error) {
	hostPath := filepath.Join(s.prefix, MachinePrefix, "machines", name)

	kvPair, err := s.store.Get(hostPath)
	if err != nil {
		return nil, err
	}

	fmt.Println("Load: ", kvPair.Key)

	host := &host.Host{
		Name: name,
	}

	return host, nil
}

func (s Kvstore) List() ([]string, error) {
	machineDir := filepath.Join(s.prefix, MachinePrefix, "machines")
	kvList, err := s.store.List(machineDir)
	if err == store.ErrKeyNotFound {
		// No machines set up
		return []string{}, nil
	} else if err != nil {
		return nil, err
	}

	hostNames := []string{}

	for _, kvPair := range kvList {
		hostNames = append(hostNames, kvPair.Key)
	}

	return hostNames, nil
}

func (s Kvstore) Remove(name string) error {
	fmt.Println("XXX: Remove")
	return nil
}

func (s Kvstore) GetMachinesDir() string {
	return filepath.Join(s.prefix, MachinePrefix, "machines")
}
