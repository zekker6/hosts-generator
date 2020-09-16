package file_writer

import (
	"strings"
	"testing"
)

const (
	memoryStoreSize = 1000
)

var (
	removalList = []string{string(rune(0)), "\n"}
)

func getAdapter() HostsAdapter {
	adapt := NewMemoryHostsAdapter(0)
	return &adapt
}

func toString(d []byte) string {
	return strings.Trim(string(d), strings.Join(removalList, ""))
}

func TestWriter_Clear(t *testing.T) {
	adapter := getAdapter()
	w := NewWriter(adapter, "\n", "")

	t.Run("check clear deletes content", func(t *testing.T) {
		w.WriteToHosts("hahahahaha")

		w.Clear()

		data := make([]byte, memoryStoreSize*2)
		adapter.Read(data)

		returned := toString(data)
		if returned != "" {
			t.Errorf("Expected string to be empty, got: %s", returned)
		}
	})

}

func TestWriter_Write(t *testing.T) {
	t.Run("test write idempotence", func(t *testing.T) {
		adapter := getAdapter()
		w := NewWriter(adapter, "\n", "")

		w.WriteToHosts("hahahahaha")

		data1 := make([]byte, memoryStoreSize*2)
		adapter.Read(data1)

		w.WriteToHosts("hahahahaha")

		data2 := make([]byte, memoryStoreSize*2)
		adapter.Read(data2)
		d1s := toString(data1)
		d2s := toString(data2)

		if d1s != d2s {
			t.Errorf("Expected string to be equals, got: %s vs %s", d1s, d2s)
		}
	})

	t.Run("test appending data correctly", func(t *testing.T) {
		adapter := getAdapter()
		w := NewWriter(adapter, "\n", "")
		originalContent := "127.0.0.1 localhost"

		adapter.Write([]byte("127.0.0.1 localhost"))
		w.WriteToHosts("hahahahaha")

		data := make([]byte, memoryStoreSize*2)
		adapter.Read(data)

		if !strings.Contains(string(data), originalContent+"\n") {
			t.Errorf("did not found original content %s at %s", originalContent, data)
		}
	})

}
