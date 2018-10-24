package main

import (
    "time"
    "net/http"
    "fmt"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type PlayerInput struct {
    x            float32
    y            float32
    clicked      bool
    disconnected bool
    conn         *websocket.Conn
    outputChan   chan StateOutput
}

type BodyState struct {
    X     float32 `json:"x"`
    Y     float32 `json:"y"`
    Theta float32 `json:"theta"`
}

type LegState struct {
    X     float32         `json:"x"`
    Y     float32         `json:"y"`
    Theta float32         `json:"theta"`
    Owner bool            `json:"owner"`
    Conn  *websocket.Conn `json:"-"`
}

type StateOutput struct {
    Body BodyState  `json:"body"`
    Legs []LegState `json:"legs"`
}

func connect(inputChan chan PlayerInput) func(http.ResponseWriter,*http.Request) {
    return func(writer http.ResponseWriter, req *http.Request) {
        fmt.Println("hit connect")
        connection, err := upgrader.Upgrade(writer, req, nil)
        if err != nil {
            fmt.Printf("connection failed: %s\n", err)
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
                if err != nil {
                    fmt.Println("read connection terminated")
                    disconnectInput := PlayerInput{conn: connection, disconnected: true}
                    connection.Close()
                    inputChan <- disconnectInput
                    break
                }
                fmt.Printf("%s\n", message)
            }
        }()

        go func() {
            fmt.Println("starting writer routine")
            for {
                output, ok := (<- outputChan)

                if !ok {
                    fmt.Println("game loop closed write connection")
                    break
                }

                err := connection.WriteJSON(output)
                if err != nil {
                    fmt.Printf("write connection terminated: %s\n", err)
                    connection.Close()
                }
            }
        }()
    }
}

func gameLoop(inputChan chan PlayerInput) {
    connMap := make(map[*websocket.Conn]chan StateOutput)
    fmt.Println("Starting game loop")
    for {
        chanLen := len(inputChan)
        if chanLen == 0 {
            chanLen = 1
        }

        for i := 0; i < chanLen; i++ {
            select {
            case input := <- inputChan:
                fmt.Println("Got some input")
                if _, exists := connMap[input.conn]; !exists {
                    fmt.Println("new connection!")
                    connMap[input.conn] = input.outputChan
                }
                if input.disconnected {
                    close(connMap[input.conn])
                    delete(connMap, input.conn)
                }
            default:
                break
            }
        }

        for _, outChan := range connMap {
            outChan <- StateOutput{Body: BodyState{X: 1, Y: 1, Theta: 1}, Legs: []LegState{LegState{X: 2, Y: 2, Theta: 2}, LegState{X: 3, Y: 3, Theta: 3}}}
        }
        time.Sleep(time.Millisecond * 30)
    }
}

func main() {
    inputChan := make(chan PlayerInput)
    go gameLoop(inputChan)

    fmt.Println("Whaddup")
    http.HandleFunc("/connect", connect(inputChan))
    fmt.Println(http.ListenAndServe("localhost:8000", nil))
}