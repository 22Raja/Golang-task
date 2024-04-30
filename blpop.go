package main

import (
	"fmt"
	"sync"
	"time"
)

type Queue struct {
	items []string
	mutex sync.Mutex
	cond  *sync.Cond
}

func NewQueue() *Queue {
	// It is to create an Object
	q := &Queue{
		items: make([]string, 0),
	}

	q.cond = sync.NewCond(&q.mutex)
	return q
}

// It adds to the index.
func (q *Queue) Push(item string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	// Adding to the list
	q.items = append(q.items, item)
	q.cond.Signal()

}

// poping the elements.
func (q *Queue) Pop(timeout time.Duration) string {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for len(q.items) == 0 {
		q.cond.Wait()
	}
	// printing the elements

	item := q.items[0]
	// changing the Index
	q.items = q.items[1:]
	return item
}

func main() {
	queue := NewQueue()

	go func() {
		time.Sleep(2 * time.Second)
		queue.Push("Raaja")
		time.Sleep(1 * time.Second)
		queue.Push("vinoth")
	}()

	go func() {
		for i := 0; i < 2; i++ {
			fmt.Println("Printing:", queue.Pop(time.Second))
		}
	}()

	time.Sleep(5 * time.Second)
}
