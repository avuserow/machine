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
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnerror"
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

func (s Kvstore) loadConfig(h *host.Host, data []byte) error {
	// Remember the machine name so we don't have to pass it through each
	// struct in the migration.
	name := h.Name

	migratedHost, migrationPerformed, err := host.MigrateHost(h, data)
	if err != nil {
		return fmt.Errorf("Error getting migrated host: %s", err)
	}

	*h = *migratedHost

	h.Name = name

	// If we end up performing a migration, we should save afterwards so we don't have to do it again on subsequent invocations.
	log.Infof("AK: migration performed: %s", migrationPerformed)
	if migrationPerformed {
		// XXX TODO do we want to save?

		//		if err := s.saveToFile(data, filepath.Join(s.GetMachinesDir(), h.Name, "config.json.bak")); err != nil {
		//			return fmt.Errorf("Error attempting to save backup after migration: %s", err)
		//		}
		//
		//		if err := s.Save(h); err != nil {
		//			return fmt.Errorf("Error saving config after migration was performed: %s", err)
		//		}
	}

	return nil
}

func (s Kvstore) Load(name string) (*host.Host, error) {
	hostPath := filepath.Join(s.prefix, MachinePrefix, "machines", name)

	if exists, err := s.Exists(name); err != nil || exists != true {
		return nil, mcnerror.ErrHostDoesNotExist{
			Name: name,
		}
	}

	kvPair, err := s.store.Get(hostPath)
	if err != nil {
		return nil, err
	}

	fmt.Println("Load: ", kvPair.Key)

	host := &host.Host{
		Name: name,
	}

	if err := s.loadConfig(host, kvPair.Value); err != nil {
		return nil, err
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
	hostPath := filepath.Join(s.prefix, MachinePrefix, "machines", name)

	err := s.store.Delete(hostPath)
	return err
}

func (s Kvstore) GetMachinesDir() string {
	return filepath.Join(s.prefix, MachinePrefix, "machines")
}
