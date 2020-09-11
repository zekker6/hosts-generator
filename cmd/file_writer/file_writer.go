package file_writer

import (
	"io/ioutil"
	"os"
)

type FileHostsAdapter struct {
	hostsLocation string
}

func NewFileHostsAdapter(hostsLocation string) FileHostsAdapter {
	return FileHostsAdapter{
		hostsLocation: hostsLocation,
	}
}

func (fa *FileHostsAdapter) Read(p []byte) (n int, err error) {
	data, err := ioutil.ReadFile(fa.hostsLocation)
	if err != nil {
		return 0, err
	}

	copy(p, data)
	return len(data), nil
}

func (fa *FileHostsAdapter) Write(p []byte) (n int, err error) {
	err = ioutil.WriteFile(fa.hostsLocation, p, 0600)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func (fa *FileHostsAdapter) Append(p []byte) error {
	f, err := os.OpenFile(fa.hostsLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(string(p)); err != nil {
		return err
	}

	return nil
}
