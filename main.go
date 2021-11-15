package main

import (
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // 接続されるクライアント
var broadcast = make(chan Message)           // メッセージブロードキャストチャネル

// アップグレーダ
var upgrader = websocket.Upgrader{}

// メッセージ用構造体
type Message struct {
    Email    string `json:"email"`
    Username string `json:"username"`
    Message  string `json:"message"`
}

func main() {
    // websockerへのルーティングを紐づけ
    http.HandleFunc("/ws", handleConnections)
    go handleMessages()
    // サーバーをlocalhostのポート8000で立ち上げる
    log.Println("http server started on :8000")
    err := http.ListenAndServe(":8000", nil)
    // エラーがあった場合ロギングする
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
    // 送られてきたGETリクエストをwebsocketにアップグレード
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Fatal(err)
    }
    // 関数が終わった際に必ずwebsocketnのコネクションを閉じる
    defer ws.Close()

    // クライアントを新しく登録
    clients[ws] = true

    for {
        var msg Message
        // 新しいメッセージをJSONとして読み込みMessageオブジェクトにマッピングする
        err := ws.ReadJSON(&msg)
        if err != nil {
            log.Printf("error: %v", err)
            delete(clients, ws)
            break
        }
        // 新しく受信されたメッセージをブロードキャストチャネルに送る
        broadcast <- msg
    }
}

func handleMessages() {
    for {
        // ブロードキャストチャネルから次のメッセージを受け取る
        msg := <-broadcast
        // 現在接続しているクライアント全てにメッセージを送信する
        for client := range clients {
            err := client.WriteJSON(msg)
            if err != nil {
                log.Printf("error: %v", err)
                client.Close()
                delete(clients, client)
            }
        }
    }
}