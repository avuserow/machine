package cert

import (
	"net/url"

	"github.com/docker/machine/libmachine/auth"
)

type CertStore interface {
	// TODO - flesh this out once we know what we need

	// This is probably too low level
	Write(filename string, data []byte) error

	// This is probably too low level
	Read(filename string) ([]byte, error)
}

func NewCertStore(authOptions *auth.Options) (CertStore, error) {
	// Determine which type of store to generate
	storeURL, err := url.Parse(authOptions.CertDir)
	if err == nil {
		// The scheme will be blank on unix paths, might be a drive letter (single char)
		// or a multi-character scheme that libkv will hopefully handle
		if len(storeURL.Scheme) > 1 {
			return NewCertKvstore(authOptions)
		}
	}
	return NewCertFilestore(authOptions)
}
