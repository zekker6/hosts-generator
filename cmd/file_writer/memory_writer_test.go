package file_writer

import "testing"

func TestMemoryHostsAdapter_Initial(t *testing.T) {
	adapter := getAdapter()
	t.Run("validate default is empty", func(t *testing.T) {
		out := make([]byte, memoryStoreSize)
		adapter.Read(out)

		if toString(out) != "" {
			t.Errorf("Expteced buffer to be empty, got: %+v", toString(out))
		}
	})
}

func TestMemoryHostsAdapter_Adapter(t *testing.T) {
	adapter := getAdapter()
	t.Run("test appending data", func(t *testing.T) {
		content := "some_data"

		adapter.Append([]byte(content))

		out := make([]byte, memoryStoreSize)
		adapter.Read(out)

		if content != toString(out) {
			t.Errorf("Expteced buffer to contain test data: %s, got: %s", content, toString(out))
		}

		adapter.Append([]byte(content))
		adapter.Read(out)

		if content+content != toString(out) {
			t.Errorf("Expteced buffer to contain test data: %s, got: %s", content+content, toString(out))
		}
	})
}

func TestMemoryHostsAdapter_Write(t *testing.T) {
	adapter := getAdapter()
	t.Run("test appending data", func(t *testing.T) {
		content := "some_data"

		adapter.Write([]byte(content))

		out := make([]byte, memoryStoreSize)
		adapter.Read(out)

		if content != toString(out) {
			t.Errorf("Expteced buffer to contain test data: %s, got: %s", content, toString(out))
		}
	})
}
