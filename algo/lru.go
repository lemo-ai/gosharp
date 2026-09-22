package algo

// LRU is a simple generic least-recently-used cache.
type LRU[K comparable, V any] struct {
	cap   int
	items map[K]*lruNode[K, V]
	head  *lruNode[K, V]
	tail  *lruNode[K, V]
}

type lruNode[K comparable, V any] struct {
	key   K
	value V
	prev  *lruNode[K, V]
	next  *lruNode[K, V]
}

// NewLRU creates an LRU cache with the given capacity (must be > 0).
func NewLRU[K comparable, V any](capacity int) *LRU[K, V] {
	if capacity <= 0 {
		capacity = 1
	}
	head := &lruNode[K, V]{}
	tail := &lruNode[K, V]{}
	head.next = tail
	tail.prev = head
	return &LRU[K, V]{
		cap:   capacity,
		items: make(map[K]*lruNode[K, V], capacity),
		head:  head,
		tail:  tail,
	}
}

// Get returns the value for key and marks it as most recently used.
func (c *LRU[K, V]) Get(key K) (V, bool) {
	n, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	c.moveToFront(n)
	return n.value, true
}

// Put inserts or updates key and marks it as most recently used.
func (c *LRU[K, V]) Put(key K, value V) {
	if n, ok := c.items[key]; ok {
		n.value = value
		c.moveToFront(n)
		return
	}
	n := &lruNode[K, V]{key: key, value: value}
	c.items[key] = n
	c.addToFront(n)
	if len(c.items) > c.cap {
		oldest := c.tail.prev
		c.remove(oldest)
		delete(c.items, oldest.key)
	}
}

// Len returns the number of entries.
func (c *LRU[K, V]) Len() int {
	return len(c.items)
}

func (c *LRU[K, V]) addToFront(n *lruNode[K, V]) {
	n.prev = c.head
	n.next = c.head.next
	c.head.next.prev = n
	c.head.next = n
}

func (c *LRU[K, V]) remove(n *lruNode[K, V]) {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.prev = nil
	n.next = nil
}

func (c *LRU[K, V]) moveToFront(n *lruNode[K, V]) {
	c.remove(n)
	c.addToFront(n)
}
