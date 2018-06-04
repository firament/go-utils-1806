# GO Lang Request Object

- http2 package, import "golang.org/x/net/http2" directly and use its ConfigureTransport and/or ConfigureServer functions
- Declares
```go
MethodGet     = "GET"
MethodPost    = "POST"
StatusContinue           = 100

type Request struct {
    Proto      string // "HTTP/1.0"
    RemoteAddr string
}
```
- `func (r *Request) Referer() string`
    - Referer returns the referring URL, if sent in the request. 
- `func (srv *Server) SetKeepAlivesEnabled(v bool)`
    - SetKeepAlivesEnabled controls whether HTTP keep-alives are enabled. 


## Cookies
- `func (r *Request) Cookies() []*Cookie`
    - Cookies parses and returns the HTTP cookies sent with the request. 
- `func (r *Request) Cookie(name string) (*Cookie, error)`
    - Cookie returns the named cookie provided in the request or ErrNoCookie if not found. If multiple cookies match the given name, only one cookie will be returned. 
- `var ErrNoCookie = errors.New("http: named cookie not present")`
    - ErrNoCookie is returned by Request's Cookie method when a cookie is not found.
- `func SetCookie(w ResponseWriter, cookie *Cookie)`
    - SetCookie adds a Set-Cookie header to the provided ResponseWriter's headers. The provided cookie must have a valid Name. Invalid cookies may be silently dropped. 
```go
type Cookie struct {
        Name  string
        Value string

        Path       string    // optional
        Domain     string    // optional
        Expires    time.Time // optional
        RawExpires string    // for reading cookies only

        // MaxAge=0 means no 'Max-Age' attribute specified.
        // MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
        // MaxAge>0 means Max-Age attribute present and given in seconds
        MaxAge   int
        Secure   bool
        HttpOnly bool
        Raw      string
        Unparsed []string // Raw text of unparsed attribute-value pairs
}
```


## Headers
- `type Header map[string][]string`
    - A Header represents the key-value pairs in an HTTP header.
- `func CanonicalHeaderKey(s string) string`
    - CanonicalHeaderKey returns the canonical format of the header key s


## FileServer
- `func FileServer(root FileSystem) Handler`
    - FileServer returns a handler that serves HTTP requests with the contents of the file system rooted at root.
    - To use the operating system's file system implementation, use http.Dir:
    - `http.Handle("/", http.FileServer(http.Dir("/tmp")))`
- `func ServeFile(w ResponseWriter, r *Request, name string)`
    - ServeFile replies to the request with the contents of the named file or directory. 

## Context
- `func (r *Request) Context() context.Context`
    - Context returns the request's context. To change the context, use WithContext.
    - The returned context is always non-nil; it defaults to the background context. 
- `func WithValue(parent Context, key, val interface{}) Context`
    - WithValue returns a copy of parent in which the value associated with key is val. 
```go
// use in rev-proxy?
ctx := context.WithValue(r.Context(), "Username", cookie.Value)
next.ServeHTTP(w, r.WithContext(ctx))
```
