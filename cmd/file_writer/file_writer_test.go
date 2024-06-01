package file_writer

import (
	"fmt"
	"sync"
	"testing"
)

func TestConcurrentWrites(t *testing.T) {
	f1 := NewFileHostsAdapter("./test")
	//f2 := NewFileHostsAdapter("./test")
	var wg sync.WaitGroup
	contents := []string{
		"111111\n222222\n333333\n",
		"444444\n555555\n666666\n",
	}

	maxConcurrent := 10
	writeAttempts := 1000
	wg.Add(maxConcurrent)
	for i := 0; i < maxConcurrent; i++ {
		go func() {
			for i := 0; i < writeAttempts; i++ {
				for _, content := range contents {
					f1.Append([]byte(content))
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
	b := make([]byte, 1024)
	fmt.Println(f1.Read(b))
	//os.Remove("./test")
}
