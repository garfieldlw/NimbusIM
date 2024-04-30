package utils

import "strings"

func IsInList(slice []string, s string) bool {
	for _, v := range slice {
		if strings.EqualFold(v, s) {
			return true
		}
	}
	return false
}

type Int32Slice []int32

func (p Int32Slice) Len() int {
	return len(p)
}
func (p Int32Slice) Less(i, j int) bool {
	return p[i] < p[j]
}
func (p Int32Slice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type Int64Slice []int64

func (p Int64Slice) Len() int {
	return len(p)
}
func (p Int64Slice) Less(i, j int) bool {
	return p[i] < p[j]
}
func (p Int64Slice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func ContainsString(v []string, c string) bool {
	for _, t := range v {
		if t == c {
			return true
		}
	}

	return false
}

func ContainsInt64(v []int64, c int64) bool {
	for _, t := range v {
		if t == c {
			return true
		}
	}

	return false
}

func ContainsInt32(v []int32, c int32) bool {
	for _, t := range v {
		if t == c {
			return true
		}
	}

	return false
}

func SliceStringAppendIndex(slice []string, value string, index int) []string {
	if slice == nil || len(slice) == 0 || len(slice) <= index {
		return append(slice, value)
	}

	slice = append(slice[:index+1], slice[index:]...) // index < len(a)
	slice[index] = value
	return slice
}

func SliceIntersectInt64(slice1, slice2 []int64) []int64 {
	res := make([]int64, 0)
	map1 := map[int64]int{}
	for _, v := range slice1 {
		map1[v] += 1
	}
	for _, v := range slice2 {
		if map1[v] > 0 {
			res = append(res, v)
			map1[v] -= 1
		}
	}
	return res
}

func SliceDifferenceInt64(slice1, slice2 []int64) []int64 {
	m := make(map[int64]int64)
	for _, v := range slice1 {
		m[v] = v
	}
	for _, v := range slice2 {
		if m[v] > 0 {
			delete(m, v)
		}
	}
	var str []int64
	for _, s2 := range m {
		str = append(str, s2)
	}
	return str
}
