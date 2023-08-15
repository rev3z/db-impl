# day one ~ three

### 1. 思想
数据库可以看作是一个文件系统

### 2. 实现
在这里我们用CSV实现一个简单的数据库。KV存取是该数据库的核心功能，可以使用哈希表，或者自建数据类型。读写的时间复杂度要尽可能控制在 $O(\log n)$ 以内。另外，这个数据库提供一个范围查询接口。简单的做法是使存储数据按键有序，复杂的做法是建立索引，用有序指针提供有序结果返回。
1. 用CSV文件作为数据库的持久化存储❌
2. 实现KV存取✔️
3. 提供范围查询接口，且支持多个键的索引✔️
4. 切片+二分查找的索引✔️

### 3. 不足
这里的简单数据库实现带来了以下问题
1. 多个键的联合排序
2. 切片实现的索引插入和删除需要移动元素，用跳表或是二叉平衡树的实现会更好
3. 没有完全支持通用的数据类型，实践后发现使用 `type Row []interface{}` 更优雅，把数据类型的解析交给DAO层

### 4. 学到了什么
1. 使用 reflect 反射读取/修改结构体属性
    - 使用 `reflect.ValueOf(obj)` 获取反射
    - 通过属性名读取结构体的字段 `reflect.Value.FieldByName("name")`
    - 通过索引读取结构体的字段 `reflect.Value.FieldByIndex([]int{1})` **切片传参**是为了获取**嵌套**结构体的属性
    - 读取的属性要调用 `reflect.Value.Interface()` 转成 `interface{}` 才能类型断言
    - 调用 `reflect.Value.Elem()` 解包 `interface{}` 或**解引用**指针
    > 在这里理解了 `interface{}` 空接口类型包含着原始值，它实际上是记录指针和类型的结构体
2. 更加理解 *go 语言* 的复制机制
    - 值类型的数据默认是浅拷贝，引用类型的数据默认是深拷贝（切片、Map、指针）
    - 通过 `arr2 := arr1` 复制的切片，是否引用相同的值取决于**底层数组**是否相同（**扩容**会变更底层数组）
    - 使用 `copy()` 函数可以深拷贝切片，另一种做法是逐个深拷贝 `arr = append(arr, node{1 ,"a", 123})`
    - 如何将指针切片转化成对应值的深拷贝切片？采用逐个解引用并深拷贝的方法，但是注意 `node2 := &(*node_ptr)` 是浅拷贝，而 `tmp = *note_ptr; node2 := &tmp` 是深拷贝
3. 使用结构体实现类似继承效果
    - 一个结构体嵌到另一个结构体，称作组合
    - 结构体可以理解为 *JSON Object*
    - 匿名组合是指外层结构体可以直接访问匿名结构体的方法与属性
    - 借助代码自动生成 UML 类图的库可以更加直观地理解复杂的接口与组合关系 [GoPlantUML V2](https://github.com/jfeliu007/goplantuml)
```go
type Live interface {
	eat()
}

type Animal struct {
	name string
	food string
}

func (animal *Animal) eat() {
	fmt.Printf("%v likes eat %v\n", animal.name, animal.food)
}

type Dog struct {
	Animal
	plays string
}

func (dog *Dog) play() {
	fmt.Printf("%v likes play %v\n", dog.name, dog.plays)
}

func main() {
	dog := Dog{Animal{"Bart", "bones"}, "balloon"}
	dog.eat()
	dog.play()
}
```
4. 了解了一个规则： *go 语言* 的 `if` 赋值语句和复合字面量的冲突
    - [官网文档#复合字面量](https://go.dev/ref/spec#Composite_literals)
```go
// 复合字面量作为操作数出现在关键字与 "if"、"for "或 "switch "语句块的开头括号之间
// 并且复合字面量没有用圆括号、方括号或大括号括起来时，就会出现解析歧义。

// 正确写法
if x == (T{a,b,c}[i]) { … }
if (x == T{a,b,c}[i]) { … }

// 以下是具体举例
//  - 错误的例子
if x := node{1, "flower"}; x != val {...}
if x != node{1, "flower"} {...}
// - 改正的例子
if x := (node{1, "flower"}); x != val {...}
if (x != node{1, "flower"}) {...}
```