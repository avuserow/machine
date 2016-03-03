package persist

import (
	"fmt"
	"net/url"
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
	"github.com/docker/machine/libmachine/host"
)

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
	fmt.Println("XXX: save NYI", host)
	return nil
}

func (s Kvstore) Exists(name string) (bool, error) {
	fmt.Println("XXX: exists NYI", name)

	return false, nil
}

func (s Kvstore) Load(name string) (*host.Host, error) {
	fmt.Println("XXX: load NYI", name)
	return nil, nil
}
