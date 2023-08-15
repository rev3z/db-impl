package day1

type Col interface{}

type Index interface {
	Insert(row *Row)
	Remove(id string)
	// 返回拷贝数据 修改返回数据不会有副作用
	Search(val Col) []Row
	SearchByArea(start, end Col) []Row
}
