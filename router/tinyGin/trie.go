package tinyGin

import "strings"

type node struct {
	pattern  string  // 完整的路由路径，只有在某一个匹配路由规则最后一个节点才有值，例如 /p/:lang
	part     string  // 路由路径中的某一部分，例如 :lang
	children []*node // 当前节点的子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，例如 part含有 : 或 * 时为true
}

// 找到第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 找到所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 开发服务时，注册路由规则，映射handler；访问时，匹配路由规则，查找到对应的handler
// 因此，Trie 树需要支持节点的插入与查询。

// 插入功能：递归查找每一层的节点，如果没有匹配到当前part的节点，则新建一个
func (n *node) insert(pattern string, parts []string, height int) {
	// 到了叶子节点，给其pattern字段赋值
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	// 获取当前height位置的part
	part := parts[height]
	// 是否有与这个part相等的child
	child := n.matchChild(part)
	if child == nil {
		// 如果没有的话就创建一个新的，即插入的新节点
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	// 递归插入
	child.insert(pattern, parts, height+1)
}

// 查询功能，同样也是递归查询每一层的节点，退出规则是，匹配到了*或者匹配到了第len(parts)层节点，匹配失败
func (n *node) search(parts []string, height int) *node {
	// 查完parts或遇到通配符，判断这个node是否有pattern(判断叶子结点)，如果有说明存在这样一条路径，如果没有说明没有这样一条路径
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	// 获取当前height位置的part
	part := parts[height]
	// 查询当前node与part相等的所有children节点
	children := n.matchChildren(part)
	// 遍历符合条件的子节点，递归查询
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
