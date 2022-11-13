package main

import (
    "fmt"
    "io"
    "log"
    "net"
    "net/http"
    "time"
)

func main() {
    server := &http.Server {
        Addr: ":8080",
        Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if r.Method == http.MethodConnect {
                handleTunneling(w, r)
            } else {
                handleHttp(w, r)
            }
        }),
    }

    log.Fatal(server.ListenAndServe())
}

func handleTunneling(w http.ResponseWriter, r *http.Request) {
    destConn, err := net.DialTimeout("tcp", r.Host, 10 * time.Second)
    if err != nil {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
        return
    }

    fmt.Println("Proxying (Tunneling) request to", r.Host)

    w.WriteHeader(http.StatusOK) 
    hijack, ok := w.(http.Hijacker)
    if !ok {
        http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
        return
    }

    clientConn, _, err := hijack.Hijack()
    if err != nil {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
    }
    go transfer(destConn, clientConn)
    go transfer(clientConn, destConn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
    defer destination.Close()
    defer source.Close()
    io.Copy(destination, source)
}

func handleHttp(w http.ResponseWriter, r *http.Request) {
    resp, err := http.DefaultTransport.RoundTrip(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
        return
    }
    defer resp.Body.Close()

    fmt.Println("Proxying (General) request to", r.Host)

    copyHeader(w.Header(), resp.Header)
    w.WriteHeader(resp.StatusCode)
    io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
    for k, vs := range src {
        for _, v := range vs {
            dst.Add(k, v)
        }
    }
}