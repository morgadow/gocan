package bus

type MessageID uint32
type MessageType uint8
type BusState uint8

const (
	DataFrame             MessageType = iota // CAN Data frame with 11-bit or 29-bit
	RemoteFrame           MessageType = iota // CAN Remote-Transfer-Request Frame (Broadcast from one node to all others)
	ErrorFrame            MessageType = iota // Signals observed errors on CAN network
	OverloadFrame         MessageType = iota // Used to add extra delay between data or remote frames (rarely used)
	FDBitRateSwitchFrame  MessageType = iota // Only CANFD
	FDErrorStateIndicator MessageType = iota // Only CANFD
)

const (
	ACTIVE  BusState = iota // read and write
	PASSIVE BusState = iota // listen-only
)

// Message CAN message for standard CAN and CAN FD
type Message struct {
	ID         MessageID
	Data       []byte
	TimeStamp  uint64      // only set when receiving message
	Type       MessageType // only set when receiving message
	DLC        uint8       // only set when receiving message
	Channel    string      // only set when receiving message
	IsExtended bool        // only set when receiving message
	IsFD       bool        // only set when receiving message
}

// Bus Interface for all main CANBus functionality. Lower device interfaces may support more functionality
type Bus interface {
	Send(*Message) error
	Recv(timeout uint32) (*Message, error)
	StatusIsOkay() (bool, error)
	Status() (uint32, error)
	State() BusState
	ReadBuffer(limit uint16) ([]Message, error)
	SetFilter(fromID MessageID, toID MessageID, mode uint8) error
	Reset() error
	Shutdown() error
}

// Config CANBus config ready to be read from any json file
type Config struct {
	BusType          string   `json:"busType"`
	Channel          string   `json:"channel"`
	BaudRate         uint32   `json:"baudRate"`
	BusState         BusState `json:"BusState"`
	IsFD             bool     `json:"isFD"`
	FDParameter      string   `json:"Parameter"`       // those parameter are only used if given bus is a FD bus
	RecvStatusFrames bool     `json:"RecvStatusrames"` // Defines behaviour on receiving status frames; if set to true, status frames can be received on Recv() call, not all CAN drivers might support status frames
	LogStatusFrames  bool     `json:"LogStatusFrames"` // Defines behaviour on receiving status frames; if set to true, status frames are logged, not all CAN drivers might support status frames
}
