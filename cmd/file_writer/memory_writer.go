package file_writer

type MemoryHostsAdapter struct {
	data []byte
}

func NewMemoryHostsAdapter(size int) MemoryHostsAdapter {
	return MemoryHostsAdapter{
		data: make([]byte, 0, size),
	}
}

func (ma *MemoryHostsAdapter) Read(p []byte) (n int, err error) {
	n = copy(p, ma.data)
	return
}

func (ma *MemoryHostsAdapter) Write(p []byte) (n int, err error) {
	ma.data = make([]byte, len(p))
	copy(ma.data, p)

	return len(ma.data), nil
}

func (ma *MemoryHostsAdapter) Append(p []byte) error {
	ma.data = append(ma.data, p...)
	return nil
}
