package sensitive

import "math/bits"

type FilterSet struct {
	elements []uint64
}

// 构造函数，初始化elements
func NewFilterSet() *FilterSet {
	// 65535 >>> 6 相当于 65535 / 64
	return &FilterSet{
		elements: make([]uint64, 1+(65535>>6)),
	}
}

// 添加一个数字到过滤器
func (fs *FilterSet) Add(no int) {
	fs.elements[no>>6] |= (1 << (no & 63))
}

// 批量添加多个数字到过滤器
func (fs *FilterSet) AddMultiple(no ...int) {
	for _, currNo := range no {
		fs.elements[currNo>>6] |= (1 << (currNo & 63))
	}
}

// 从过滤器中移除一个数字
func (fs *FilterSet) Remove(no int) {
	fs.elements[no>>6] &= ^(1 << (no & 63))
}

// 添加一个数字并返回是否是新添加
func (fs *FilterSet) AddAndNotify(no int) bool {
	eWordNum := no >> 6
	oldElements := fs.elements[eWordNum]
	fs.elements[eWordNum] |= (1 << (no & 63))
	return fs.elements[eWordNum] != oldElements
}

// 从过滤器中移除一个数字并返回是否成功移除
func (fs *FilterSet) RemoveAndNotify(no int) bool {
	eWordNum := no >> 6
	oldElements := fs.elements[eWordNum]
	fs.elements[eWordNum] &= ^(1 << (no & 63))
	return fs.elements[eWordNum] != oldElements
}

// 检查过滤器中是否包含某个数字
func (fs *FilterSet) Contains(no int) bool {
	return (fs.elements[no>>6] & (1 << (no & 63))) != 0
}

// 检查过滤器中是否包含多个数字
func (fs *FilterSet) ContainsAll(no ...int) bool {
	if len(no) == 0 {
		return true
	}
	for _, currNo := range no {
		if (fs.elements[currNo>>6] & (1 << (currNo & 63))) == 0 {
			return false
		}
	}
	return true
}

// 一个非高效的方式来检查多个数字是否都包含
func (fs *FilterSet) ContainsAllUselessWay(no ...int) bool {
	tempElements := make([]uint64, len(fs.elements))
	for _, currNo := range no {
		tempElements[currNo>>6] |= (1 << (currNo & 63))
	}

	for i := 0; i < len(tempElements); i++ {
		if (tempElements[i] & ^fs.elements[i]) != 0 {
			return false
		}
	}
	return true
}

// 获取过滤器中包含的元素数量
func (fs *FilterSet) Size() int {
	size := 0
	for _, element := range fs.elements {
		size += bits.OnesCount64(element)
	}
	return size
}
