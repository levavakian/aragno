package ecs

// EntityId alias for entity ids
type EntityId uint64

// EntityManager manages the creation of new entities
type EntityManager struct {
	LastId EntityId
}

// NewId generates a new unique id from an entity manager
func (em *EntityManager) NewId() EntityId {
	em.LastId = em.LastId + 1
	return em.LastId
}

// PipeType alias for pipe types
type PipeType uint64

// System interface for types that run updates at every frame and can be ordered by priority
type System interface {
	Update(dt float32)
	Priority() uint64
}

// Message used for communicating between systems
type Message struct {
	Pipe               PipeType
	ActingEntityID     EntityId
	SecondaryEntityIds []EntityId
	Data               interface{}
}

// Messager pub/sub for systems
type Messenger struct {
	Callbacks map[PipeType][]func(message *Message)
}

// AddCallback adds a callback for a given PipeType
func (messenger *Messenger) AddCallback(msgType PipeType, fn func(*Message)) {
	messenger.Callbacks[msgType] = append(messenger.Callbacks[msgType], fn)
}

// SendMessage sends a message to the PipeType specified in msg
func (messenger *Messenger) SendMessage(msg *Message) {
	for _, fn := range messenger.Callbacks[msg.Pipe] {
		fn(msg)
	}
}
