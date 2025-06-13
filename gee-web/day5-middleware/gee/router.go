package gee

import (
	"net/http"
	"strings"
)

// router 路由管理器，负责路由的注册、匹配和请求处理
type router struct {
	roots    map[string]*node       // 各HTTP方法对应的路由前缀树根节点
	handlers map[string]HandlerFunc // 路由模式到处理函数的映射
}

// newRouter 创建并初始化路由管理器实例
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// parsePattern 解析路由模式，分割路径并处理通配符
// 参数:
//
//	pattern - 原始路由模式字符串（如 /user/:name/*）
//
// 返回:
//
//	[]string - 解析后的路径段数组（如 ["user", ":name", "*"]）
//
// 注意:
//   - 每个模式只允许包含一个通配符*
//   - 遇到通配符后停止解析后续路径
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' { // 遇到通配符提前终止
				break
			}
		}
	}
	return parts
}

// addRoute 注册路由到路由树
// 参数:
//
//	method  - HTTP方法（GET/POST等）
//	pattern - 路由模式字符串
//	handler - 对应的处理函数
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern // 构建路由唯一标识
	// 初始化对应方法的路由树
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	// 插入路由到前缀树
	r.roots[method].insert(pattern, parts, 0)
	// 存储路由处理函数
	r.handlers[key] = handler
}

// getRoute 根据请求路径查找路由节点并解析参数
// 参数:
//
//	method - HTTP方法
//	path   - 请求路径
//
// 返回:
//
//	*node          - 匹配的路由节点
//	map[string]string - 解析到的路径参数
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		// 解析路径参数
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' { // 解析冒号参数
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 { // 解析通配符参数
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

// getRoutes 获取指定方法的所有注册路由节点
// 参数:
//
//	method - HTTP方法
//
// 返回:
//
//	[]*node - 路由节点列表
func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes) // 遍历路由树收集所有节点
	return nodes
}

// handle 处理请求入口
// 参数:
//
//	c - 请求上下文
func (r *router) handle(c *Context) {
	// 查找匹配的路由
	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		// 找到路由时设置参数并添加处理函数
		key := c.Method + "-" + n.pattern
		c.Params = params
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		// 未找到路由时添加404处理函数
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	// 启动中间件/处理函数执行链
	c.Next()
}
