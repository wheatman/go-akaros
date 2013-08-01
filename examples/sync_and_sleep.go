package main

import (
	"fmt"
	"time"
	"sync"
	"sync/atomic"
)

var m sync.Mutex
var w sync.WaitGroup
const NUM_GS = 0x100
const NUM_LOOPS = 0x200
var done int32

func test(id int) {
	for i := 0; i < NUM_LOOPS; i++ {
		time.Sleep(100 * time.Millisecond)
		m.Lock()
		fmt.Printf("Go Work: (%d, %d)\n", id, i)
		m.Unlock()
	}
	m.Lock()
	fmt.Printf("Go Done: %d\n", id)
	m.Unlock()
	atomic.AddInt32(&done, 1)
    w.Done()
}

func main() {
	w.Add(NUM_GS)
	for i := 0; i < NUM_GS; i++ {
		m.Lock()
		fmt.Printf("Go Begin: %d\n", i)
		m.Unlock()
		go test(i)
	}
	w.Wait()
	fmt.Printf("Program Finished: %d\n", done)
}
