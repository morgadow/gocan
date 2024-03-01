package pcan

import (
	"errors"
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"github.com/morgadow/gocan"
)

// defines and singleton values
const StandardLanguage = LanguageNeutral // selected language for error messages
const PositionStateInDataStatusFrame = 3 // position of TPCANStatus inside a StatusFrame message
var bootTimeEpoch uint64 = 0             // Message epoch time / TODO implement this to be not zero but the correct epoch time (datasheet or maybe python implementation)
var hasEvents = true                     // indicates if WaitForSingleObject can be used to reduce CPU load while waiting for messages /

// errors
var (
	ErrInvalidChannel  = errors.New("invalid channel selected")
	ErrInvalidBaudRate = errors.New("invalid baudrate selected")
)

// All available standard baud rates for PCAN Channels (defining custom is theoretically possible)
// Defined this way to improve config file suppport
var IntToBaudrate = map[uint32]TPCANBaudrate{
	1000000: PCAN_BAUD_1M,
	800000:  PCAN_BAUD_800K,
	500000:  PCAN_BAUD_500K,
	250000:  PCAN_BAUD_250K,
	125000:  PCAN_BAUD_125K,
	100000:  PCAN_BAUD_100K,
	95000:   PCAN_BAUD_95K,
	83000:   PCAN_BAUD_83K,
	50000:   PCAN_BAUD_50K,
	47000:   PCAN_BAUD_47K,
	33000:   PCAN_BAUD_33K,
	20000:   PCAN_BAUD_20K,
	10000:   PCAN_BAUD_10K,
	5000:    PCAN_BAUD_5K,
}
var BaudrateToInt = map[TPCANBaudrate]uint32{
	PCAN_BAUD_1M:   1000000,
	PCAN_BAUD_800K: 800000,
	PCAN_BAUD_500K: 500000,
	PCAN_BAUD_250K: 250000,
	PCAN_BAUD_125K: 125000,
	PCAN_BAUD_100K: 100000,
	PCAN_BAUD_95K:  95000,
	PCAN_BAUD_83K:  83000,
	PCAN_BAUD_50K:  50000,
	PCAN_BAUD_47K:  47000,
	PCAN_BAUD_33K:  33000,
	PCAN_BAUD_20K:  20000,
	PCAN_BAUD_10K:  10000,
	PCAN_BAUD_5K:   5000,
}

// All available PCAN Channels
// Defined this way to improve config file suppport
var StringToChannel = map[string]TPCANHandle{
	"PCAN_NONEBUS": PCAN_NONEBUS,

	"PCAN_ISABUS1": PCAN_ISABUS1,
	"PCAN_ISABUS2": PCAN_ISABUS2,
	"PCAN_ISABUS3": PCAN_ISABUS3,
	"PCAN_ISABUS4": PCAN_ISABUS4,
	"PCAN_ISABUS5": PCAN_ISABUS5,
	"PCAN_ISABUS6": PCAN_ISABUS6,
	"PCAN_ISABUS7": PCAN_ISABUS7,
	"PCAN_ISABUS8": PCAN_ISABUS8,

	"PCAN_DNGBUS1": PCAN_DNGBUS1,

	"PCAN_PCIBUS1":  PCAN_PCIBUS1,
	"PCAN_PCIBUS2":  PCAN_PCIBUS2,
	"PCAN_PCIBUS3":  PCAN_PCIBUS3,
	"PCAN_PCIBUS4":  PCAN_PCIBUS4,
	"PCAN_PCIBUS5":  PCAN_PCIBUS5,
	"PCAN_PCIBUS6":  PCAN_PCIBUS6,
	"PCAN_PCIBUS7":  PCAN_PCIBUS7,
	"PCAN_PCIBUS8":  PCAN_PCIBUS8,
	"PCAN_PCIBUS9":  PCAN_PCIBUS9,
	"PCAN_PCIBUS10": PCAN_PCIBUS10,
	"PCAN_PCIBUS11": PCAN_PCIBUS11,
	"PCAN_PCIBUS12": PCAN_PCIBUS12,
	"PCAN_PCIBUS13": PCAN_PCIBUS13,
	"PCAN_PCIBUS14": PCAN_PCIBUS14,
	"PCAN_PCIBUS15": PCAN_PCIBUS15,
	"PCAN_PCIBUS16": PCAN_PCIBUS16,

	"PCAN_USBBUS1":  PCAN_USBBUS1,
	"PCAN_USBBUS2":  PCAN_USBBUS2,
	"PCAN_USBBUS3":  PCAN_USBBUS3,
	"PCAN_USBBUS4":  PCAN_USBBUS4,
	"PCAN_USBBUS5":  PCAN_USBBUS5,
	"PCAN_USBBUS6":  PCAN_USBBUS6,
	"PCAN_USBBUS7":  PCAN_USBBUS7,
	"PCAN_USBBUS8":  PCAN_USBBUS8,
	"PCAN_USBBUS9":  PCAN_USBBUS9,
	"PCAN_USBBUS10": PCAN_USBBUS10,
	"PCAN_USBBUS11": PCAN_USBBUS11,
	"PCAN_USBBUS12": PCAN_USBBUS12,
	"PCAN_USBBUS13": PCAN_USBBUS13,
	"PCAN_USBBUS14": PCAN_USBBUS14,
	"PCAN_USBBUS15": PCAN_USBBUS15,
	"PCAN_USBBUS16": PCAN_USBBUS16,

	"PCAN_PCCBUS1": PCAN_PCCBUS1,
	"PCAN_PCCBUS2": PCAN_PCCBUS2,

	"PCAN_LANBUS1":  PCAN_LANBUS1,
	"PCAN_LANBUS2":  PCAN_LANBUS2,
	"PCAN_LANBUS3":  PCAN_LANBUS3,
	"PCAN_LANBUS4":  PCAN_LANBUS4,
	"PCAN_LANBUS5":  PCAN_LANBUS5,
	"PCAN_LANBUS6":  PCAN_LANBUS6,
	"PCAN_LANBUS7":  PCAN_LANBUS7,
	"PCAN_LANBUS8":  PCAN_LANBUS8,
	"PCAN_LANBUS9":  PCAN_LANBUS9,
	"PCAN_LANBUS10": PCAN_LANBUS10,
	"PCAN_LANBUS11": PCAN_LANBUS11,
	"PCAN_LANBUS12": PCAN_LANBUS12,
	"PCAN_LANBUS13": PCAN_LANBUS13,
	"PCAN_LANBUS14": PCAN_LANBUS14,
	"PCAN_LANBUS15": PCAN_LANBUS15,
	"PCAN_LANBUS16": PCAN_LANBUS16,
}
var ChannelToString = map[TPCANHandle]string{
	PCAN_NONEBUS:  "PCAN_NONEBUS",
	PCAN_ISABUS1:  "PCAN_ISABUS1",
	PCAN_ISABUS2:  "PCAN_ISABUS2",
	PCAN_ISABUS3:  "PCAN_ISABUS3",
	PCAN_ISABUS4:  "PCAN_ISABUS4",
	PCAN_ISABUS5:  "PCAN_ISABUS5",
	PCAN_ISABUS6:  "PCAN_ISABUS6",
	PCAN_ISABUS7:  "PCAN_ISABUS7",
	PCAN_ISABUS8:  "PCAN_ISABUS8",
	PCAN_DNGBUS1:  "PCAN_DNGBUS1",
	PCAN_PCIBUS1:  "PCAN_PCIBUS1",
	PCAN_PCIBUS2:  "PCAN_PCIBUS2",
	PCAN_PCIBUS3:  "PCAN_PCIBUS3",
	PCAN_PCIBUS4:  "PCAN_PCIBUS4",
	PCAN_PCIBUS5:  "PCAN_PCIBUS5",
	PCAN_PCIBUS6:  "PCAN_PCIBUS6",
	PCAN_PCIBUS7:  "PCAN_PCIBUS7",
	PCAN_PCIBUS8:  "PCAN_PCIBUS8",
	PCAN_PCIBUS9:  "PCAN_PCIBUS9",
	PCAN_PCIBUS10: "PCAN_PCIBUS10",
	PCAN_PCIBUS11: "PCAN_PCIBUS11",
	PCAN_PCIBUS12: "PCAN_PCIBUS12",
	PCAN_PCIBUS13: "PCAN_PCIBUS13",
	PCAN_PCIBUS14: "PCAN_PCIBUS14",
	PCAN_PCIBUS15: "PCAN_PCIBUS15",
	PCAN_PCIBUS16: "PCAN_PCIBUS16",
	PCAN_USBBUS1:  "PCAN_USBBUS1",
	PCAN_USBBUS2:  "PCAN_USBBUS2",
	PCAN_USBBUS3:  "PCAN_USBBUS3",
	PCAN_USBBUS4:  "PCAN_USBBUS4",
	PCAN_USBBUS5:  "PCAN_USBBUS5",
	PCAN_USBBUS6:  "PCAN_USBBUS6",
	PCAN_USBBUS7:  "PCAN_USBBUS7",
	PCAN_USBBUS8:  "PCAN_USBBUS8",
	PCAN_USBBUS9:  "PCAN_USBBUS9",
	PCAN_USBBUS10: "PCAN_USBBUS10",
	PCAN_USBBUS11: "PCAN_USBBUS11",
	PCAN_USBBUS12: "PCAN_USBBUS12",
	PCAN_USBBUS13: "PCAN_USBBUS13",
	PCAN_USBBUS14: "PCAN_USBBUS14",
	PCAN_USBBUS15: "PCAN_USBBUS15",
	PCAN_USBBUS16: "PCAN_USBBUS16",
	PCAN_PCCBUS1:  "PCAN_PCCBUS1",
	PCAN_PCCBUS2:  "PCAN_PCCBUS2",
	PCAN_LANBUS1:  "PCAN_LANBUS1",
	PCAN_LANBUS2:  "PCAN_LANBUS2",
	PCAN_LANBUS3:  "PCAN_LANBUS3",
	PCAN_LANBUS4:  "PCAN_LANBUS4",
	PCAN_LANBUS5:  "PCAN_LANBUS5",
	PCAN_LANBUS6:  "PCAN_LANBUS6",
	PCAN_LANBUS7:  "PCAN_LANBUS7",
	PCAN_LANBUS8:  "PCAN_LANBUS8",
	PCAN_LANBUS9:  "PCAN_LANBUS9",
	PCAN_LANBUS10: "PCAN_LANBUS10",
	PCAN_LANBUS11: "PCAN_LANBUS11",
	PCAN_LANBUS12: "PCAN_LANBUS12",
	PCAN_LANBUS13: "PCAN_LANBUS13",
	PCAN_LANBUS14: "PCAN_LANBUS14",
	PCAN_LANBUS15: "PCAN_LANBUS15",
	PCAN_LANBUS16: "PCAN_LANBUS16",
}

// List of valid data lengths for a CAN FD message
var CAN_FD_DLC = [...]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 12, 16, 20, 24, 32, 48, 64}

// pcanBus PCAN BusIf capable of sending and reading CAN messages
type pcanBus struct {
	Config    gocan.Config
	Handle    TPCANHandle
	Bitrate   TPCANBaudrate  // only set if not a FD channel
	BitrateFD TPCANBitrateFD // only set if a FD channel
	HWType    TPCANType      // only for non plug´n´play devices and currently not used
	IOPort    uint32         // only for non plug´n´play devices and currently not used
	Interrupt uint16         // only for non plug´n´play devices and currently not used
	recvEvent syscall.Handle
}

// Convenient method for creating and initiating a pcanBus with multiple default parameters channel
func NewPCANBus(config *gocan.Config) (gocan.Bus, error) {

	var baud TPCANBaudrate
	var bitrateFD TPCANBitrateFD
	var handle TPCANHandle
	var ok = false

	// load api if not done already
	if !apiLoaded {
		err := LoadAPI()
		if err != nil {
			return nil, err
		}
	}

	// create bus
	if config.IsFD {
		return nil, errors.New("CANFD not implemented error")
	} else {

		if handle, ok = StringToChannel[config.Channel]; !ok {
			return nil, ErrInvalidChannel
		}
		if baud, ok = IntToBaudrate[config.BaudRate]; !ok {
			return nil, ErrInvalidBaudRate
		}

		newBus := &pcanBus{
			Config:    *config,
			Handle:    handle,
			Bitrate:   baud,
			BitrateFD: bitrateFD,
			HWType:    PCAN_TYPE_ISA, // default value, might not work for all types of handles
			IOPort:    0x02A0,        // default value, might not work for all types of handles
			Interrupt: 11,            // default value, might not work for all types of handles
		}
		err := newBus.Initialize()
		if err != nil {
			return nil, err
		}

		// set bus parameter depending on config
		var states = map[gocan.BusState]TPCANParameterValue{gocan.ACTIVE: PCAN_PARAMETER_OFF, gocan.PASSIVE: PCAN_PARAMETER_ON}
		SetParameter(newBus.Handle, PCAN_LISTEN_ONLY, states[config.BusState])

		// setting for receiving functions
		var conv = map[bool]TPCANParameterValue{false: PCAN_PARAMETER_OFF, true: PCAN_PARAMETER_ON}
		SetParameter(newBus.Handle, PCAN_ALLOW_STATUS_FRAMES, conv[config.RecvStatusFrames])
		SetParameter(newBus.Handle, PCAN_ALLOW_RTR_FRAMES, conv[config.RecvRTRFrames])
		SetParameter(newBus.Handle, PCAN_ALLOW_ERROR_FRAMES, conv[config.RecvErrorFrames])
		SetParameter(newBus.Handle, PCAN_ALLOW_ECHO_FRAMES, conv[config.RecvEchoFrames])

		return newBus, err
	}
}

// Initializes PCANStandardBus channel
func (p *pcanBus) Initialize() error {

	var ret = PCAN_ERROR_UNKNOWN
	var err error = nil

	if p.Config.IsFD {
		return errors.New("pcan FD not implemented")
	} else {
		ret, err = Initialize(p.Handle, p.Bitrate, p.HWType, p.IOPort, p.Interrupt)
		err = evalRetval(ret, err)
		if err != nil {
			return err
		}
	}

	// prepare WaitForSingleObject implementation when waiting for CAN messages (currently only windows support)
	p.recvEvent = 0
	if hasEvents {
		modkernel32, errLoad := syscall.LoadLibrary("kernel32.dll")
		procCreateEvent, errOpen := syscall.GetProcAddress(modkernel32, "CreateEventW")
		if errLoad == nil && errOpen == nil && procCreateEvent != 0 {
			r0, _, errno := syscall.SyscallN(procCreateEvent)
			if errno == 0 && r0 != 0 && syscall.Handle(r0) != syscall.InvalidHandle {
				p.recvEvent = syscall.Handle(r0)
				retVal, errVal := SetParameter(p.Handle, PCAN_RECEIVE_EVENT, TPCANParameterValue(r0))
				if retVal != PCAN_ERROR_OK || errVal != nil {
					hasEvents = false
					_ = syscall.CloseHandle(p.recvEvent)
					p.recvEvent = 0
				}
			}
		}
		// just for safety
		if p.recvEvent == 0 || p.recvEvent == syscall.InvalidHandle {
			hasEvents = false
		}
	}

	return nil
}

// Returns message from PCANStandardBus
// timeout: Timeout for receiving message from CAN bus in milliseconds (if set below zero, no timeout is set)
func (p *pcanBus) Recv(timeout int) (*gocan.Message, error) {

	var ret = PCAN_ERROR_UNKNOWN
	var msg *gocan.Message = nil
	var err error = nil

	// timeout handling: a negative timeout sets timeout to infinity
	if timeout < 0 {
		timeout = syscall.INFINITE
	}
	var timeoutU32 = uint32(timeout)
	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	endTime := startTime + int64(timeout)

	// receive message
	for msg == nil {
		ret, msg, err = p.recvSingleMessage()
		if ret == PCAN_ERROR_QRCVEMPTY {
			if hasEvents {
				val, errWait := syscall.WaitForSingleObject(p.recvEvent, timeoutU32)
				switch val {
				case syscall.WAIT_OBJECT_0:
					break
				case syscall.WAIT_FAILED:
					return nil, errWait
				case syscall.WAIT_TIMEOUT:
					return nil, errWait
				default:
					return nil, errWait
				}
			} else {
				// timeout handling
				if time.Now().UnixNano()/int64(time.Millisecond) > endTime {
					return nil, err
				}
				time.Sleep(250 * time.Microsecond)
			}
		}
	}

	return msg, err
}

// Reads single message from PCAN CAN gocan.
func (p *pcanBus) recvSingleMessage() (TPCANStatus, *gocan.Message, error) {

	var newMsg gocan.Message
	var msgType gocan.MessageType
	var ret = PCAN_ERROR_UNKNOWN
	var msg TPCANMsg
	var msgFD TPCANMsgFD
	var timestamp TPCANTimestamp
	var timestampFD TPCANTimestampFD
	var rxData []byte              // buffer for uniform handling FD or std messages
	var rxMsgType TPCANMessageType // buffer for uniform handling FD or std messages
	var rxTimeStamp uint64         // buffer for uniform handling FD or std messages
	var rxDLC uint8                // buffer for uniform handling FD or std messages
	var err error = nil

	// receive single message, already converted to gocan.Message
	if p.Config.IsFD {
		ret, msgFD, timestampFD, err = ReadFD(p.Handle)
		if err != nil || ret == PCAN_ERROR_QRCVEMPTY {
			return ret, nil, err
		}
		rxDLC = msgFD.DLC
		rxMsgType = msgFD.MsgType
		rxTimeStamp = bootTimeEpoch + uint64(timestampFD)/(1000.0*1000.0)
		rxData = msgFD.Data[:getLengthFromDLC(rxDLC)] // only return the suggested message length, even if full message is held in buffer with up to 64 byte
	} else {
		ret, msg, timestamp, err = Read(p.Handle)

		if err != nil || ret == PCAN_ERROR_QRCVEMPTY {
			return ret, nil, err
		}

		rxDLC = msg.DLC
		rxMsgType = msg.MsgType
		rxTimeStamp = bootTimeEpoch + ((uint64(timestamp.Micros) + 1000*uint64(timestamp.Millis) + uint64(0x100000000)*1000*uint64(timestamp.MillisOverflow)) / (1000.0 * 1000.0))
		rxData = msg.Data[:getLengthFromDLC(rxDLC)] // only return the suggested message length, even if full message is held in buffer with 8 byte
	}

	// determine message frame type	// TODO rewrite switch?
	if rxMsgType == PCAN_MESSAGE_STANDARD || rxMsgType == PCAN_MESSAGE_EXTENDED || rxMsgType == PCAN_MESSAGE_FD {
		msgType = gocan.DataFrame
	} else if rxMsgType == PCAN_MESSAGE_RTR {
		msgType = gocan.RemoteFrame
	} else if rxMsgType == PCAN_MESSAGE_ERRFRAME || rxMsgType == PCAN_MESSAGE_STATUS {
		msgType = gocan.ErrorFrame
	} else if rxMsgType == PCAN_MESSAGE_BRS {
		msgType = gocan.FDBitRateSwitchFrame
	} else if rxMsgType == PCAN_MESSAGE_ESI {
		msgType = gocan.FDErrorStateIndicator
	}

	// save message data
	newMsg = gocan.Message{
		ID:         gocan.MessageID(msg.ID),
		TimeStamp:  rxTimeStamp,
		Type:       msgType,
		Data:       rxData,
		DLC:        rxDLC,
		Channel:    p.Config.Channel,
		IsFD:       rxMsgType == PCAN_MESSAGE_FD || rxMsgType == PCAN_MESSAGE_ESI || rxMsgType == PCAN_MESSAGE_BRS,
		IsExtended: rxMsgType == PCAN_MESSAGE_EXTENDED,
	}

	return ret, &newMsg, nil
}

// Sends message over PCAN channel
func (p *pcanBus) Send(msg *gocan.Message) error {

	var ret = PCAN_ERROR_UNKNOWN
	var err error = nil

	// CAN FD copy to CAN FD message and send
	if p.Config.IsFD {
		var pcanMsg = TPCANMsgFD{
			ID:      TPCANMsgID(msg.ID),
			MsgType: TPCANMessageType(msg.Type),
			DLC:     uint8(len(msg.Data)),
		}
		copy(pcanMsg.Data[:], msg.Data)

		ret, err = WriteFD(p.Handle, pcanMsg)

		// Standard CAN copy to CAN message and send
	} else {
		msgType := PCAN_MESSAGE_STANDARD
		if msg.IsExtended {
			msgType = PCAN_MESSAGE_EXTENDED
		}
		var pcanMsg = TPCANMsg{
			ID:      TPCANMsgID(msg.ID),
			MsgType: msgType,
			DLC:     getDLCFromLength(len(msg.Data)),
		}
		copy(pcanMsg.Data[:], msg.Data)

		ret, err = Write(p.Handle, pcanMsg)
	}

	return evalRetval(ret, err)
}

// Convenient function to check PCANStandardBus for PCAN_ERROR_OK Status, other bus errors are ignored
func (p *pcanBus) StatusIsOkay() (bool, error) {
	ret, err := GetStatus(p.Handle)
	return ret == PCAN_ERROR_OK, err
}

// Returns Status of PCANStandardBus channel
func (p *pcanBus) Status() (uint32, error) {
	state, err := GetStatus(p.Handle)
	return uint32(state), evalRetval(state, err)
}

// Returns State of PCANStandardBus channel
func (p *pcanBus) State() gocan.BusState {
	return p.Config.BusState
}

// Reads from device buffer until it has no more messages stored with an optional message limit
// If limit is set to zero, no limit will will be used
func (p *pcanBus) ReadBuffer(limit uint16) ([]gocan.Message, error) {

	var ret = PCAN_ERROR_UNKNOWN
	var msg *gocan.Message
	var err error = nil
	var msgs []gocan.Message

	// read until buffer empty is returned
	for {
		ret, msg, err = p.recvSingleMessage()
		if ret == PCAN_ERROR_QRCVEMPTY {
			return msgs, err
		}
		if msg != nil {
			msgs = append(msgs, *msg)
			if limit != 0 && len(msgs) >= int(limit) {
				return msgs, err
			}
		}
	}
}

// Retrieves a TPCANParameter value from channel or device (only work for simple parameters)
func (p *pcanBus) GetParameter(param TPCANParameter) (TPCANParameterValue, error) {
	state, val, err := GetParameter(p.Handle, param)
	return val, evalRetval(state, err)
}

// Configures a TPCANParameter from channel or device (only work for simple parameters)
func (p *pcanBus) SetParameter(param TPCANParameter, val TPCANParameterValue) error {
	state, err := SetParameter(p.Handle, param, val)
	return evalRetval(state, err)
}

// Retrieves a TPCANParameter value from channel or device
func (p *pcanBus) GetValue(param TPCANParameter) (uint32, error) {
	var buf uint32
	state, err := GetValue(p.Handle, param, unsafe.Pointer(&buf), uint32(unsafe.Sizeof(buf)))
	return buf, evalRetval(state, err)
}

// Configures a TPCANParameter from channel or device
func (p *pcanBus) SetValue(param TPCANParameter, val TPCANParameterValue) error {
	state, err := SetParameter(p.Handle, param, val)
	return evalRetval(state, err)
}

// Apply message filter to PCANStandardBus channel
func (p *pcanBus) SetFilter(fromID gocan.MessageID, toID gocan.MessageID, mode uint8) error {
	state, err := FilterMessages(p.Handle, TPCANMsgID(fromID), TPCANMsgID(toID), TPCANMode(mode))
	return evalRetval(state, err)
}

// Resets PCANStandardBus in order to gain PCAN_ERROR_OK Status
func (p *pcanBus) Reset() error {
	state, err := Reset(p.Handle)
	return evalRetval(state, err)
}

// Shuts channel down and closes connection
func (p *pcanBus) Shutdown() error {

	state, err := Uninitialize(p.Handle)
	if p.recvEvent != 0 { // close handle
		_ = syscall.CloseHandle(p.recvEvent)
	}
	return evalRetval(state, err)
}

// Turn on or off flashing of the device's LED for physical identification purposes
func (p *pcanBus) SetLEDState(ledState bool) error {
	val := PCAN_PARAMETER_OFF
	if ledState {
		val = PCAN_PARAMETER_ON
	}
	state, err := SetParameter(p.Handle, PCAN_CHANNEL_IDENTIFYING, val)
	return evalRetval(state, err)
}

// Returns the channel condition as a level for availablity
func (p *pcanBus) ChannelCondition2() (gocan.ChannelCondition, error) {
	var cond gocan.ChannelCondition = gocan.Invalid
	state, val, err := GetParameter(p.Handle, PCAN_CHANNEL_CONDITION)

	condition := TPCANCHannelCondition(val)
	if (condition & PCAN_CHANNEL_AVAILABLE) == PCAN_CHANNEL_AVAILABLE {
		cond = gocan.Available
	} else if (condition & PCAN_CHANNEL_OCCUPIED) == PCAN_CHANNEL_OCCUPIED {
		cond = gocan.Occupied
	} else if (condition & PCAN_CHANNEL_PCANVIEW) == PCAN_CHANNEL_PCANVIEW {
		cond = gocan.Occupied
	} else if (condition & PCAN_CHANNEL_UNAVAILABLE) == PCAN_CHANNEL_UNAVAILABLE {
		cond = gocan.Unavailable
	}

	return cond, evalRetval(state, err)
}

// Returns the channel condition as a level for availablity
func (p *pcanBus) ChannelCondition() (gocan.ChannelCondition, error) {
	var buf uint32
	var cond gocan.ChannelCondition = gocan.Invalid
	state, err := GetValue(p.Handle, PCAN_CHANNEL_CONDITION, unsafe.Pointer(&buf), uint32(unsafe.Sizeof(buf)))

	if (buf & uint32(PCAN_CHANNEL_AVAILABLE)) == uint32(PCAN_CHANNEL_AVAILABLE) {
		cond = gocan.Available
	} else if (buf & uint32(PCAN_CHANNEL_OCCUPIED)) == uint32(PCAN_CHANNEL_OCCUPIED) {
		cond = gocan.Occupied
	} else if (buf & uint32(PCAN_CHANNEL_PCANVIEW)) == uint32(PCAN_CHANNEL_PCANVIEW) {
		cond = gocan.Occupied
	} else if (buf & uint32(PCAN_CHANNEL_UNAVAILABLE)) == uint32(PCAN_CHANNEL_UNAVAILABLE) {
		cond = gocan.Unavailable
	}

	return cond, evalRetval(state, err)
}

// Starts recording a trace on given path with a max file size in MB
// maxFileSize: trace file is splitted in files with this maximum size of file in MB; set to zero to have a infinite large trace file (max is 100 MB)
// Note: A trace file only gets filled if the Recv() function is called!
func (p *pcanBus) TraceStart(filePath string, maxFileSize uint32) error {

	if maxFileSize > MAX_TRACE_FILE_SIZE_ACCEPTED {
		return fmt.Errorf("maximum size of a trace file is %v MB", MAX_TRACE_FILE_SIZE_ACCEPTED)
	}

	// configure trace configuration (only file size is set, the other options are always used)
	cfg := TRACE_FILE_DATE | TRACE_FILE_TIME | TRACE_FILE_OVERWRITE
	if maxFileSize > 0 {
		cfg |= TRACE_FILE_SEGMENTED
	} else {
		cfg |= TRACE_FILE_SINGLE
	}
	state, err := SetParameter(p.Handle, PCAN_TRACE_CONFIGURE, TPCANParameterValue(cfg))
	if err != nil || state != PCAN_ERROR_OK {
		return evalRetval(state, err)
	}
	if maxFileSize > 0 {
		state, err := SetValue(p.Handle, PCAN_TRACE_SIZE, unsafe.Pointer(&maxFileSize), 4)
		if err != nil || state != PCAN_ERROR_OK {
			return evalRetval(state, err)
		}
	}

	// configure trace file path
	if len(filePath) > MAX_LENGHT_STRING_BUFFER {
		return fmt.Errorf("filepath exceeds max length of %v", MAX_LENGHT_STRING_BUFFER)
	}

	// convert path to a fixed buffer size as pcan wants it that way
	var buffer [MAX_LENGHT_STRING_BUFFER]byte
	for i := range filePath {
		buffer[i] = byte(filePath[i])
	}
	state, err = SetValue(p.Handle, PCAN_TRACE_LOCATION, unsafe.Pointer(&buffer), unsafe.Sizeof(buffer))
	if err != nil || state != PCAN_ERROR_OK {
		return evalRetval(state, err)
	}

	// start tracing
	state, err = SetParameter(p.Handle, PCAN_TRACE_STATUS, PCAN_PARAMETER_ON)
	return evalRetval(state, err)
}

// Stops recording currently running trace
func (p *pcanBus) TraceStop() error {
	state, err := SetParameter(p.Handle, PCAN_TRACE_STATUS, PCAN_PARAMETER_OFF)
	return evalRetval(state, err)
}

// Uninitializes all PCAN Channels initialized by CAN_Initialize
func ShutdownAllHandles() error {
	state, err := Uninitialize(PCAN_NONEBUS)
	return evalRetval(state, err)
}

// Returns list of all existing PCAN channels on a system in a single call, regardless of their current availability
func AttachedChannels() ([]TPCANHandle, error) {
	posChannels := [...]TPCANHandle{PCAN_USBBUS1, PCAN_USBBUS2, PCAN_USBBUS3, PCAN_USBBUS4,
		PCAN_USBBUS5, PCAN_USBBUS6, PCAN_USBBUS7, PCAN_USBBUS8,
		PCAN_USBBUS9, PCAN_USBBUS10, PCAN_USBBUS11, PCAN_USBBUS12,
		PCAN_USBBUS13, PCAN_USBBUS14, PCAN_USBBUS15, PCAN_USBBUS16}
	attachedChannels := []TPCANHandle{}

	// iterate through channels and check for every channel if available
	for i := range posChannels {
		state, cond, err := GetParameter(posChannels[i], PCAN_CHANNEL_CONDITION)
		if state != PCAN_ERROR_OK || err != nil {
			return nil, err
		}
		if cond == TPCANParameterValue(PCAN_CHANNEL_AVAILABLE) ||
			cond == TPCANParameterValue(PCAN_CHANNEL_OCCUPIED) ||
			cond == TPCANParameterValue(PCAN_CHANNEL_PCANVIEW) {
			attachedChannels = append(attachedChannels, posChannels[i])
		}
	}

	return attachedChannels, nil
}

// Returns list of all existing PCAN channels on a system in a single call, regardless of their current availability
// TODO This function is not working correctly, as the given information does not matched connected hardware
func AttachedChannels_Extended() ([]TPCANChannelInformation, error) {
	count, err := AttachedChannelsCount()
	if err != nil || count == 0 { // size calculation not possible with a slice len of 0
		return nil, err
	}

	buf := make([]TPCANChannelInformation, count)
	size := uintptr(len(buf)) * unsafe.Sizeof(buf[0])
	state, err := GetValue(PCAN_NONEBUS, PCAN_ATTACHED_CHANNELS, unsafe.Pointer(&buf[0]), uint32(size))

	return buf, evalRetval(state, err)
}

// Returns list of all existing PCAN channels on a system in a single call, regardless of their current availability
func AttachedChannelsNames() ([]string, error) {
	channels, err := AttachedChannels()
	if err != nil { // size calculation not possible with a slice len of 0
		return nil, err
	}

	names := make([]string, len(channels))
	for i := range channels {
		names[i] = ChannelToString[channels[i]]
	}

	return names, nil
}

// Gets information about all existing PCAN channels on a system in a single call, regardless of their current availability.
func AttachedChannelsCount() (uint32, error) {
	var channelCount uint32
	ret, err := GetValue(PCAN_NONEBUS, PCAN_ATTACHED_CHANNELS_COUNT, unsafe.Pointer(&channelCount), uint32(unsafe.Sizeof(channelCount)))
	if err != nil {
		return channelCount, err
	}
	return channelCount, getFormattedError(ret)
}

// Uses the API function to get string description for any TPCANStatus
// status: PCAN Status returned from api
func getFormattedError(status TPCANStatus) error {

	// if no error, do not convert to error but return nil
	if status == PCAN_ERROR_OK {
		return nil
	}

	ret, buffer, err := GetErrorText(status, StandardLanguage)
	if err != nil {
		return err
	}
	if ret != PCAN_ERROR_OK {
		return fmt.Errorf("could not retrieve error text for pcan error %v", status)
	}

	var numBytes = 0
	for i, b := range buffer {
		if b == 0 {
			numBytes = i
			break
		}
	}
	return errors.New(string(buffer[:numBytes]))
}

// helper function handles the API return bus state and error and evaluates final error
// status: PCAN Status returned from api
// err: error returned from api and syscall
func evalRetval(status TPCANStatus, err error) error {
	if err != nil {
		return err
	}
	if status != PCAN_ERROR_OK {
		err = getFormattedError(status)
	}
	return err
}

// Convert a CAN DLC value into the actual data length of the CAN/CAN-FD frame.
// dlc: A value between 0 and 15 (CAN and FD DLC range)
// returns: The message length represented by the DLC
func getLengthFromDLC(dlc uint8) int {

	if dlc <= 8 {
		return int(dlc)
	} else if dlc >= 15 {
		return LENGTH_DATA_CANFD_MESSAGE
	} else {
		switch dlc {
		case 9:
			return 12
		case 10:
			return 16
		case 11:
			return 20
		case 12:
			return 24
		case 13:
			return 32
		case 14:
			return 48
		default:
			return int(dlc)
		}
	}
}

// Convert the the actual data length of the CAN/CAN-FD into the frame CAN DLC value.
// length: Length in number of bytes (0-64). Only supported until MaxLengthDataCANFDMessage!
// returns: The message length represented by the DLC
func getDLCFromLength(length int) uint8 {

	if length <= 8 {
		return uint8(length)
	}

	for dlc, nofBytes := range CAN_FD_DLC {
		if nofBytes >= uint8(length) {
			return uint8(dlc)
		}
	}
	return 15 // max DLC possible in a CAN FD message
}
