#  Hgo是一个HTTP框架

##  HgoMutex实现了自定义的路由转发功能
```
//GET 方法处理的路由映射及方法
func (hm *ServMux) GET(pattern string, handle HgoHandleFunc) {
	registe(hm, "GET", pattern, handle)
}

//HEAD 方法处理的路由映射及方法
func (hm *ServMux) HEAD(pattern string, handle HgoHandleFunc) {
	registe(hm, "HEAD", pattern, handle)
}

//POST 方法处理的路由映射及方法
func (hm *ServMux) POST(pattern string, handle HgoHandleFunc) {
	registe(hm, "POST", pattern, handle)
}

//PUT 方法处理的路由映射及方法
func (hm *ServMux) PUT(pattern string, handle HgoHandleFunc) {
	registe(hm, "PUT", pattern, handle)
}

//DELETE 方法处理的路由映射及方法
func (hm *ServMux) DELETE(pattern string, handle HgoHandleFunc) {
	registe(hm, "DELETE", pattern, handle)
}

//CONNECT 方法处理的路由映射及方法
func (hm *ServMux) CONNECT(pattern string, handle HgoHandleFunc) {
	registe(hm, "CONNECT", pattern, handle)
}

//OPTIONS 方法处理的路由映射及方法
func (hm *ServMux) OPTIONS(pattern string, handle HgoHandleFunc) {
	registe(hm, "OPTIONS", pattern, handle)
}

//TRACE 方法处理的路由映射及方法
func (hm *ServMux) TRACE(pattern string, handle HgoHandleFunc) {
	registe(hm, "TRACE", pattern, handle)
}

//PATCH 方法处理的路由映射及方法
func (hm *ServMux) PATCH(pattern string, handle HgoHandleFunc) {
	registe(hm, "PATCH", pattern, handle)
}
```

使用这些函数进行注册，注册时可以使用正则表达式进行参数匹配

例如/user/:name/:id([0-9]+)/:num([0-9]+)中id的格式为([0-9]+)，这样处理之后id会在HgoContext中的Param中保存。

可以使用如下代码进行初始化

```
	hm := ServMux.NewServMux()
	hm.GET("/user/:name/:year", sayHello)
	hm.POST("/user/:name([a-z])+", sayHi)
	http.ListenAndServe(":4000", hm)
```
该框架类似于httpRouter但是实现了更复杂的参数匹配