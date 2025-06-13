package gee

import (
	"log"
	"net/http"
	"strings"
)

// HandlerFunc 定义请求处理函数类型，接收Context参数
type HandlerFunc func(*Context)

// Engine 实现http.Handler接口，作为Web框架核心
type (
	// RouterGroup 路由组，支持中间件和嵌套路由组
	RouterGroup struct {
		prefix      string        // 路由组前缀（支持嵌套）
		middlewares []HandlerFunc // 注册的中间件列表
		parent      *RouterGroup  // 父级路由组（用于嵌套）
		engine      *Engine       // 所有路由组共享的引擎实例
	}

	// Engine 引擎，聚合路由组并管理所有路由
	Engine struct {
		*RouterGroup                // 根路由组
		router       *router        // 路由树（实际路由存储）
		groups       []*RouterGroup // 维护所有路由组列表
	}
)

// New 创建并初始化Gee引擎实例
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup} // 初始化时包含根路由组
	return engine
}

// Group 创建新的路由组（支持嵌套路由组）
// 参数prefix为当前组的路径前缀，会自动与父级前缀拼接
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix, // 拼接完整前缀
		parent: group,                 // 设置父级路由组
		engine: engine,                // 共享引擎实例
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// Use 为路由组注册中间件（支持多个中间件）
// 中间件按照注册顺序执行
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// addRoute 私有方法：注册路由到路由树
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp // 拼接完整路由路径
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET 注册GET请求路由（快捷方法）
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST 注册POST请求路由（快捷方法）
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Run 启动HTTP服务器（封装http.ListenAndServe）
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP 处理HTTP请求（实现http.Handler接口）
// 1. 合并匹配的中间件
// 2. 创建请求上下文
// 3. 执行处理链
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	// 收集所有匹配请求路径前缀的路由组中间件
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	// 创建请求上下文
	c := newContext(w, req)
	// 将中间件链赋值给上下文
	c.handlers = middlewares
	// 将请求交给路由处理
	engine.router.handle(c)
}
