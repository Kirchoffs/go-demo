package main

import (
    "net/http"
    "simple-server/api"
)

func main() {
    srv := api.NewServer()
    if err := http.ListenAndServe(":8080", srv); err == nil {
        panic(err)
    }
}