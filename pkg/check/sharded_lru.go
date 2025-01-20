// pkg/check/sharded_lru.go

package check

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

type SlabAllocator struct {
	pools map[int]*sync.Pool
	mu    sync.Mutex
}

func NewSlabAllocator() *SlabAllocator {
	return &SlabAllocator{
		pools: make(map[int]*sync.Pool),
	}
}

func (sa *SlabAllocator) Allocate(size int) interface{} {
	sa.mu.Lock()
	defer sa.mu.Unlock()

	if pool, exists := sa.pools[size]; exists {
		return pool.Get()
	}

	pool := &sync.Pool{
		New: func() interface{} {
			return make([]byte, size)
		},
	}
	sa.pools[size] = pool
	return pool.Get()
}

func (sa *SlabAllocator) Free(size int, obj interface{}) {
	sa.mu.Lock()
	defer sa.mu.Unlock()

	if pool, exists := sa.pools[size]; exists {
		pool.Put(obj)
	}
}

type entry struct {
	key   string
	value *KorcenResult
}

var entryPool = sync.Pool{
	New: func() interface{} {
		return &entry{}
	},
}

func getEntry(key string, value *KorcenResult) *entry {
	e := entryPool.Get().(*entry)
	e.key = key
	e.value = value
	return e
}

func putEntry(e *entry) {
	e.key = ""
	e.value = nil
	entryPool.Put(e)
}

type GoroutineManager struct {
	numWorkers int
	jobChan    chan func()
	stopChan   chan struct{}
	wg         sync.WaitGroup
}

func NewGoroutineManager(workerCount int) *GoroutineManager {
	manager := &GoroutineManager{
		numWorkers: workerCount,
		jobChan:    make(chan func(), 1000),
		stopChan:   make(chan struct{}),
	}

	manager.wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer manager.wg.Done()
			for {
				select {
				case job := <-manager.jobChan:
					job()
				case <-manager.stopChan:
					return
				}
			}
		}()
	}

	return manager
}

func (gm *GoroutineManager) Submit(job func()) {
	gm.jobChan <- job
}

func (gm *GoroutineManager) Shutdown() {
	close(gm.stopChan)
	gm.wg.Wait()
}

type UsageStats struct {
	Hits       uint64
	Misses     uint64
	Entries    uint64
	Requests   uint64
	AvgLatency int64
}

type TrackableLRUCache struct {
	lru          *LRUCache
	stats        UsageStats
	latencySum   int64
	requestCount uint64
}

func (t *TrackableLRUCache) Get(key string) (*KorcenResult, bool) {
	start := time.Now()
	value, ok := t.lru.Get(key)
	duration := time.Since(start).Nanoseconds()

	atomic.AddInt64(&t.latencySum, duration)
	atomic.AddUint64(&t.requestCount, 1)

	if ok {
		atomic.AddUint64(&t.stats.Hits, 1)
	} else {
		atomic.AddUint64(&t.stats.Misses, 1)
	}
	return value, ok
}

func (t *TrackableLRUCache) Set(key string, value *KorcenResult) error {
	start := time.Now()
	t.lru.Set(key, value)
	duration := time.Since(start).Nanoseconds()

	atomic.AddInt64(&t.latencySum, duration)
	atomic.AddUint64(&t.requestCount, 1)
	atomic.AddUint64(&t.stats.Entries, 1)

	return nil
}

func (t *TrackableLRUCache) GetStats() UsageStats {
	reqCount := atomic.LoadUint64(&t.requestCount)
	latency := atomic.LoadInt64(&t.latencySum)

	var avgLatency int64
	if reqCount > 0 {
		avgLatency = latency / int64(reqCount)
	}

	return UsageStats{
		Hits:       atomic.LoadUint64(&t.stats.Hits),
		Misses:     atomic.LoadUint64(&t.stats.Misses),
		Entries:    atomic.LoadUint64(&t.stats.Entries),
		Requests:   reqCount,
		AvgLatency: avgLatency,
	}
}

type LRUCache struct {
	capacity  int
	mu        sync.Mutex
	ll        *list.List
	cache     map[string]*list.Element
	allocator *SlabAllocator
}

func NewLRUCache(capacity int, allocator *SlabAllocator) *LRUCache {
	return &LRUCache{
		capacity:  capacity,
		ll:        list.New(),
		cache:     make(map[string]*list.Element, capacity),
		allocator: allocator,
	}
}

func (c *LRUCache) Get(key string) (*KorcenResult, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[key]; ok {
		c.ll.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}
	return nil, false
}

func (c *LRUCache) Set(key string, value *KorcenResult) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[key]; ok {
		c.ll.MoveToFront(elem)
		elem.Value.(*entry).value = value
		return
	}

	e := getEntry(key, value)
	elem := c.ll.PushFront(e)
	c.cache[key] = elem

	if c.ll.Len() > c.capacity {
		c.removeOldest()
	}
}

func (c *LRUCache) removeOldest() {
	elem := c.ll.Back()
	if elem != nil {
		c.ll.Remove(elem)
		ent, ok := elem.Value.(*entry)
		if !ok {
			return
		}
		delete(c.cache, ent.key)
		freeKorcenResult(ent.value)
		putEntry(ent)
	}
}

type ShardedLRUCache struct {
	shards     []*TrackableLRUCache
	shardCount int
	manager    *GoroutineManager
	allocator  *SlabAllocator
}

func NewShardedLRUCache(shardCount int, perShardCapacity int) *ShardedLRUCache {
	if shardCount <= 0 || (shardCount&(shardCount-1)) != 0 {
		panic("shardCount must be a power of two")
	}

	allocator := NewSlabAllocator()
	shards := make([]*TrackableLRUCache, shardCount)
	for i := 0; i < shardCount; i++ {
		lru := NewLRUCache(perShardCapacity, allocator)
		trackedLRU := &TrackableLRUCache{
			lru: lru,
		}
		shards[i] = trackedLRU
	}

	manager := NewGoroutineManager(shardCount)

	return &ShardedLRUCache{
		shards:     shards,
		shardCount: shardCount,
		manager:    manager,
		allocator:  allocator,
	}
}

func (c *ShardedLRUCache) getShard(key string) *TrackableLRUCache {
	hash := Murmur3Hash(key)
	return c.shards[hash&(uint32(c.shardCount)-1)]
}

func (c *ShardedLRUCache) Get(key string) (*KorcenResult, bool) {
	shard := c.getShard(key)
	return shard.Get(key)
}

func (c *ShardedLRUCache) Set(key string, value *KorcenResult) error {
	shard := c.getShard(key)
	done := make(chan error, 1)

	c.manager.Submit(func() {
		err := shard.Set(key, value)
		done <- err
	})

	err := <-done
	if err != nil {
		return err
	}
	return nil
}

func (c *ShardedLRUCache) GetStats() UsageStats {
	var totalStats UsageStats
	for _, shard := range c.shards {
		stats := shard.GetStats()
		totalStats.Hits += stats.Hits
		totalStats.Misses += stats.Misses
		totalStats.Entries += stats.Entries
		totalStats.Requests += stats.Requests
		totalStats.AvgLatency += stats.AvgLatency
	}

	if totalStats.Requests > 0 {
		totalStats.AvgLatency /= int64(c.shardCount)
	}
	return totalStats
}

func (c *ShardedLRUCache) Stop() {
	c.manager.Shutdown()
}
