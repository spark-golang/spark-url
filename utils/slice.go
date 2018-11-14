package utils

import (
	"sort"
	"strconv"
)

// StrInSlice 判断string是否在slice中
func StrInSlice(str string, arr []string) bool {
	if len(arr) == 0 {
		return false
	}
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// IntInSlice 判断int是否在slice中
func IntInSlice(num int, arr []int) bool {
	if len(arr) == 0 {
		return false
	}
	for _, v := range arr {
		if v == num {
			return true
		}
	}
	return false
}

// Interface2String interface的slice转string的slice
func Interface2String(arr []interface{}) []string {
	length := len(arr)
	strSlice := make([]string, length)
	if length == 0 {
		return strSlice
	}
	for k, v := range arr {
		strSlice[k] = v.(string)
	}
	return strSlice
}

// StringToInt string的slice转int
func StringToInt(arr []string) []int {
	length := len(arr)
	intSlice := make([]int, length)
	if length == 0 {
		return intSlice
	}
	for k, v := range arr {
		var ok error
		intSlice[k], ok = strconv.Atoi(v)
		if ok != nil {
			intSlice[k] = 0
		}
	}
	return intSlice
}

// StringToUint32 string的slice转uint32
func StringToUint32(arr []string) []uint32 {
	length := len(arr)
	intSlice := make([]uint32, length)
	if length == 0 {
		return intSlice
	}
	for k, v := range arr {
		tmp, _ := strconv.ParseInt(v, 10, 32)
		intSlice[k] = uint32(tmp)
	}
	return intSlice
}

func StrSliceRemove(slice []string, i int) []string {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

func ReverseStringSlice(slice []string) []string {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

// 用于uid去重
type Uint32Slice []uint32

func (p Uint32Slice) Len() int           { return len(p) }
func (p Uint32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func RemoveRep(slc []uint32) []uint32 {
	if len(slc) <= 1 {
		return slc
	}
	sort.Sort(Uint32Slice(slc))

	var d int
	for i := 1; i < len(slc); i++ {
		if slc[d] != slc[i] {
			d++
			slc[d] = slc[i]
		}
	}
	return slc[:d+1]
}

type Uint64Slice []uint64

func (p Uint64Slice) Len() int           { return len(p) }
func (p Uint64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func RemoveRepUint64(slc []uint64) []uint64 {
	if len(slc) <= 1 {
		return slc
	}
	sort.Sort(Uint64Slice(slc))

	var d int
	for i := 1; i < len(slc); i++ {
		if slc[d] != slc[i] {
			d++
			slc[d] = slc[i]
		}
	}
	return slc[:d+1]
}

func ReOrderSlice(org []uint32, dst []uint32) {
	if len(org) == 0 || len(dst) == 0 {
		return
	}

	var idx int
	for _, did := range dst {
		for i := idx; i < len(org); i++ {
			if did == org[i] {
				org[idx], org[i] = org[i], org[idx]
				idx++
				break
			}
		}
	}
}
