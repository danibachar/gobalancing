package gobalancing

import (
	"sync"
)

// SmoothWeightedRRThreadSafe - Load Balancer implementation, Smooth Weighted Round Robin
type SmoothWeightedRRThreadSafe struct {
	lb *SmoothWeightedRR
	sync.RWMutex
}

// NewSWRRhreadSafe - NewSWRRhreadSafe - Load Balancer constructor
func NewSWRRhreadSafe() *SmoothWeightedRRThreadSafe {
	return &SmoothWeightedRRThreadSafe{
		lb: NewSWRR(),
	}
}

func (lb *SmoothWeightedRRThreadSafe) itemsCount() int {
	lb.RLock()
	defer lb.RUnlock()

	return len(lb.lb.items)
}

// Update - updates the weight of an existing item
func (lb *SmoothWeightedRRThreadSafe) Update(item interface{}, weight float64) (err error) {
	lb.Lock()
	defer lb.Unlock()

	return lb.lb.Update(item, weight)
}

// Add - adds a new unique item to the list
func (lb *SmoothWeightedRRThreadSafe) Add(item interface{}, weight float64) (err error) {
	lb.Lock()
	defer lb.Unlock()

	return lb.lb.Add(item, weight)

}

// RemoveAll - removes all items and reset state
func (lb *SmoothWeightedRRThreadSafe) RemoveAll() {
	lb.Lock()
	defer lb.Unlock()

	lb.lb.RemoveAll()
}

// Reset - reset load balancing state, like adding all items with their original weights
func (lb *SmoothWeightedRRThreadSafe) Reset() {
	lb.Lock()
	defer lb.Unlock()

	lb.lb.Reset()
}

// All - returns a map all items as keys and their weights as values
func (lb *SmoothWeightedRRThreadSafe) All() map[interface{}]float64 {
	lb.RLock()
	defer lb.RUnlock()

	return lb.lb.All()
}

// Next - fetches next item according to the smooth weighted round robin fashion
func (lb *SmoothWeightedRRThreadSafe) Next() interface{} {
	lb.RLock()
	defer lb.RUnlock()

	return lb.lb.Next()
}
