package cache

type Cache interface {
	Get(key int64) (value interface{}, exists bool)
	Remove(key int64)
	Add(key int64, value interface{})
}

// New returns a new cache instance
// which MUST be safe for concurrent use
func New(maxEntries int) Cache {
	switch maxEntries {
	case 0:
		return newNop()
	default:
		return newLRUCache(maxEntries)
	}
}
