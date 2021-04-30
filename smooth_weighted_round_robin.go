package gobalancing

type weightedItem struct {
	Item            interface{}
	Weight          int
	CurrentWeight   int
	EffectiveWeight int
}

type SmoothWeightedRR struct {
	items []*weightedItem
	n     int
}

func (lb *SmoothWeightedRR) Add(item interface{}, weight int) {
	weightedItem := &weightedItem{Item: item, Weight: weight, EffectiveWeight: weight}
	lb.items = append(lb.items, weightedItem)
	lb.n++
}

func (lb *SmoothWeightedRR) RemoveAll() {
	lb.items = lb.items[:0]
	lb.n = 0
}

func (lb *SmoothWeightedRR) Reset() {
	for _, s := range lb.items {
		s.EffectiveWeight = s.Weight
		s.CurrentWeight = 0
	}
}

func (lb *SmoothWeightedRR) All() map[interface{}]int {
	m := make(map[interface{}]int)
	for _, i := range lb.items {
		m[i.Item] = i.Weight
	}
	return m
}

func (lb *SmoothWeightedRR) Next() interface{} {
	i := lb.nextWeightedItem()
	if i == nil {
		return nil
	}
	return i.Item
}

func (lb *SmoothWeightedRR) nextWeightedItem() *weightedItem {
	if lb.n == 0 {
		return nil
	}
	if lb.n == 1 {
		return lb.items[0]
	}

	return nextSmoothWeightedItem(lb.items)
}

func nextSmoothWeightedItem(items []*weightedItem) (best *weightedItem) {
	total := 0

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
