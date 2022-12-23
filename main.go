package main

import (
    "os"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
)

// NewProxy takes target host and creates a reverse proxy
// NewProxy 拿到 targetHost 后，创建一个反向代理
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
    url, err := url.Parse(targetHost)
    if err != nil {
        return nil, err
    }

    return httputil.NewSingleHostReverseProxy(url), nil
}

// ProxyRequestHandler handles the http request using proxy
// ProxyRequestHandler 使用 proxy 处理请求
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        proxy.ServeHTTP(w, r)
    }
}

func main() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("/home/write-data/writefreely/static/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("/home/write-data/writefreely/static/js"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("/home/write-data/writefreely/static/img"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("/home/write-data/writefreely/static/fonts"))))

    // initialize a reverse proxy and pass the actual backend server url here
    // 初始化反向代理并传入真正后端服务的地址
    proxy, err := NewProxy("http://localhost:8080")
    if err != nil {
        panic(err)
    }

    // handle all requests to your server using the proxy
    // 使用 proxy 处理所有请求到你的服务
    port := os.Getenv("PORT")
    http.HandleFunc("/", ProxyRequestHandler(proxy))
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

