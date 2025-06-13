package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H 类型简化 JSON 数据生成，定义为键值对映射
type H map[string]interface{}

// Context 封装 HTTP 请求上下文，包含请求处理全过程所需信息
type Context struct {
	// 原始对象
	Writer http.ResponseWriter // HTTP 响应写入器
	Req    *http.Request       // HTTP 请求对象

	// 请求信息
	Path   string            // 请求路径
	Method string            // 请求方法
	Params map[string]string // 路由参数

	// 响应信息
	StatusCode int // HTTP 状态码

	// 中间件处理
	handlers []HandlerFunc // 中间件处理器列表
	index    int           // 当前执行中间件索引
}

// newContext 创建并初始化 Context 实例
// 参数:
//
//	w - HTTP 响应写入器
//	req - HTTP 请求对象
//
// 返回:
//
//	*Context - 初始化后的上下文指针
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Path:   req.URL.Path,
		Method: req.Method,
		Req:    req,
		Writer: w,
		index:  -1,
	}
}

// Next 执行下一个中间件处理器，支持中间件的链式调用
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// Fail 中断中间件执行链并返回错误响应
// 参数:
//
//	code - HTTP 状态码
//	err - 错误信息
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

// Param 获取路由参数值
// 参数:
//
//	key - 参数键
//
// 返回:
//
//	string - 参数值
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// PostForm 获取表单参数值
// 参数:
//
//	key - 表单字段键
//
// 返回:
//
//	string - 表单字段值
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 获取 URL 查询参数值
// 参数:
//
//	key - 查询参数键
//
// 返回:
//
//	string - 查询参数值
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status 设置 HTTP 响应状态码
// 参数:
//
//	code - HTTP 状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 设置 HTTP 响应头
// 参数:
//
//	key - 头字段键
//	value - 头字段值
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String 返回纯文本格式响应
// 参数:
//
//	code - HTTP 状态码
//	format - 格式化字符串
//	values - 格式化参数
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON 返回 JSON 格式响应
// 参数:
//
//	code - HTTP 状态码
//	obj - 要序列化的对象
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data 返回原始二进制数据响应
// 参数:
//
//	code - HTTP 状态码
//	data - 二进制数据
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML 返回 HTML 格式响应
// 参数:
//
//	code - HTTP 状态码
//	html - HTML 内容
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
