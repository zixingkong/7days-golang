package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 添加路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	// 1. 解析请求路径为分段数组（如/user/profile → ["user", "profile"]）
	searchParts := parsePattern(path)

	// 2. 创建动态参数存储map
	params := make(map[string]string)

	// 3. 获取对应HTTP方法的路由树根节点
	root, ok := r.roots[method]
	if !ok {
		return nil, nil // 方法不存在时返回空
	}

	// 4. 在路由树中搜索匹配节点（核心搜索逻辑）
	n := root.search(searchParts, 0)

	if n != nil {
		// 5. 解析匹配节点的原始模式（如/user/:id）
		parts := parsePattern(n.pattern)

		// 6. 提取动态参数：
		for index, part := range parts {
			// 6.1 处理冒号参数（如:name）
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			// 6.2 处理通配符参数（如*filepath）
			if part[0] == '*' && len(part) > 1 {
				// 捕获剩余所有路径段
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break // 通配符后不再解析
			}
		}
		return n, params // 返回匹配节点和参数
	}

	return nil, nil // 无匹配路由
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
