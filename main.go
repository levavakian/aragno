package main

import (
    "time"
    "net/http"
    "fmt"
    "sync"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type PlayerInput struct {
    X            float32
    Y            float32
    Clicked      bool
    Disconnected bool
    Conn         *websocket.Conn
    OutputChan   chan StateOutput
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
    Body      BodyState  `json:"body"`
    Legs      []LegState `json:"legs"`
    PlayerIdx int        `json:"player_idx"`
}

type CollectedInputs struct {
    InputChan chan PlayerInput
    Inputs    map[*websocket.Conn]PlayerInput
    CloseSig  chan struct{}
    mux       sync.Mutex
}

func (collector *CollectedInputs) Add(input PlayerInput) {
    collector.mux.Lock()
    defer collector.mux.Unlock()
    collector.Inputs[input.Conn] = input
}

func (collector *CollectedInputs) Pop() map[*websocket.Conn]PlayerInput {
    collector.mux.Lock()
    defer collector.mux.Unlock()

    inputs := collector.Inputs
    collector.Inputs = make(map[*websocket.Conn]PlayerInput)
    return inputs
}

func (collecter *CollectedInputs) Collect() {
    for {
        select {
        case input := <- collecter.InputChan:
            collecter.Add(input)
        case _, ok := <- collecter.CloseSig:
            if !ok {
                fmt.Println("Shutting down input collecter")
                break;
            }
        default:
            // Do nothing
        }
    }
}

type OutputDispatch struct {
    OutputChans map[*websocket.Conn]chan StateOutput
}

func (dispatch *OutputDispatch) ConsumeInputs(inputs map[*websocket.Conn]PlayerInput) {
    for inConn, inInput := range inputs {
        fmt.Printf("Consuming input from %s\n", &inConn)

        _, exists := dispatch.OutputChans[inConn]
        
        if !exists {
            dispatch.OutputChans[inConn] = inInput.OutputChan
        }

        if inInput.Disconnected {
            close(dispatch.OutputChans[inConn])
            delete(dispatch.OutputChans, inConn)
        }
    }
}

func (dispatch *OutputDispatch) Dispatch(state StateOutput) {
    for _, outChan := range dispatch.OutputChans {
        outChan <- state
    }
}

func gameLoop(collector *CollectedInputs) {
    fmt.Println("Starting game loop")

    // Set up output dispatcher
    outputDispatch := OutputDispatch{OutputChans: make(map[*websocket.Conn]chan StateOutput)}

    // Game loop
    for {
        // Collect and consume inputs
        collectedInputs := collector.Pop()
        outputDispatch.ConsumeInputs(collectedInputs)

        // Do game stuff
        state := StateOutput{Body: BodyState{X: 1, Y: 1, Theta: 1}, Legs: []LegState{LegState{X: 2, Y: 2, Theta: 2}, LegState{X: 3, Y: 3, Theta: 3}}}

        // Output game state to websockets
        outputDispatch.Dispatch(state)

        // Server run rate at ~30Hz
        time.Sleep(time.Millisecond * 33)
    }
}

func connect(inputChan chan PlayerInput) func(http.ResponseWriter,*http.Request) {
    return func(writer http.ResponseWriter, req *http.Request) {
        connection, err := upgrader.Upgrade(writer, req, nil)
        if err != nil {
            fmt.Printf("Connection failed: %s\n", err)
            return
        }

        outputChan := make(chan StateOutput)
        go func() {
            fmt.Printf("Starting reader for %s", &connection)
            initInput := PlayerInput{X: 0, Y: 0, Clicked: false, Disconnected: false,
                                     Conn: connection, OutputChan: outputChan}
            inputChan <- initInput

            for {
                _, message, err := connection.ReadMessage()
                if err != nil {
                    fmt.Printf("Read connection closed for %s", &connection)
                    disconnectInput := PlayerInput{Conn: connection, Disconnected: true}
                    connection.Close()
                    inputChan <- disconnectInput
                    break
                }
                fmt.Printf("%s\n", message)
            }
        }()

        go func() {
            fmt.Printf("Starting writer for %s", &connection)
            for {
                output, ok := (<- outputChan)

                if !ok {
                    fmt.Printf("Close requested for %s\n", &connection)
                    break
                }

                err := connection.WriteJSON(output)
                if err != nil {
                    fmt.Printf("Write connection failed with %s for %s\n", err, &connection)
                    connection.Close()
                }
            }
        }()
    }
}

func main() {
    // Set up collecter to grab outputs
    collector := &CollectedInputs{InputChan: make(chan PlayerInput), Inputs: make(map[*websocket.Conn]PlayerInput), CloseSig: make(chan struct{})}
    go collector.Collect()

    go gameLoop(collector)

    http.HandleFunc("/connect", connect(collector.InputChan))
    fmt.Println(http.ListenAndServe("localhost:8000", nil))
}