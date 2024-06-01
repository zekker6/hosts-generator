package cmd

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"hosts-generator/cmd/file_writer"
	"hosts-generator/cmd/parsers"
)

type fakeClient struct {
	addresses []string

	m sync.Mutex

	itemsCount   int
	refreshEvery time.Duration
	changeRate   float64
}

func newFakeClient(itemsCount int, refreshEvery time.Duration, changeRate float64) *fakeClient {
	c := &fakeClient{
		itemsCount:   itemsCount,
		refreshEvery: refreshEvery,
		changeRate:   changeRate,
		m:            sync.Mutex{},
	}
	go c.run()

	return c
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
func (c *fakeClient) run() {
	randAddr := func(i int) string {
		return fmt.Sprintf("%s-addr-%d", randString(10), i)
	}

	c.m.Lock()
	if c.addresses == nil {
		c.addresses = make([]string, c.itemsCount)
		for i := 0; i < c.itemsCount; i++ {
			c.addresses[i] = randAddr(i)
		}
	}
	c.m.Unlock()

	toUpdate := int(float64(c.itemsCount) * c.changeRate)
	for {
		time.Sleep(c.refreshEvery)
		c.m.Lock()
		for i := 0; i < toUpdate; i++ {
			idx := rand.Intn(c.itemsCount)
			c.addresses[idx] = randAddr(idx)
		}
		c.m.Unlock()
	}
}

func (c *fakeClient) Get() ([]string, error) {
	c.m.Lock()
	defer c.m.Unlock()

	r := make([]string, len(c.addresses))
	copy(r, c.addresses)
	return r, nil
}

func TestAppRun(t *testing.T) {
	c := newFakeClient(100, 10*time.Millisecond, 0.1)
	c2 := newFakeClient(100, 10*time.Millisecond, 0.1)

	writer := file_writer.NewMemoryHostsAdapter(10000)
	newline := "\n"

	initialContent := `
# Standard host addresses
127.0.0.1        localhost
::1              localhost
ff02::1          ip6-allnodes
ff02::2          ip6-allrouters
`
	writer.Write([]byte(initialContent))

	adapter := file_writer.NewWriter(&writer, newline, "")

	app := NewApp([]parsers.Parser{c, c2}, adapter, newline, "127.0.0.1", 1*time.Millisecond, true, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := app.Run(ctx)
	if err != nil {
		t.Errorf("Expected no error, got: %+v", err)
	}

	out := make([]byte, 10000)
	_, _ = writer.Read(out)
	fmt.Println(string(out))
}
