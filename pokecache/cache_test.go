package pokecache

import (
	"fmt"
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

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}
