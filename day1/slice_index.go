package day1

import "reflect"

type SliceIndex struct {
	Content []*Row
	Pos     int
}

func CreateSliceIndex(at int) Index {
	return &SliceIndex{make([]*Row, 0), at}
}

func (idx *SliceIndex) Insert(row *Row) {
	// 不解决ID的unique特性 问题在于存在多个记录的时候 如何确定删除哪个？
	// 删除的过程需要联动各个索引
	// unique索引是一个特殊的索引 其他索引必定持有unique索引
	// 当下最好的做法是把binarySearch方法和索引的位置解耦 让它可以使用ID
	val := getProperty(*row, idx.Pos)
	at1, exist := idx.binarySearchFirst(row.Id, true)
	// 实际应该插入的位置
	at2, _ := idx.binarySearchFirst(val, false)
	if exist {
		idx.Content[at1] = row
		idx.move(at1, at2)
		return
	}
	// 不按ID乱序插入
	idx.Content = append(idx.Content, row)
	idx.move(len(idx.Content)-1, at2)
}

func (idx *SliceIndex) Remove(id string) {
	at, ok := idx.binarySearchFirst(id, true)
	size := len(idx.Content)
	if ok {
		// 把要删除的元素移动到末尾再删除 保持数组紧凑
		idx.move(at, size-1)
	}
	idx.Content = idx.Content[:size-1]
}

// 返回拷贝数据 修改返回数据不会有副作用
func (idx *SliceIndex) Search(val Col) []Row {
	return idx.SearchByArea(val, val)
}

func (idx *SliceIndex) SearchByArea(start Col, end Col) []Row {
	if LessCol(end, start) {
		return []Row{}
	}
	at1, ok := idx.binarySearchFirst(start, false)
	at2, _ := idx.binarySearchLast(end, false)
	if !ok || at2 < at1 {
		return []Row{}
	}
	size := at2 - at1 + 1
	result := make([]Row, size)
	for i := 0; i < size; i++ {
		result[i] = *idx.Content[at1+i]
	}
	return result
}

func (idx *SliceIndex) move(src, dst int) {
	if src == dst {
		return
	}
	tmp := idx.Content[src]
	for i := src; i < dst; i++ {
		idx.Content[i] = idx.Content[i+1]
	}
	for i := src; i > dst; i-- {
		idx.Content[i] = idx.Content[i-1]
	}
	idx.Content[dst] = tmp
}

// 返回插入位置 以及是否存在
func (idx *SliceIndex) binarySearchFirst(val Col, useId bool) (int, bool) {
	arr, at := idx.Content, idx.Pos
	if useId {
		at = 0
	}
	size := len(arr)
	low, high := 0, size-1
	// 提供一个非零初始值 且不占空间 (已取消) 最后的比对环节不一定对
	// var midVal Col = struct{}{}
	for low <= high {
		mid := low + (high-low)>>1
		midVal := getProperty(*arr[mid], at)
		// 为什么不用三段式？ 因为存在重复的情况 为了总是定位第一个
		if LessCol(midVal, val) {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	if low < size {
		fval := getProperty(*arr[low], at)
		return low, fval == val
	}
	return low, false
}

func (idx *SliceIndex) binarySearchLast(val Col, useId bool) (int, bool) {
	arr, at := idx.Content, idx.Pos
	if useId {
		at = 0
	}

	size := len(arr)
	low, high := 0, size-1
	for low <= high {
		mid := low + (high-low)>>1
		midVal := getProperty(*arr[mid], at)
		// 存在重复的情况 定位第一个
		if GreaterCol(midVal, val) {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	if low < size {
		fval := getProperty(*arr[low], at)
		return high, fval == val
	}
	return high, false
}

// reflect.Value 需要转化成 interface{} 才能类型断言
func getProperty(s Row, at int) Col {
	return reflect.ValueOf(s).FieldByIndex([]int{at}).Interface()
}

func LessCol(idx1, idx2 Col) bool {
	switch t1 := idx1.(type) {
	case int:
		return t1 < idx2.(int)
	case string:
		return t1 < idx2.(string)
	default:
		return false
	}
}

func GreaterCol(idx1, idx2 Col) bool {
	switch t1 := idx1.(type) {
	case int:
		return t1 > idx2.(int)
	case string:
		return t1 > idx2.(string)
	default:
		return false
	}
}
