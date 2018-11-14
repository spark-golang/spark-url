package utils

import "sort"

// Element 用于map的Ksort
type Element struct {
	Key   string
	Value string
}

// Ksort 按照map的key排序
func Ksort(sortMap map[string]string) []Element {
	length := len(sortMap)
	keys := make([]string, length)
	returnMap := make([]Element, length)
	i := 0
	for k := range sortMap {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for k, v := range keys {
		returnMap[k].Value = sortMap[v]
		returnMap[k].Key = v
	}
	return returnMap
}
