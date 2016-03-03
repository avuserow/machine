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
	store store.Store
}

func NewKvstore(path string, certsDir string) *Kvstore {
	fmt.Printf(`XXX NewKvstore("%s", "%s")`, path, certsDir)
	kvurl, err := url.Parse(path)
	var kvStore store.Store
	if err != nil {
		switch kvurl.Scheme {
		case "etcd":
			kvStore, err = libkv.NewStore(
				store.ETCD,
				[]string{kvurl.Host},
				&store.Config{
					ConnectionTimeout: 10 * time.Second,
				},
			)
			// TODO other KV store types
		}
	}

	return &Kvstore{store: kvStore}
}

func (s Kvstore) Save(host *host.Host) error {
	data, err := json.Marshal(host)
	if err != nil {
		return err
	}

	hostPath := filepath.Join(MachinePrefix, "machines", host.Name)
	err = s.store.Put(hostPath, data, nil)
	return err
}

func (s Kvstore) Exists(name string) (bool, error) {
	hostPath := filepath.Join(MachinePrefix, "machines", name)
	return s.store.Exists(hostPath)
}

func (s Kvstore) Load(name string) (*host.Host, error) {
	hostPath := filepath.Join(MachinePrefix, "machines", name)

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
	machineDir := filepath.Join(MachinePrefix, "machines")
	kvList, err := s.store.List(machineDir)
	if err != nil {
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
	return filepath.Join(MachinePrefix, "machines")
}
