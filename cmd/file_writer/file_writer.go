package file_writer

import (
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
	data, err := os.ReadFile(fa.hostsLocation)
	if err != nil {
		return 0, err
	}

	copy(p, data)
	return len(data), nil
}

func (fa *FileHostsAdapter) Write(p []byte) (n int, err error) {
	err = os.WriteFile(fa.hostsLocation, p, 0600)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func (fa *FileHostsAdapter) Append(p []byte) error {
	data, err := os.ReadFile(fa.hostsLocation)
	if err != nil {
		return err
	}

	data = append(data, p...)
	return os.WriteFile(fa.hostsLocation, data, 0600)
}
