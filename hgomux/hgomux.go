package hgomux

import (
	// "fmt"

	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

//HgoHandleFunc 定义了处理请求的函数格式
type HgoHandleFunc func(hc *HgoContext)

//HgoMux 实现hgo的mux
type HgoMux struct {
	mu      sync.RWMutex
	routers map[string]hgoRouter //map[METHOD]map[path]hgoEntry
}

type hgoRouter map[string]*hgoEntry //路由分配 map[path]hgoEntry

//hgoEntry 保存路由分配的信息
type hgoEntry struct {
	reg    *regexp.Regexp //编译好的正则表达式，用来匹配请求路径
	handle HgoHandleFunc
	parms  []string //路由匹配的名字,例如/user/:name中的name
}

//HgoContext 保存http的writer和reader传入HgoHandleFunc
type HgoContext struct {
	W     http.ResponseWriter
	R     *http.Request
	Parms map[string]string //url中匹配的名字和对应的内容
}

//NewHgoMux 生成新的HgoMux
func NewHgoMux() *HgoMux {
	return &HgoMux{
		routers: make(map[string]hgoRouter),
	}
}

//处理http请求以及路由转发
func (hm *HgoMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//判断是否存在对应的方法池
	router, ok := hm.routers[r.Method]
	if !ok {
		http.NotFound(w, r)
		log.Println("serveHTTP: unexpected request---", r.RequestURI)
		return
	}

	//初始化context
	hc := new(HgoContext)
	hc.W = w
	hc.R = r
	hc.Parms = make(map[string]string)

	for _, entry := range router {
		if entry.reg.MatchString(r.URL.RequestURI()) {
			res := entry.reg.FindAllStringSubmatch(r.URL.RequestURI(), -1)
			// fmt.Println(res)
			for i, parm := range res[0][1:] {
				hc.Parms[entry.parms[i+1]] = parm
			}
			entry.handle(hc)
			break
		}
	}
}

//把方法和路径以及处理函数进行注册到HgoNutex中
// type hgoEntry struct {
// 	reg    *regexp.Regexp //编译好的正则表达式，用来匹配请求路径
// 	handle HgoHandleFunc
// }
func registe(hm *HgoMux, method, pattern string, handle HgoHandleFunc) {

	//匹配path信息处理，构成regexp和一个parm信息
	//例如/user/:name/:id([0-9]+)/:num([0-9]+)
	//regexp编译的为/user/([[:word:]\u4E00-\u9FA5\\s]+)/([0-9]+)/([0-9]+)
	//parm为["", name, id, num]
	//parm第一个为""是为了和regexp.FindAllString()返回的slice的映射对齐,
	//比如返回[[/user/liming/1000/2222 liming 1000 2222]]
	var parms = []string{"zero"}
	var exprs []string

	for _, part := range strings.Split(pattern, "/") {
		if strings.HasPrefix(part, ":") {
			part = part[1:]
			var expr = "([[:word:]\u4E00-\u9FA5\\s]+)"
			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[:index]
			}
			exprs = append(exprs, expr)
			parms = append(parms, part)
		} else {
			exprs = append(exprs, part)
		}
	}
	patt := strings.Join(exprs, "/")
	// fmt.Println(pattern)
	re, err := regexp.Compile(patt)
	if err != nil {
		log.Println("registe:", err)
		return
	}

	hm.mu.Lock()
	defer hm.mu.Unlock()
	//生成METHOD对应的方法池
	router, ok := hm.routers[method]
	if !ok {
		router = make(map[string]*hgoEntry)
	}

	// fmt.Println(parms)
	entry := new(hgoEntry)
	entry.reg = re
	entry.parms = parms
	entry.handle = handle

	router[pattern] = entry

	hm.routers[method] = router
}

//GET 方法处理的路由映射及方法
func (hm *HgoMux) GET(pattern string, handle HgoHandleFunc) {
	registe(hm, "GET", pattern, handle)
}

//HEAD 方法处理的路由映射及方法
func (hm *HgoMux) HEAD(pattern string, handle HgoHandleFunc) {
	registe(hm, "HEAD", pattern, handle)
}

//POST 方法处理的路由映射及方法
func (hm *HgoMux) POST(pattern string, handle HgoHandleFunc) {
	registe(hm, "POST", pattern, handle)
}

//PUT 方法处理的路由映射及方法
func (hm *HgoMux) PUT(pattern string, handle HgoHandleFunc) {
	registe(hm, "PUT", pattern, handle)
}

//DELETE 方法处理的路由映射及方法
func (hm *HgoMux) DELETE(pattern string, handle HgoHandleFunc) {
	registe(hm, "DELETE", pattern, handle)
}

//CONNECT 方法处理的路由映射及方法
func (hm *HgoMux) CONNECT(pattern string, handle HgoHandleFunc) {
	registe(hm, "CONNECT", pattern, handle)
}

//OPTIONS 方法处理的路由映射及方法
func (hm *HgoMux) OPTIONS(pattern string, handle HgoHandleFunc) {
	registe(hm, "OPTIONS", pattern, handle)
}

//TRACE 方法处理的路由映射及方法
func (hm *HgoMux) TRACE(pattern string, handle HgoHandleFunc) {
	registe(hm, "TRACE", pattern, handle)
}

//PATCH 方法处理的路由映射及方法
func (hm *HgoMux) PATCH(pattern string, handle HgoHandleFunc) {
	registe(hm, "PATCH", pattern, handle)
}
