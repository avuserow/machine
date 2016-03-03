package cert

import (
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
)

const MachinePrefix = "machine/v0"

func init() {
	etcd.Register()
	//consul.Register()
	//zookeeper.Register()
	//boltdb.Register()
}

type CertKvstore struct {
	store  store.Store
	prefix string
}

func NewCertKvstore(path string) (*CertKvstore, error) {
	fmt.Printf(`XXX NewCertKvstore("%s")`, path)
	var kvStore store.Store
	kvurl, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("Malformed store path: %s %s", path, err)
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
		return nil, fmt.Errorf("Unsupporetd KV store type: %s", kvurl.Scheme)
	}

	return &CertKvstore{
		store:  kvStore,
		prefix: kvurl.Path,
	}, nil
}

func (s CertKvstore) Write(filename string, data []byte) error {
	fmt.Printf(`XXX Write("%s", <data>)
`, filename)

	key := filepath.Join(s.prefix, MachinePrefix, filename)
	err := s.store.Put(key, data, nil)
	return err
}

func (s CertKvstore) Read(filename string) ([]byte, error) {
	fmt.Printf(`XXX Read("%s")
`, filename)
	key := filepath.Join(s.prefix, MachinePrefix, filename)
	kvpair, err := s.store.Get(key)
	if err != nil {
		return nil, err
	}
	return kvpair.Value, nil
}
