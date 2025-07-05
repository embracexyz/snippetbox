# day1

1. 项目组织结构构建
   1. cmd
   2. internal
   3. ui
2. http server3个组成部分
   1. httpServer（处理http的连接管理）
   2. serverMux（管理路由和handler的映射）
   3. handler（业务处理函数）
3. 如何限制request的请求方法
4. http请求
   1. 请求参数query 参数的查询方法
5. http响应
   1. 如何写回响应
   2. 如何定义响应状态码：w.WriteHeader()
   3. 如何定义响应头部信息：w.Header()；其本质是`map[string][]string`
   4. http.Error()的helper函数使用
6. 如何限制/ 是否只允许精确匹配
7. 模版渲染
   1. 

# request context/session

1. 只在本次reqeust的生命周期中，只在内存里存储
2. session是短期存储在外部存储（redis、mysql中）

```bash
app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
# 这句其实是将authenticatedUserID：id存储为session data，生成对应的session token，存在外部存储中，比如redis，然后通过r返回给客户端，并不把数据存在r.Context()中，只在其中已有session data情况下，用于更新

return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
# 这里也是从context中获取session token，然后从session token取session data，再从session data中取authenticatedUserID




```

## context 用于认证

```bash
func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}
```

1. 如上方法，有个缺陷，当session data中取出authenticatedUserID时，该用户已经在数据库删除了，这时是感知不到的
2. 所以需要通过查数据库实现，但调用方可能有多次，对数据库有压力，那么就需要一次查询，后续handler通过一个标志位判断结果true or false就行了，于是加入一个中间件用于查库，然后将标志信息set 到context，后续handler只需在context查询即可
3. 该中间件authenticate做的确保用户存于库中，所以还需要先从session data取出authenticatedUserID，拿到id，才能作为查询条件，而没有authenticatedUserID的，说明还没登录，直接放行（后续有requireAuth 中间件拦截，该中间件又依赖authenticate设置在context的标志位，所以是负责2个不同逻辑的中间件：一个是负责设置标志位，一个是根据标志位选择性拦截）