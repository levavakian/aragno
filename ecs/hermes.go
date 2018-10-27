package ecs

// PipeType alias for pipe types
type PipeType uint64

// Message used for communicating between systems
type Message struct {
	Pipe               PipeType
	EntityId           EntityId
	Targets            []EntityId
	Data               interface{}
}

// Messager pub/sub for systems
type Hermes struct {
	Callbacks map[PipeType][]func(message *Message)
}

func NewHermes() *Hermes {
	return &Hermes{make(map[PipeType][]func(message *Message))}
}

// AddCallback adds a callback for a given PipeType
func (hermes *Hermes) AddCallback(msgType PipeType, fn func(*Message)) {
	hermes.Callbacks[msgType] = append(hermes.Callbacks[msgType], fn)
}

// SendMessage sends a message to the PipeType specified in msg
func (hermes *Hermes) Send(msg *Message) {
	for _, fn := range hermes.Callbacks[msg.Pipe] {
		fn(msg)
	}
}
