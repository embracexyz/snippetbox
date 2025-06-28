# 路由注册

**注意！/static/需要添加/，表示前缀匹配，不然就只会精准匹配到/static 路径，导致找不到css js等静态文件**



```bash
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
```

# 中间件

# 表单处理

1. 表单url
2. 如何做表单validation，报错时返回提示并携带之前填入的信息

| method | url               | handler           |
| ------ | ----------------- | ----------------- |
| GET    | /snippet/create   | snippetCreate     |
| POST   | /snippet/create   | snippetCreatePost |
| GET    | /snippet/view/:id | snippetView       |

