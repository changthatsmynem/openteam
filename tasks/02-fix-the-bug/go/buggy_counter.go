package counter

import (
	"sync"
	"time"
)

// NextID generates a unique ID by incrementing a current number, so if there are multiple 
// goroutines calling NextID concurrently, "race condition" can occur, 
// becasue multiple goroutines will read and write to the shared variable "current" at the same time, leading to duplicate IDs.
// To fix this, I use a "mutex" to ensure that only one goroutine can access the variable at a time

var current int64
var mu sync.Mutex

func NextID() int64 {
	mu.Lock()
	defer mu.Unlock()
	id := current
	time.Sleep(0)
	current++
	return id
}
