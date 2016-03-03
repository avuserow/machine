package cert

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type CertFilestore struct {
	Path string
}

func NewCertFilestore(path string) (*CertFilestore, error) {
	fmt.Printf(`XXX NewCertFilestore("%s")
`, path)

	return &CertFilestore{
		Path: path,
	}, nil
}

func (s CertFilestore) Write(filename string, data []byte) error {
	fmt.Printf(`XXX Write("%s", <data>)
`, filename)

	// TODO - audit/verify this impl if we keep it
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return ioutil.WriteFile(filename, data, 0600)
	}

	tmpfi, err := ioutil.TempFile(filepath.Dir(filename), ".tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfi.Name())

	if err = ioutil.WriteFile(tmpfi.Name(), data, 0600); err != nil {
		return err
	}

	if err = tmpfi.Close(); err != nil {
		return err
	}

	if err = os.Remove(filename); err != nil {
		return err
	}

	if err = os.Rename(tmpfi.Name(), filename); err != nil {
		return err
	}
	return nil
}

func (s CertFilestore) Read(filename string) ([]byte, error) {
	fmt.Printf(`XXX Read("%s")
`, filename)
	return ioutil.ReadFile(filename)
}
