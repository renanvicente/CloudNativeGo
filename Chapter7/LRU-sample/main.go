package main

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
)

var cache *lru.Cache

func init() {
	cache, _ = lru.NewWithEvict(2,
		func(key interface{}, value interface{}) {
			fmt.Printf("Evicted: key=%v value=%v\n", key, value)
		})
}

//func ThreadSafeWrite(key, value string)  {
//	cache.Lock()                                    // Establish write lock
//	cache.data[key] = value
//	cache.Unlock()                                  // Release write lock
//}
func main() {
	cache.Add(1, "a")		// adds 1
	cache.Add(2, "b")		// adds 2; cache is now at capacity

	fmt.Println(cache.Get(1))		// "a true"; 1 now most recently used

	cache.Add(3,"c")			// adds 3, evicts key 2

	fmt.Println(cache.Get(2))


}