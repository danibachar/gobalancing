package gobalancing

import (
	"errors"
	"sync"
)

type weightedItem struct {
	Item            interface{}
	Weight          float64
	CurrentWeight   float64
	EffectiveWeight float64
}

// type balanceableLocation struct {
// 	item  *weightedItem
// 	index int
// }

type SmoothWeightedRR struct {
	items    []*weightedItem
	itemsMap map[interface{}](*weightedItem)
	sync.RWMutex
}

func NewSWRR() *SmoothWeightedRR {
	return &SmoothWeightedRR{
		items:    make([]*weightedItem, 0),
		itemsMap: make(map[interface{}]*weightedItem),
	}
}

func (lb *SmoothWeightedRR) itemsCount() int {
	return len(lb.items)
}

func (lb *SmoothWeightedRR) Update(item interface{}, weight float64) (err error) {
	lb.Lock()
	defer lb.Unlock()

	var itemToUpdate = lb.itemsMap[item]
	if itemToUpdate == nil {
		return errors.New("No item to update")
	}
	itemToUpdate.Weight = weight
	return nil
}

func (lb *SmoothWeightedRR) Add(item interface{}, weight float64) (err error) {
	lb.Lock()
	defer lb.Unlock()

	if lb.itemsMap[item] != nil {
		return errors.New("Item already in queue")
	}
	weightedItem := &weightedItem{Item: item, Weight: weight, EffectiveWeight: weight}
	lb.itemsMap[item] = weightedItem
	lb.items = append(lb.items, weightedItem)
	return nil
}

func (lb *SmoothWeightedRR) RemoveAll() {
	lb.items = lb.items[:0]
	lb.itemsMap = make(map[interface{}]*weightedItem)
}

func (lb *SmoothWeightedRR) Reset() {
	lb.Lock()
	defer lb.Unlock()

	for _, s := range lb.items {
		s.EffectiveWeight = s.Weight
		s.CurrentWeight = 0
	}
}

func (lb *SmoothWeightedRR) All() map[interface{}]float64 {
	lb.RLock()
	defer lb.RUnlock()

	m := make(map[interface{}]float64)
	for _, i := range lb.items {
		m[i.Item] = i.Weight
	}
	return m
}

func (lb *SmoothWeightedRR) Next() interface{} {
	lb.RLock()
	defer lb.RUnlock()

	i := lb.nextWeightedItem()
	if i == nil {
		return nil
	}
	return i.Item
}

func (lb *SmoothWeightedRR) nextWeightedItem() *weightedItem {
	if lb.itemsCount() == 0 {
		return nil
	}
	if lb.itemsCount() == 1 {
		return lb.items[0]
	}

	return nextSmoothWeightedItem(lb.items)
}

func nextSmoothWeightedItem(items []*weightedItem) (best *weightedItem) {
	total := 0.0

	for i := 0; i < len(items); i++ {
		item := items[i]

		if item == nil {
			continue
		}

		item.CurrentWeight += item.EffectiveWeight
		total += item.EffectiveWeight
		if item.EffectiveWeight < item.Weight {
			item.EffectiveWeight++
		}

		if best == nil || item.CurrentWeight > best.CurrentWeight {
			best = item
		}

	}

	if best == nil {
		return nil
	}

	best.CurrentWeight -= total
	return best
}
