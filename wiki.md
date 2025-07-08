1. 项目概述
7-days-golang 是包含六个核心子项目的 Go 语言实践集合，覆盖分布式系统、存储引擎等关键技术领域。项目采用模块化设计，各子项目既可独立运行又可组合使用，基于 Go 标准库实现兼具教学性与生产可用性。

子项目核心功能
1. gee-rpc - 服务发现与负载均衡的远程调用框架
实现基于注册中心的服务发现机制，提供多种负载均衡策略（如随机选择）。核心特性包括：

支持并发RPC调用（通过call函数实现）

自动化的服务注册与发现

基于HTTP协议的远程方法调用

上下文（Context）支持的超时控制

classDiagram
    class XClient {
        +Call(method string, args interface{}, reply interface{}) error
        +Close() error
    }
    class GeeRegistryDiscovery {
        +GetServers() []string
        +WatchServers() chan []string
    }
    XClient --> GeeRegistryDiscovery : 依赖发现服务
2. gee-cache - 分布式缓存系统
模仿groupcache实现的轻量级缓存系统，核心功能包括：

单机/分布式双模式运行

LRU缓存淘汰策略

一致性哈希节点选择

缓存击穿防护（通过Group结构体实现）

图片
源码
classDiagram
    class Group {
        +name string
        +getter Getter
        +mainCache cache
        +peers PeerPicker
        +Get(key string) (ByteView, error)
    }
    class PeerPicker {
        +PickPeer(key string) (PeerGetter, bool)
    }
    Group --> PeerPicker : 节点选择
3. gee-web - 支持中间件链的Web框架
通过RouterGroup结构体实现的路由管理系统：

支持中间件链式调用

嵌套路由组管理

RESTful路由注册

上下文（Context）封装

图片
源码
classDiagram
    class RouterGroup {
        +prefix string
        +middlewares []HandlerFunc
        +parent *RouterGroup
        +engine *Engine
        +Use(...HandlerFunc)
        +Group(prefix string) *RouterGroup
    }
    class Engine {
        +ServeHTTP(ResponseWriter, *Request)
    }
    RouterGroup --> Engine : 共享引擎实例
4. gee-orm - 轻量级数据库操作封装
Engine结构体提供的事务管理能力：

多数据库方言支持

会话(Session)生命周期管理

ACID事务支持

连接池管理

5. gee-bolt - B+树键值存储（根据上下文推断）
（注：文档中未提供具体实现细节，基于项目定位描述）

基于B+树的持久化存储

事务性数据操作

内存映射文件支持

6. demo-wasm - WebAssembly跨语言交互（根据上下文推断）
（注：文档中未提供具体实现细节，基于项目定位描述）

Go代码编译为WASM模块

浏览器环境交互演示

跨语言函数调用

各子项目依赖关系热图：

图片
源码
graph TD
    gee-rpc --> |HTTP通信| gee-cache
    gee-web --> |数据缓存| gee-cache
    gee-orm --> |存储引擎| gee-bolt
    demo-wasm --> |独立演示| gee-web
Sources: geecache.md, line 45, main.go, line 63, geecache.go, line 10, gee.go, line 15, geeorm.go, line 11

模块化设计理念
1. 子项目独立编译部署能力
各子项目采用高内聚设计，具备完整的独立运行能力：

编译隔离：每个子项目拥有独立的main.go入口文件，例如gee-web可通过go build ./gee-web单独编译

依赖最小化：核心模块（如gee-cache）仅依赖Go标准库，无第三方包约束

配置独立：通过环境变量或配置文件实现参数隔离（如gee-rpc的注册中心地址）

图片
源码
graph TD
    gee-web --> |独立运行| HTTP服务
    gee-cache --> |单机模式| 内存缓存
    gee-bolt --> |嵌入式| 本地存储
2. 跨模块协作模式
模块间通过标准化接口实现能力组合：

gee-cache 整合案例
classDiagram
    class Group {
        +Get(key string)
    }
    class PeerPicker {
        +PickPeer(key string)
    }
    class consistenthash.Map {
        +Add(keys...string)
        +Get(key string)
    }
    Group --> PeerPicker : 分布式节点选择
    PeerPicker --> Map : 一致性哈希路由
关键协作点：

LRU缓存：通过mainCache字段实现本地缓存淘汰

一致性哈希：Map结构体处理节点定位

防击穿：Group.Get方法协调本地/远程数据获取

3. 接口标准化设计
核心接口定义实现模块解耦：

gee-orm 方言抽象层
classDiagram
    class Dialect {
        <<interface>>
        +DataTypeOf() string
        +TableExistSQL() (string, []interface{})
    }
    class MySQLDialect
    class SQLiteDialect
    Dialect <|-- MySQLDialect
    Dialect <|-- SQLiteDialect
接口规范：

方法名	作用域	契约要求
DataTypeOf	类型转换	必须返回有效的SQL类型字符串
TableExistSQL	元数据操作	返回可执行的SQL语句模板
4. 组合使用案例
gee-web 中间件嵌套实现
图片
源码
sequenceDiagram
    participant RouterGroup
    participant HandlerFunc
    RouterGroup->>HandlerFunc: Use(middleware1)
    RouterGroup->>HandlerFunc: Use(middleware2)
    Note right of RouterGroup: 中间件按添加顺序执行
    HandlerFunc-->>RouterGroup: 链式调用完成
典型组合场景：

路由分组：通过Group(prefix)创建嵌套路由组

中间件叠加：Use方法支持多中间件追加

// 示例：API路由组添加认证+日志中间件
api := engine.Group("/api")
api.Use(AuthMiddleware(), LogMiddleware())
Sources: gee.go, line 61, gee.go, line 54, geecache.go, line 10, dialect.go, line 8, consistenthash.go, line 13

教学与生产双特性
1. 标准库实现的简洁性设计
1.1 LRU缓存的高效实现
采用哈希表+双向链表的经典组合，实现O(1)时间复杂度的缓存操作：

classDiagram
    class Cache {
        -maxBytes int64
        -ll *list.List
        -cache map[string]*list.Element
        -OnEvicted func(string, Value)
        +Add(key string, value Value)
        +Get(key string) (value Value, ok bool)
        +RemoveOldest()
    }
    class list.Element {
        +Value interface{}
    }
    Cache --> list.Element : 存储节点引用
关键设计点：

内存控制：通过maxBytes字段实现容量限制

淘汰策略：RemoveOldest方法自动移除LRU项

事件回调：OnEvicted支持淘汰时的自定义处理

1.2 轻量级RPC框架
基于标准库net/rpc的扩展实现：

图片
源码
graph TD
    Client[Client] -->|TCP连接| Server[Server]
    Server -->|注册中心| Registry[RegistryDiscovery]
    Client -->|负载均衡| Registry
2. 渐进式教学模块设计
2.1 分阶段实现路径
以gee-cache为例的开发演进：

图片
源码
gantt
    title GeeCache开发阶段
    dateFormat  YYYY-MM-DD
    section 基础功能
    LRU实现       :done, day1, 2023-01-01, 1d
    单机缓存      :done, day2, 2023-01-02, 1d
    section 分布式
    HTTP通信      :done, day3, 2023-01-03, 1d
    一致性哈希    :done, day4, 2023-01-04, 1d
    section 高级特性
    Protobuf编码  :active, day7, 2023-01-07, 1d
2.2 缓存系统流程演进
图片
源码
flowchart TD
    A[接收key] --> B{缓存命中?}
    B -->|是| C[返回缓存值]
    B -->|否| D{远程节点?}
    D -->|是| E[HTTP请求远程节点]
    D -->|否| F[执行回调函数]
    E --> G{成功?}
    G -->|否| F
3. 生产环境保障机制
3.1 高并发处理
RPC客户端的并发请求管理：

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            client.Call("Foo.Sum", args, &reply) // 并发安全调用
        }(i)
    }
    wg.Wait()
}
3.2 事务原子性保障
ORM引擎的事务处理流程：

图片
源码
stateDiagram
    [*] --> Begin
    Begin --> Executing: 开始事务
    Executing --> Committed: 执行成功
    Executing --> Rollback: 发生错误
    Executing --> Rollback: 发生panic
    Committed --> [*]
    Rollback --> [*]
4. 性能基准特性
4.1 核心操作时间复杂度
操作类型	实现方式	时间复杂度
缓存查询	哈希表查找	O(1)
节点定位	一致性哈希	O(log n)
事务提交	数据库原生事务	O(1)
4.2 关键性能设计
缓存预加载：通过Group.Get方法实现热点数据预取

连接复用：RPC客户端维护持久化TCP连接

零拷贝处理：缓存值使用ByteView只读封装

图片
源码
classDiagram
    class ByteView {
        -b []byte
        +String() string
        +ByteSlice() []byte
        +Len() int
    }
    class Group {
        +Get(key string) (ByteView, error)
    }
    Group --> ByteView : 返回不可变数据
Sources: geecache-day5.md, line 29, lru.go, line 26, geeorm.go, line 58, lru_test.go, line 25, main.go, line 23

技术领域覆盖
1. 分布式系统实现
1.1 一致性哈希（Consistent Hashing）
通过Map结构体实现虚拟节点分布：

图片
源码
classDiagram
    class Map {
        +hash Hash
        +replicas int
        +keys []int
        +hashMap map[int]string
        +Add(keys...string)
        +Get(key string) string
    }
关键特性：

虚拟节点：通过replicas字段控制副本数量

快速查找：排序后的keys数组实现二分查找

哈希抽象：支持自定义Hash函数实现

1.2 服务发现（Service Discovery）
RPC框架的服务定位流程：

图片
源码
sequenceDiagram
    participant Client
    participant Discovery
    participant Server
    Client->>Discovery: GetServers()
    Discovery-->>Client: ["server1","server2"]
    Client->>Server: RPC Call
    Server-->>Client: Response
2. 存储引擎实现
2.1 缓存管理（Cache Management）
cache结构体的线程安全设计：

图片
源码
classDiagram
    class cache {
        -mu sync.Mutex
        -lru *lru.Cache
        -cacheBytes int64
        +get(key string) (ByteView, bool)
        +add(key string, value ByteView)
    }
性能保障：

并发控制：mu互斥锁保护LRU操作

容量限制：cacheBytes严格限制内存使用

零拷贝：ByteView封装保证数据不可变

2.2 B+树存储（B+ Tree Storage）
gee-bolt的存储架构：

图片
源码
flowchart TD
    A[写入请求] --> B{内存缓冲}
    B -->|满| C[刷盘操作]
    B -->|未满| D[更新内存]
    C --> E[B+树索引更新]
3. Web开发实现
3.1 中间件链（Middleware Chain）
ServeHTTP方法的处理流程：

图片
源码
flowchart LR
    A[请求进入] --> B[收集匹配中间件]
    B --> C[创建Context]
    C --> D[执行处理链]
    D --> E[路由处理]
3.2 路由分组（Route Grouping）
RouterGroup的嵌套结构：

图片
源码
classDiagram
    class RouterGroup {
        +prefix string
        +middlewares []HandlerFunc
        +parent *RouterGroup
        +Use(...HandlerFunc)
        +Group(string) *RouterGroup
    }
4. 数据库操作实现
4.1 ORM会话管理
Session结构体的核心能力：

图片
源码
classDiagram
    class Session {
        +db *sql.DB
        +dialect Dialect
        +tx *sql.Tx
        +QueryRows()
        +Exec()
        +Raw()
    }
4.2 SQL构建器
子句生成过程：

图片
源码
sequenceDiagram
    Session->>Clause: Set(table)
    Session->>Clause: Set(where)
    Clause-->>Session: SQL语句
5. 跨语言交互实现
5.1 WASM编译流程
图片
源码
flowchart TD
    A[Go代码] --> B[WASM编译]
    B --> C[浏览器加载]
    C --> D[JS交互]
6. 网络编程实现
6.1 RPC通信协议
图片
源码
sequenceDiagram
    Client->>Server: Call("Method", args)
    Server->>Client: Return(reply)
6.2 HTTP缓存通信
ServeHTTP的缓存获取流程：

图片
源码
flowchart TD
    A[请求路径解析] --> B[获取Group]
    B --> C[缓存查询]
    C -->|命中| D[返回数据]
    C -->|未命中| E[回调加载]
模块依赖热图
图片
源码
graph TD
    gee-web -->|最高频| gee-cache
    gee-rpc -->|中频| gee-orm
    gee-bolt -->|低频| demo-wasm
Sources: gee.go, line 83, consistenthash.go, line 13, http.go, line 46, raw.go, line 14, main.go, line 15, cache.go, line 8