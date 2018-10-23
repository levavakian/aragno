package main

import (
    "net/http"
    "fmt"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func connect(writer http.ResponseWriter, req *http.Request) {
	connection, err := upgrader.Upgrade(writer, req, nil)
	if err != nil {
		fmt.Println("no connect")
		return
	}
	defer connection.Close()

	for {
		mt, message, err := connection.ReadMessage()
		if err != nil {
			fmt.Println("no read:", err)
			break
		}

		fmt.Printf("recv: %s", message)
		err = connection.WriteMessage(mt, message)
		if err != nil {
			fmt.Println("no write:", err)
			break
		}
	}
}

func main() {
	fmt.Println("Whaddup")
	http.HandleFunc("/connect", connect)
	fmt.Println(http.ListenAndServe("localhost:8000", nil))
}