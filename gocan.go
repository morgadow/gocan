package gocan

type MessageID uint32
type MessageType uint8
type BusState uint8
type ChannelCondition uint8

// Different types of available frames on a CAN bus connection
const (
	DataFrame             MessageType = iota // CAN Data frame with 11-bit or 29-bit
	RemoteFrame           MessageType = iota // CAN Remote-Transfer-Request Frame (Broadcast from one node to all others)
	ErrorFrame            MessageType = iota // Signals observed errors on CAN network
	OverloadFrame         MessageType = iota // Used to add extra delay between data or remote frames (rarely used)
	FDBitRateSwitchFrame  MessageType = iota // Only CANFD
	FDErrorStateIndicator MessageType = iota // Only CANFD
)

// Bus state
const (
	ACTIVE  BusState = iota // read and write
	PASSIVE BusState = iota // listen-only
)

// Condition for a single handle
const (
	Available   ChannelCondition = iota // Channel is available for a connection
	Occupied    ChannelCondition = iota // Channel is already occupied by a connection, a connection may be possible depending on interface type
	Unavailable ChannelCondition = iota // Channel is not available or not connected
	Invalid     ChannelCondition = iota // Invalid state or not able to retrieve state for this interface
)

// CAN message for standard CAN and CAN FD
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

// Interface for all main CANBus functionality. Lower device interfaces may support more functionality
type Bus interface {
	Send(*Message) error                                          // Send a single message on the CAN bus
	Recv(timeout int) (*Message, error)                           // Receive single message from CAN bus with timeout in [ms], a timeout below zero is treated as no timeout
	StatusIsOkay() (bool, error)                                  // Check function if the connection state is okay
	Status() (uint32, error)                                      // Returns the CAN status code, which can differ between different devices
	State() BusState                                              // Returns the bus state (ACTIVE or PASSIVE)
	ReadBuffer(limit uint16) ([]Message, error)                   // Empties the internal CAN hardware message buffer is device supports this feature with a maximum message count
	SetFilter(fromID MessageID, toID MessageID, mode uint8) error // Set a message id filter on hardware if supported by device
	Reset() error                                                 // Reset rx and tx buffer, does not reset hardware
	Shutdown() error                                              // Disconnect from device
	ChannelCondition() (ChannelCondition, error)                  // Returns channel condition
}

// CANBus config ready to be read from any json file
type Config struct {
	BusType          string   `json:"busType"`
	Channel          string   `json:"channel"`
	BaudRate         uint32   `json:"baudRate"`
	BusState         BusState `json:"BusState"`
	IsFD             bool     `json:"isFD"`
	FDParameter      string   `json:"Parameter"`        // those parameter are only used if given bus is a FD bus
	RecvStatusFrames bool     `json:"RecvStatusFrames"` // If set to true, status frames can be received on Recv() call
	RecvRTRFrames    bool     `json:"RecvRTRFrames"`    // If set to true, remote transmission frames can be received on Recv() call
	RecvErrorFrames  bool     `json:"RecvErrorFrames"`  // If set to true, error frames can be received on Recv() call
	RecvEchoFrames   bool     `json:"RecvEchoFrames"`   // If set to true, echo frames can be received on Recv() call
}
