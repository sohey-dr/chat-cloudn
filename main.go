package main

import (
    "log"
    "net/http"
    "cloudn/handler"
)

func main() {
    http.HandleFunc("/ws", handler.HandleConnections)
    go handler.HandleMessages()
    log.Println("http server started on :8000")
    err := http.ListenAndServe(":8000", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
