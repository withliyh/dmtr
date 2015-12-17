package dmtr

type Item struct {
	Key string
	Val float64
}

type Sorter struct {
	Sorter []Item
}

func NewSorter() *Sorter {
	r := Sorter{Sorter: make([]Item, 0)}
	return &r
}

func (ms *Sorter) Add(key string, val float64) {
	ms.Sorter = append(ms.Sorter, Item{key, val})
}

func (ms *Sorter) Len() int {
	return len(ms.Sorter)
}

func (ms *Sorter) Less(i, j int) bool {
	return ms.Sorter[i].Val < ms.Sorter[j].Val
}

func (ms *Sorter) Swap(i, j int) {
	ms.Sorter[i], ms.Sorter[j] = ms.Sorter[j], ms.Sorter[i]
}
