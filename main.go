package main

import (
    "net/http"
    "fmt"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type PlayerInput struct {
    x float32
    y float32
    clicked bool
    disconnected bool
    conn *websocket.Conn
    outputChan chan StateOutput
}

type StateOutput struct {
    message string
}

func connect(inputChan chan PlayerInput) func(http.ResponseWriter,*http.Request) {
    return func(writer http.ResponseWriter, req *http.Request) {
        fmt.Println("hit connect")
        connection, err := upgrader.Upgrade(writer, req, nil)
        if err != nil {
            fmt.Println("no connect")
            return
        }
        fmt.Println("no connect error")
        outputChan := make(chan StateOutput)
        go func() {
            fmt.Println("starting reader routine")
            initInput := PlayerInput{x: 0, y: 0, clicked: false, disconnected: false,
                                     conn: connection, outputChan: outputChan}
            inputChan <- initInput

            for {
                _, message, err := connection.ReadMessage()
                fmt.Println("got message")
                if err != nil {
                    fmt.Println("bad connection")
                    disconnectInput := PlayerInput{conn: connection, disconnected: true}
                    connection.Close()
                    inputChan <- disconnectInput
                    break
                }
                fmt.Printf("%s", message)
            }
        }()

        go func() {
            fmt.Println("starting writer routine")
            for {
                output, ok := (<- outputChan)

                if !ok {
                    fmt.Println("exiting write routine")
                    break
                }

                err := connection.WriteMessage(0, []byte(output.message))
                if err != nil {
                    connection.Close()
                }
            }
        }()
    }
}

func gameLoop(inputChan chan PlayerInput) {
    connMap := make(map[*websocket.Conn]chan StateOutput)

    for {
        chanLen := len(inputChan)
        for i := 0; i < chanLen; i++ {
            select {
            case input := <- inputChan:
                fmt.Println("Got some input")
                if _, exists := connMap[input.conn]; !exists {
                    fmt.Println("new connection!")
                    connMap[input.conn] = input.outputChan
                }
            default:
                break
            }
        }
    }
}

func main() {
    inputChan := make(chan PlayerInput)

    go gameLoop(inputChan)

    fmt.Println("Whaddup")
    http.HandleFunc("/connect", connect(inputChan))
    fmt.Println(http.ListenAndServe("localhost:8000", nil))
}