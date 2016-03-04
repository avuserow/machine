package cert

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
	"github.com/docker/machine/libmachine/auth"
)

const MachinePrefix = "machine/v0"

func init() {
	etcd.Register()
	//consul.Register()
	//zookeeper.Register()
	//boltdb.Register()
}

type CertKvstore struct {
	authOptions *auth.Options
	store       store.Store
	prefix      string
}

func NewCertKvstore(authOptions *auth.Options) (*CertKvstore, error) {
	fmt.Printf(`XXX NewCertKvstore("%s")
`, authOptions.CertDir)
	var kvStore store.Store
	kvurl, err := url.Parse(authOptions.CertDir)
	if err != nil {
		return nil, fmt.Errorf("Malformed store path: %s %s", authOptions.CertDir, err)
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

	// XXX This feels super messy - there's got to be a better way.
	//prefix := fmt.Sprintf("%s:/%s", kvurl.Scheme, kvurl.Host)
	//authOptions.CertDir = strings.TrimPrefix(authOptions.CertDir, prefix)
	//authOptions.CaCertPath = strings.TrimPrefix(authOptions.CaCertPath, prefix)
	//authOptions.CaPrivateKeyPath = strings.TrimPrefix(authOptions.CaPrivateKeyPath, prefix)
	//authOptions.CaCertRemotePath = strings.TrimPrefix(authOptions.CaCertRemotePath, prefix)
	//authOptions.ServerCertPath = strings.TrimPrefix(authOptions.ServerCertPath, prefix)
	//authOptions.ServerKeyPath = strings.TrimPrefix(authOptions.ServerKeyPath, prefix)
	//authOptions.ClientKeyPath = strings.TrimPrefix(authOptions.ClientKeyPath, prefix)
	//authOptions.ServerCertRemotePath = strings.TrimPrefix(authOptions.ServerCertRemotePath, prefix)
	//authOptions.ServerKeyRemotePath = strings.TrimPrefix(authOptions.ServerKeyRemotePath, prefix)
	//authOptions.ClientCertPath = strings.TrimPrefix(authOptions.ClientCertPath, prefix)
	//authOptions.StorePath = strings.TrimPrefix(authOptions.StorePath, prefix)
	fmt.Printf("XXX CertDir: %s\n", authOptions.CertDir)
	fmt.Printf("XXX CaCertPath: %s\n", authOptions.CaCertPath)

	return &CertKvstore{
		store:       kvStore,
		prefix:      kvurl.Path,
		authOptions: authOptions,
	}, nil
}

func (s CertKvstore) Write(filename string, data []byte, flag int, perm os.FileMode) error {

	key := filepath.Join("/", MachinePrefix, s.prefix, filename)
	fmt.Printf("XXX KV Write -> %s\n", key)
	err := s.store.Put(key, data, nil)
	fmt.Printf("XXX err: %s\n", err)
	return err
}

func (s CertKvstore) Read(filename string) ([]byte, error) {
	key := filepath.Join("/", MachinePrefix, s.prefix, filename)
	fmt.Printf("XXX KV Read -> %s\n", key)
	kvpair, err := s.store.Get(key)
	if err != nil {
		return nil, err
	}
	fmt.Printf("XXX err: %s\n", err)
	return kvpair.Value, nil
}
func (s CertKvstore) Exists(filename string) bool {
	key := filepath.Join("/", MachinePrefix, s.prefix, filename)
	fmt.Printf("XXX KV Exists -> %s\n", key)
	exists, err := s.store.Exists(key)
	if err != nil {
		// TODO log a better message on other errors
		fmt.Printf("KV lookup failure on %s: %s\n", filename, err)
		return false
	}
	fmt.Printf("XXX err: %s\n", err)
	fmt.Printf("XXX exists: %v\n", exists)
	return exists
}
