package pokecache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := NewCache(2 * time.Second)
	key := "testKey"
	val := []byte("testData")

	c.Add(key, val)

	cachedVal, ok := c.Get(key)
	if !ok || string(cachedVal) != string(val) {
		t.Errorf("Expected %s, got %s", val, cachedVal)
	}

	time.Sleep(3 * time.Second)
	_, ok = c.Get(key)
	if ok {
		t.Error("Expected cache to expire, but entry still exists")
	}
}
