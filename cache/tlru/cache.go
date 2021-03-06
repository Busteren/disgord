// tlru (time aware least recently lastUsed) has the same overwriting strategy as LRU, but adds a lifetime to objects
// on creation. Objects whose lifetime is outdated, are considered dead and can be removed.
package tlru

import (
	"sync"
	"time"

	"github.com/andersfylling/disgord/cache/interfaces"
	"github.com/andersfylling/snowflake/v2"
)

type Snowflake = snowflake.Snowflake

func NewCacheItem(content interface{}, lifetime time.Duration) *CacheItem {
	return &CacheItem{
		item:  content,
		death: time.Now().Add(lifetime).UnixNano(),
	}
}

type CacheItem struct {
	item interface{}

	// unix timestamp when item is considered outdated/dead.
	// update on creation/changes
	death int64

	// allows for least recently lastUsed monitoring
	lastUsed int64
}

func (i *CacheItem) Object() interface{} {
	return i.item
}

func (i *CacheItem) Set(v interface{}) {
	i.item = v
}

func (i *CacheItem) update() {
	now := time.Now()
	i.lastUsed = now.UnixNano()
}

func (i *CacheItem) dead(now time.Time) bool {
	return i.death <= now.UnixNano()
}

func NewCacheList(size uint, lifetime time.Duration) *CacheList {
	return &CacheList{
		items:    make(map[Snowflake]*CacheItem, size),
		limit:    size,
		lifetime: lifetime,
	}
}

type CacheList struct {
	sync.RWMutex
	items    map[Snowflake]*CacheItem
	limit    uint          // 0 == unlimited
	lifetime time.Duration // 0 == unlimited

	misses uint64 // opposite of cache hits
	hits   uint64
}

func (list *CacheList) size() uint {
	return uint(len(list.items))
}

func (list *CacheList) First() (item *CacheItem, key Snowflake) {
	for key, item = range list.items {
		return
	}

	return
}

// set adds a new item to the list or returns false if the item already exists
func (list *CacheList) Set(id Snowflake, newItemI interfaces.CacheableItem) {
	newItem := newItemI.(*CacheItem)
	newItem.update()
	if item, exists := list.items[id]; exists { // check if it points to a diff item
		if item.item != newItem.item || item != newItem {
			*item = *newItem
		}
		return
	} else {
		list.items[id] = newItem
	}

	if list.limit == 0 || list.size() <= list.limit {
		return
	}
	// if limit is reached, replace the content of the least recently lastUsed (lru)
	list.removeLRU(id)
}

func (list *CacheList) removeLRU(exception Snowflake) {
	lru, lruKey := list.First()
	for key, item := range list.items {
		if key != exception && item.lastUsed < lru.lastUsed {
			// TODO: create an lru map, for later?
			lru = item
			lruKey = key
		}
	}

	delete(list.items, lruKey)
}

func (list *CacheList) RefreshAfterDiscordUpdate(itemI interfaces.CacheableItem) {
	item := itemI.(*CacheItem)
	item.update()
	item.death = time.Now().Add(list.lifetime).UnixNano()
}

// get an item from the list.
func (list *CacheList) Get(id Snowflake) (ret interfaces.CacheableItem, exists bool) {
	var item *CacheItem
	if item, exists = list.items[id]; exists {
		ret = item
		item.update()
		list.hits++
	} else {
		list.misses++
	}
	return
}

func (list *CacheList) Delete(id Snowflake) {
	if _, exists := list.items[id]; exists {
		delete(list.items, id)
	}
}

func (list *CacheList) CreateCacheableItem(content interface{}) interfaces.CacheableItem {
	return NewCacheItem(content, list.lifetime)
}

// Efficiency ...
func (list *CacheList) Efficiency() float64 {
	return float64(list.hits) / float64(list.misses+list.hits)
}

var _ interfaces.CacheAlger = (*CacheList)(nil)
var _ interfaces.CacheableItem = (*CacheItem)(nil)
