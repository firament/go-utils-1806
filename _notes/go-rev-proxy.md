# HTTP/HTTPS revrese proxy

## Notes
- Server Http content
```go
webserver.Handle("/", http.FileServer(http.Dir(httpContentBasePath))
```

## ReverseProxy
- https://golang.org/pkg/net/http/httputil/
- type **ReverseProxy**
    -  ReverseProxy is an HTTP Handler that takes an incoming request and sends it to another server, proxying the response back to the client. 
    ```go
    type ReverseProxy struct {
            // Director must be a function which modifies
            // the request into a new request to be sent
            // using Transport. Its response is then copied
            // back to the original client unmodified.
            // Director must not access the provided Request
            // after returning.
            Director func(*http.Request)

            // The transport used to perform proxy requests.
            // If nil, http.DefaultTransport is used.
            Transport http.RoundTripper

            // FlushInterval specifies the flush interval
            // to flush to the client while copying the
            // response body.
            // If zero, no periodic flushing is done.
            FlushInterval time.Duration

            // ErrorLog specifies an optional logger for errors
            // that occur when attempting to proxy the request.
            // If nil, logging goes to os.Stderr via the log package's
            // standard logger.
            ErrorLog *log.Logger

            // BufferPool optionally specifies a buffer pool to
            // get byte slices for use by io.CopyBuffer when
            // copying HTTP response bodies.
            BufferPool BufferPool

            // ModifyResponse is an optional function that
            // modifies the Response from the backend.
            // If it returns an error, the proxy returns a StatusBadGateway error.
            ModifyResponse func(*http.Response) error
    }
    ```
    - `func NewSingleHostReverseProxy(target *url.URL) *ReverseProxy`
        - NewSingleHostReverseProxy returns a new ReverseProxy that routes URLs to the scheme, host, and base path provided in target. 
        - If the target's path is "/base" and the incoming request was for "/dir", the target request will be for /base/dir. NewSingleHostReverseProxy does not rewrite the Host header. 
        - To rewrite Host headers, use ReverseProxy directly with a custom Director policy.
    - `func (p *ReverseProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request)`

### Director
- https://www.integralist.co.uk/posts/golang-reverse-proxy/
    ```go
    reverseProxy := httputil.NewSingleHostReverseProxy(origin)

    reverseProxy.Director = func(req *http.Request) {
        req.Header.Add("X-Forwarded-Host", req.Host)
        req.Header.Add("X-Origin-Host", origin.Host)
        req.URL.Scheme = origin.Scheme
        req.URL.Host = origin.Host

        req.URL.Path = proxyPath
    }
    ```

### Transport
- httputil.ReverseProxy has a Transport field. You can use it to modify the response
- Now httputil/reverseproxy supports
    ```go
    type ReverseProxy struct {
            ...

            // ModifyResponse is an optional function that
            // modifies the Response from the backend
            // If it returns an error, the proxy returns a StatusBadGateway error.
            ModifyResponse func(*http.Response) error
        }

    // Sample
    func rewriteBody(resp *http.Response) (err error) {
        b, err := ioutil.ReadAll(resp.Body) //Read html
        if err != nil {
            return  err
        }
        err = resp.Body.Close()
        if err != nil {
            return err
        }
        b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1) // replace html
        body := ioutil.NopCloser(bytes.NewReader(b))
        resp.Body = body
        resp.ContentLength = int64(len(b))
        resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
        return nil
    }

    // ...
    target, _ := url.Parse("http://example.com")
    proxy := httputil.NewSingleHostReverseProxy(target)
    proxy.ModifyResponse = rewriteBody
    ```

#### Handler Chaining
```go
func addCORS(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
        handler.ServeHTTP(w, r)
    })
}
```

### Dump Request
- https://golang.org/pkg/net/http/httputil/
    - func DumpRequestOut(req *http.Request, body bool) ([]byte, error)
    - type ReverseProxy struct {
    - func NewSingleHostReverseProxy(target *url.URL) *ReverseProxy
        - NewSingleHostReverseProxy returns a new ReverseProxy that routes URLs to the scheme, host, and base path provided in target. 
        - If the target's path is "/base" and the incoming request was for "/dir", the target request will be for /base/dir. 
        - NewSingleHostReverseProxy does not rewrite the Host header. 
        - To rewrite Host headers, use ReverseProxy directly with a custom Director policy. 
    - func (p *ReverseProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request)
- https://medium.com/ymedialabs-innovation/reverse-proxy-in-go-d26482acbcad
- https://letsencrypt.org/
    - wildcard certificates support
