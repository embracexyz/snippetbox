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