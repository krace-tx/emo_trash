package sensitive

type WordNode struct {
	value    int         // 节点名称
	subNodes []*WordNode // 子节点
	isLast   bool        // 是否是敏感词的末尾，默认false
}

// 创建新的 WordNode
func NewWordNode(value int, isLast bool) *WordNode {
	return &WordNode{
		value:    value,
		isLast:   isLast,
		subNodes: nil,
	}
}

// 添加子节点，如果子节点存在，则直接返回子节点，如果不存在则创建并添加新节点
func (node *WordNode) AddIfNoExist(value int, isLast bool) *WordNode {
	// 如果子节点为空，直接添加新的子节点
	if node.subNodes == nil {
		return node.addSubNode(&WordNode{value: value, isLast: isLast})
	}

	// 遍历子节点，判断是否已经存在该节点
	for _, subNode := range node.subNodes {
		if subNode.value == value {
			// 如果存在，并且当前节点还不是末尾，而新增节点是末尾，则更新为末尾
			if !subNode.isLast && isLast {
				subNode.isLast = true
			}
			return subNode
		}
	}

	// 如果不存在，则创建新节点并添加
	return node.addSubNode(&WordNode{value: value, isLast: isLast})
}

// 添加子节点
func (node *WordNode) addSubNode(subNode *WordNode) *WordNode {
	node.subNodes = append(node.subNodes, subNode)
	return subNode
}

// 查询子节点，返回匹配的子节点，如果没有匹配则返回nil
func (node *WordNode) QuerySub(value int) *WordNode {
	if node.subNodes == nil {
		return nil
	}

	for _, subNode := range node.subNodes {
		if subNode.value == value {
			return subNode
		}
	}
	return nil
}

// 判断当前节点是否为最后一个敏感词节点
func (node *WordNode) IsLast() bool {
	return node.isLast
}

// 设置节点为最后一个敏感词节点
func (node *WordNode) SetLast(isLast bool) {
	node.isLast = isLast
}
