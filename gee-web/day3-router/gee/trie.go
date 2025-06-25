package gee

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// 插入节点
/*

[思考中]
这段代码实现了在Trie树中插入路由模式的功能。让我一步步解释：

1. **方法签名**：
```go
func (n *node) insert(pattern string, parts []string, height int)
```
- `n *node`：当前操作的Trie节点
- `pattern`：完整的路由模式（如`/user/:name`）
- `parts`：分割后的路由片段（如`["user", ":name"]`）
- `height`：当前处理的层级深度

2. **终止条件**：
```go
if len(parts) == height {
    n.pattern = pattern
    return
}
```
- 当处理深度等于parts长度时，说明已到达路由末端
- 将完整路由模式存储在节点的pattern字段
- 结束递归

3. **获取当前层片段**：
```go
part := parts[height]
```
- 根据当前深度height获取对应的路由片段

4. **查找匹配子节点**：
```go
child := n.matchChild(part)
```
- 在子节点中查找与当前片段匹配的节点
- `matchChild`方法（未展示）应处理通配符匹配逻辑

5. **创建新节点**：
```go
if child == nil {
    child = &node{
        part: part,
        isWild: part[0] == ':' || part[0] == '*'
    }
    n.children = append(n.children, child)
}
```
- 若无匹配子节点，创建新节点
- `isWild`标记该片段是否为通配符（以`:`或`*`开头）
- 将新节点加入当前节点的children列表

6. **递归插入**：
```go
child.insert(pattern, parts, height+1)
```
- 递归处理下一层片段
- 深度height+1进入下一层级

**核心功能**：将路由模式按层级插入Trie树，动态创建节点结构，支持通配符路由的存储和匹配。

是否需要我进一步解释某个具体部分？
*/
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 匹配路由，查询Trie树
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
