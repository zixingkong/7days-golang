package gee

import (
	"fmt"
	"strings"
)

// node 表示路由前缀树节点，用于高效路由匹配
type node struct {
	pattern  string  // 完整路由路径（仅在叶子节点设置）
	part     string  // 当前节点对应的路由部分
	children []*node // 子节点列表
	isWild   bool    // 是否是通配节点（包含:或*）
}

// String 实现节点字符串表示，用于调试
func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// insert 递归插入路由路径到前缀树
// pattern: 完整路由路径（如/user/:name）
// parts: 路径分割后的部分（如["user", ":name"]）
// height: 当前处理的部分索引
func (n *node) insert(pattern string, parts []string, height int) {
	// 到达路径末尾时设置完整模式
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	// 查找匹配的子节点
	child := n.matchChild(part)

	// 没有匹配则创建新节点
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*', // 判断是否是参数节点或通配节点
		}
		n.children = append(n.children, child)
	}

	// 递归插入下一层路径
	child.insert(pattern, parts, height+1)
}

// search 递归搜索匹配路由路径的节点
// parts: 请求路径分割后的部分
// height: 当前处理的部分索引
func (n *node) search(parts []string, height int) *node {
	// 终止条件：到达路径末尾或遇到通配符*
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil // 非叶子节点不匹配
		}
		return n
	}

	part := parts[height]
	// 获取所有可能匹配的子节点（包含通配节点）
	children := n.matchChildren(part)

	// 递归搜索所有匹配的子节点
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

// travel 遍历收集所有路由节点（用于调试或展示路由列表）
func (n *node) travel(list *([]*node)) {
	// 收集有完整模式的叶子节点
	if n.pattern != "" {
		*list = append(*list, n)
	}
	// 递归遍历子节点
	for _, child := range n.children {
		child.travel(list)
	}
}

// matchChild 查找第一个匹配的子节点（用于插入）
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren 查找所有匹配的子节点（用于搜索）
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
