package model

import (
	"sort"
)

func Sort[T any](slice []T, f func(T) bool) []T {

	var result []T
	for _, item := range slice {
		if f(item) {
			result = append(result, item)
		}
	}

	return result
}

type ShowSorterFn = func(any, func(i int, j int) bool)

func NewShowSorter(s []*SortBy) ShowSorterFn {

	return sort.Slice
}

// sort.Slice(list, func(i, j int) bool {
// 	return list[i].Name < list[j].Name
// })
