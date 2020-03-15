package hgo

import (
	"crypto/tls"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

// Context contains elements for request and response and so on
type Context struct {
	// get all the method of request
	r *http.Request
	// get all the method of responsewriter
	w http.ResponseWriter

	// contains the params mathched in url
	Param map[string]string

	// Method inform the request method from client(GET, POST or so on)
	Method string

	// URL specifies either the URI being requested (for server
	// requests) or the URL to access (for client requests).
	//
	// For server requests, the URL is parsed from the URI
	// supplied on the Request-Line as stored in RequestURI.  For
	// most requests, fields other than Path and RawQuery will be
	// empty. (See RFC 7230, Section 5.3)
	URL *url.URL

	// The protocol version for incoming server requests.
	//
	// For client requests, these fields are ignored. The HTTP
	// client code always uses either HTTP/1.1 or HTTP/2.
	// See the docs on Transport for details.
	Proto      string // "HTTP/1.0"
	ProtoMajor int    // 1
	ProtoMinor int    // 0

	// Header contains the request header fields either received
	// by the server or to be sent by the client.
	//
	// If a server received a request with header lines,
	//
	//	Host: example.com
	//	accept-encoding: gzip, deflate
	//	Accept-Language: en-us
	//	fOO: Bar
	//	foo: two
	//
	// then
	//
	//	Header = map[string][]string{
	//		"Accept-Encoding": {"gzip, deflate"},
	//		"Accept-Language": {"en-us"},
	//		"Foo": {"Bar", "two"},
	//	}
	//
	// For incoming requests, the Host header is promoted to the
	// Request.Host field and removed from the Header map.
	//
	// HTTP defines that header names are case-insensitive. The
	// request parser implements this by using CanonicalHeaderKey,
	// making the first character and any characters following a
	// hyphen uppercase and the rest lowercase.
	//
	// For client requests, certain headers such as Content-Length
	// and Connection are automatically written when needed and
	// values in Header may be ignored. See the documentation
	// for the Request.Write method.
	// in this file, header is of request
	RHeader http.Header

	// Body is the request's body.
	//
	// For client requests, a nil body means the request has no
	// body, such as a GET request. The HTTP Client's Transport
	// is responsible for calling the Close method.
	//
	// For server requests, the Request Body is always non-nil
	// but will return EOF immediately when no body is present.
	// The Server will close the request body. The ServeHTTP
	// Handler does not need to.
	Body io.ReadCloser

	// ContentLength records the length of the associated content.
	// The value -1 indicates that the length is unknown.
	// Values >= 0 indicate that the given number of bytes may
	// be read from Body.
	ContentLength int64

	// TransferEncoding lists the transfer encodings from outermost to
	// innermost. An empty list denotes the "identity" encoding.
	// TransferEncoding can usually be ignored; chunked encoding is
	// automatically added and removed as necessary when sending and
	// receiving requests.
	TransferEncoding []string

	// For server requests, Host specifies the host on which the
	// URL is sought. For HTTP/1 (per RFC 7230, section 5.4), this
	// is either the value of the "Host" header or the host name
	// given in the URL itself. For HTTP/2, it is the value of the
	// ":authority" pseudo-header field.
	// It may be of the form "host:port". For international domain
	// names, Host may be in Punycode or Unicode form. Use
	// golang.org/x/net/idna to convert it to either format if
	// needed.
	// To prevent DNS rebinding attacks, server Handlers should
	// validate that the Host header has a value for which the
	// Handler considers itself authoritative. The included
	// ServeMux supports patterns registered to particular host
	// names and thus protects its registered Handlers.
	Host string

	// paramters set in the URI string, such as URI: /?name=jone&name=can&pass=0011
	// so the QueryForm can be that
	//map[string][]string{
	// 	"name": []string{"jone", "can"},
	// 	"pass": []string{"0011"},
	// }
	QueryForm url.Values

	// PostForm contains the parsed form data from PATCH, POST
	// or PUT body parameters.
	//
	// This field is only available after ParseForm is called.
	// The HTTP client ignores PostForm and uses Body instead.
	PostForm url.Values

	// MultipartForm is the parsed multipart form, including file uploads.
	// This field is only available after ParseMultipartForm is called.
	// The HTTP client ignores MultipartForm and uses Body instead.
	MultipartForm *multipart.Form

	// RemoteAddr allows HTTP servers and other software to record
	// the network address that sent the request, usually for
	// logging. This field is not filled in by ReadRequest and
	// has no defined format. The HTTP server in this package
	// sets RemoteAddr to an "IP:port" address before invoking a
	// handler.
	// This field is ignored by the HTTP client.
	RemoteAddr string

	// RequestURI is the unmodified request-target of the
	// Request-Line (RFC 7230, Section 3.1.1) as sent by the client
	// to a server. Usually the URL field should be used instead.
	// It is an error to set this field in an HTTP client request.
	RequestURI string

	// TLS allows HTTP servers and other software to record
	// information about the TLS connection on which the request
	// was received. This field is not filled in by ReadRequest.
	// The HTTP server in this package sets the field for
	// TLS-enabled connections before invoking a handler;
	// otherwise it leaves the field nil.
	// This field is ignored by the HTTP client.
	TLS *tls.ConnectionState

	// Cookies from request
	RCookies []*http.Cookie
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w:                w,
		r:                r,
		Param:            newParam(),
		Method:           r.Method,
		URL:              r.URL,
		Proto:            r.Proto,
		ProtoMajor:       r.ProtoMajor,
		ProtoMinor:       r.ProtoMinor,
		RHeader:          r.Header,
		Body:             r.Body,
		ContentLength:    r.ContentLength,
		TransferEncoding: r.TransferEncoding,
		Host:             r.Host,
		QueryForm:        r.URL.Query(),
		//PostForm
		//MultipartForm
		RemoteAddr: r.RemoteAddr,
		RequestURI: r.RequestURI,
		TLS:        r.TLS,
		RCookies:   r.Cookies(),
	}
}

// Request return the http.Request of context
func (c *Context) Request() *http.Request {
	return c.r
}

// ResponseWriter return the w of context
func (c *Context) ResponseWriter() http.ResponseWriter {
	return c.w
}

// UserAgent return the useragent in request's cookie
func (c *Context) UserAgent() string {
	return c.r.UserAgent()
}

// Referer return the router information of request
func (c *Context) Referer() string {
	return c.r.Referer()
}

// RCookie return the cookie in request
func (c *Context) RCookie(name string) (*http.Cookie, error) {
	return c.r.Cookie(name)
}

// ParsePostForm parse the form of POST method
func (c *Context) ParsePostForm() (url.Values, error) {
	err := c.r.ParseForm()
	if err != nil {
		return nil, err
	}
	c.PostForm = c.r.PostForm
	return c.PostForm, nil
}

// ParseMultipartForm parse the files or something uploaded
func (c *Context) ParseMultipartForm(maxMemory int64) (*multipart.Form, error) {
	err := c.r.ParseForm()
	if err != nil {
		return nil, err
	}

	err = c.r.ParseMultipartForm(maxMemory)
	if err != nil {
		return nil, err
	}
	c.MultipartForm = c.r.MultipartForm
	return c.MultipartForm, nil
}
func newParam() map[string]string {
	return make(map[string]string)
}

// Write implements Writer, so the context can be used directly to write to Http.ResponseWriter
func (c *Context) Write(buf []byte) (int, error) {
	return c.w.Write(buf)
}

// Read implements Reader, so the context can be used directly to read the body of Http.Request
func (c *Context) Read(buf []byte) (int, error) {
	return c.r.Body.Read(buf)
}

// Header return the Header for responseWriter
func (c *Context) Header() http.Header {
	return c.w.Header()
}

// WriteHeader write the status code in header
func (c *Context) WriteHeader(code int) {
	c.w.WriteHeader(code)
}

// String is used to write string to ResponseWriter
func (c *Context) String(statusCode int, msg string, a ...interface{}) error {
	c.WriteHeader(statusCode)
	_, err := fmt.Fprintf(c, msg, a...)
	return err
}
