package main

import (
	"sync"
	"time"
)

type Resource struct {
	url			string
	polling		bool
	lastPooled	int64
}

type Resources struct {
	data	[]*Resource
	lock	*sync.Mutex
}

func Pooler(res *Resources)  {
	for  {
		// Get the least recently polled Resource and mark it as being polled
		res.lock.Lock()
		var r *Resource

		for _, v := range res.data {
			if v.polling {
				continue
			}
			if r == nil || v.lastPooled < r.lastPooled {
				r = v
			}
		}
		if r != nil {
			r.polling = true
		}
		res.lock.Lock()

		if r == nil {
			continue
		}

		// Poll the URL

		// Update the Resource's polling and lastPolled
		res.lock.Lock()
		r.polling = false
		r.lastPooled = time.Now().UnixNano()
		res.lock.Unlock()
	}
	
}