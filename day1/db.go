package day1

import "fmt"

type DB struct {
	// 载入内存：这些节点必须要载入内存。
	// 快速读写：可以使用哈希表，哈希表只需要存储对应的指针/下标。
	// 范围查询：数据有序，我们可以使用索引代其排序。
	// 实现的方式有：数组+二分查找 | 二叉搜索树 | 跳表
	Store     map[string]*Row
	IncId     int
	IdIndex   Index
	TimeIndex Index
}

type Row struct {
	Id   string
	Name string
	Time int
}

func CreateDb() *DB {
	db := &DB{make(map[string]*Row), 0, &SliceIndex{make([]*Row, 0), 0}, &SliceIndex{make([]*Row, 0), 2}}
	return db
}

// 插入或更新记录 以ID作为唯一标识
// 如果是更新记录，则将索引删除再插入
func (db *DB) Insert(row Row) {
	// 不用删除旧索引 底层会更新
	// 如果没有ID 则自增分配
	if len(row.Id) == 0 {
		// ID自增
		db.IncId += 1
		// 新记录存入哈希表
		row.Id = fmt.Sprint(db.IncId)
	}
	db.Store[row.Id] = &row
	// 插入索引
	db.IdIndex.Insert(&row)
	db.TimeIndex.Insert(&row)
}

// 删除记录
func (db *DB) Remove(id string) {
	if _, ok := db.Store[id]; ok {
		db.IdIndex.Remove(id)
		db.TimeIndex.Remove(id)
		delete(db.Store, id)
	} else {
		fmt.Println("要删除的记录不存在！")
	}
}

// 按ID查询
func (db *DB) QueryById(id string) []Row {
	return db.IdIndex.Search(id)
}

// 按ID范围查询 返回结果按ID排序
// ID-记录 二分查找到第一个大于begin 和第一个大于end的记录
func (db *DB) QueryByIdArea(begin, end int) []Row {
	return db.IdIndex.SearchByArea(begin, end)
}

// 按时间范围查询 返回结果按时间排序
func (db *DB) QueryByTimeArea(begin, end int) []Row {
	return db.TimeIndex.SearchByArea(begin, end)
}

// 问题
// 1. go语言如何自动匹配CSV的数据类型构建结构体
// 2. 联合键值排序
