package gobalancing

type LoadBalance interface {
	// Next gets next selected item.
	// Next is not goroutine-safe. You MUST use the snchronization primitive to protect it in concurrent cases.
	Next() (item interface{})
	// Add adds a weighted item for selection.
	Add(item interface{}, weight int)
	// All returns all items.
	All() map[interface{}]int
	// RemoveAll removes all weighted items.
	RemoveAll()
	// Reset resets the balancing algorithm.
	Reset()
}
