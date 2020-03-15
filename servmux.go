package hgo

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

// ServMux used for http router
type ServMux struct {
	mu sync.RWMutex
	rs map[string]router
}

// router distribution data
type router map[string]*entry

func newRouter() router {
	return make(router)
}

// HandleFunc defines the handle function
type HandleFunc func(c *Context)

// entry contains each router's data
type entry struct {
	reg    *regexp.Regexp // compiled to match the url path
	handle HandleFunc
	param  []string //name matched in url' path, for instance: name in /user/:name
}

func newEntry() *entry {
	return new(entry)
}

// NewServMux create a new ServMux
func NewServMux() *ServMux {
	return &ServMux{
		rs: make(map[string]router),
	}
}

// ServeHTTP handle the http request and process router distribution
func (sm *ServMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r.RequestURI)

	// if there is method pool
	rou, ok := sm.rs[r.Method]
	if !ok {
		http.NotFound(w, r)
		log.Println("serveHTTP: unexpected request---", r.RequestURI)
		return
	}

	// context initialize
	c := newContext(w, r)

	// regexp match processing
	for _, en := range rou {
		if en.reg.MatchString(r.URL.RequestURI()) {
			res := en.reg.FindAllStringSubmatch(r.URL.RequestURI(), -1)
			// fmt.Println(res)
			for i, parm := range res[0][1:] {
				c.Param[en.param[i+1]] = parm
			}
			en.handle(c)
			break
		}
	}
}

// registe the method and the router entry into ServMux
// type hgoEntry struct {
// 	reg    *regexp.Regexp //编译好的正则表达式，用来匹配请求路径
// 	handle HandleFunc
// }
func registe(sm *ServMux, method, pattern string, handle HandleFunc) {

	//process the path matching data, forming regexp and parm
	//such as : '/user/:name/:id([0-9]+)/:num([0-9]+)'
	//regexp compiled: '/user/([[:word:]\u4E00-\u9FA5\\s]+)/([0-9]+)/([0-9]+)'
	//parm is ["", name, id, num]
	//parm[0] is ""to match up with the slice returened by regexp.FindAllString()
	//as return [[/user/liming/1000/2222 liming 1000 2222]]
	var parms = []string{""}
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

	sm.mu.Lock()
	defer sm.mu.Unlock()
	//method pool for "METHOD"
	rou, ok := sm.rs[method]
	if !ok {
		rou = newRouter()
	}

	// fmt.Println(parms)
	en := newEntry()
	en.reg = re
	en.param = parms
	en.handle = handle

	rou[pattern] = en

	sm.rs[method] = rou
}

//GET method and router distribution
func (sm *ServMux) GET(pattern string, handle HandleFunc) {
	registe(sm, "GET", pattern, handle)
}

//HEAD method and router distribution
func (sm *ServMux) HEAD(pattern string, handle HandleFunc) {
	registe(sm, "HEAD", pattern, handle)
}

//POST method and router distribution
func (sm *ServMux) POST(pattern string, handle HandleFunc) {
	registe(sm, "POST", pattern, handle)
}

//PUT method and router distribution
func (sm *ServMux) PUT(pattern string, handle HandleFunc) {
	registe(sm, "PUT", pattern, handle)
}

//DELETE method and router distribution
func (sm *ServMux) DELETE(pattern string, handle HandleFunc) {
	registe(sm, "DELETE", pattern, handle)
}

//CONNECT method and router distribution
func (sm *ServMux) CONNECT(pattern string, handle HandleFunc) {
	registe(sm, "CONNECT", pattern, handle)
}

//OPTIONS method and router distribution
func (sm *ServMux) OPTIONS(pattern string, handle HandleFunc) {
	registe(sm, "OPTIONS", pattern, handle)
}

//TRACE method and router distribution
func (sm *ServMux) TRACE(pattern string, handle HandleFunc) {
	registe(sm, "TRACE", pattern, handle)
}

//PATCH method and router distribution
func (sm *ServMux) PATCH(pattern string, handle HandleFunc) {
	registe(sm, "PATCH", pattern, handle)
}
