package mapsort

import "sort"

type MyMapSorter struct {
	Keys   []string
	Values []int
}

func NewMapSorter(m map[string]int) *MyMapSorter {
	vs := &MyMapSorter{
		Keys:   make([]string, 0, len(m)),
		Values: make([]int, 0, len(m)),
	}
	for k, v := range m {
		vs.Keys = append(vs.Keys, k)
		vs.Values = append(vs.Values, v)
	}
	return vs
}

func (vs *MyMapSorter) Sort() {
	sort.Sort(sort.Reverse(vs))
}

func (vs *MyMapSorter) Len() int {
	return len(vs.Values)
}

func (vs *MyMapSorter) Less(i, j int) bool {
	return vs.Values[i] < vs.Values[j]
}

func (vs *MyMapSorter) Swap(i, j int) {
	vs.Values[i], vs.Values[j] = vs.Values[j], vs.Values[i]
	vs.Keys[i], vs.Keys[j] = vs.Keys[j], vs.Keys[i]
}
