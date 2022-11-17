package pcan

import (
	"errors"
	"fmt"
	"syscall"
	"time"

	"github.com/morgadow/gocan"
	log "github.com/sirupsen/logrus"
)

const StandardLanguage = LanguageNeutral // selected language for error messages
const PositionStateInDataStatusFrame = 3 // position of TPCANStatus inside a StatusFrame message
var bootTimeEpoch uint64 = 0             // todo implement this to be not zero but the correct epoch time (datasheet or maybe python implementation)
var hasEvents = true                     // indicates if WaitForSingleObject can be used to reduce CPU load while waiting for messages

// PCANFDParameterList Parameters which must be set for Initializing a PCANBusFD
var PCANFDParameterList = [...]string{"f_clock", "data_brp", "data_sjw", "data_tseg1", "data_tseg2", "nom_brp", "nom_sjw", "nom_tseg1", "nom_tseg2"} // todo implement CANFD and provide examples

// PCANFDOptionalParameterList Parameters which can be set for Initializing a PCANBusFD
var PCANFDOptionalParameterList = [...]string{"data_ssp_offset", "nom_sam"}

// AvailableBaudRates All available standard baud rates for PCAN Channels (defining custom is theoretically possible)
// Defined this way to improve config file suppport
var AvailableBaudRates = map[uint32]TPCANBaudrate{
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
var ErrInvalidBaudRate = errors.New("invalid baudrate selected, choose one of pcan.AvailableBaudRates")

// AvailableChannels All available PCAN Channels
// Defined this way to improve config file suppport
var AvailableChannels = map[string]TPCANHandle{
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
var ErrInvalidChannel = errors.New("invalid channel selected, choose one of pcan.AvailableChannels")

// CAN_FD_DLC List of valid data lengths for a CAN FD message
var CAN_FD_DLC = [...]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 12, 16, 20, 24, 32, 48, 64}

// pcanBus PCAN BusIf capable of sending and reading CAN messages
type pcanBus struct {
	Config    gocan.Config
	Handle    TPCANHandle
	Bitrate   TPCANBaudrate  // only set if not a FD channel
	BitrateFD TPCANBitrateFD // only set if a FD channel
	HWType    TPCANType
	IOPort    uint32
	Interrupt uint16
	recvEvent syscall.Handle
}

// NewPCANBus Convenient method for creating and initiating a pcanBus with multiple default parameters channel
func NewPCANBus(config *gocan.Config) (gocan.Bus, error) {

	var baud TPCANBaudrate
	var bitrateFD TPCANBitrateFD
	var handle TPCANHandle
	var ok = false

	// load api if not done already
	if !APILoaded {
		err := LoadAPI()
		if err != nil {
			return nil, err
		}
	}

	// create bus
	if config.IsFD {
		return nil, errors.New("CANFD not implemented error")
	} else {

		if handle, ok = AvailableChannels[config.Channel]; !ok {
			return nil, ErrInvalidChannel
		}
		if baud, ok = AvailableBaudRates[config.BaudRate]; !ok {
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

		// set bus state parameter
		switch config.BusState {
		case gocan.ACTIVE:
			err = newBus.SetValue(PCAN_LISTEN_ONLY, uint32(PCAN_PARAMETER_OFF))
		case gocan.PASSIVE:
			err = newBus.SetValue(PCAN_LISTEN_ONLY, uint32(PCAN_PARAMETER_ON))
		default:
			err = fmt.Errorf("invalid busstate: %v", err)
		}

		return newBus, err
	}
}

// Initialize Initializes PCANStandardBus channel
func (p *pcanBus) Initialize() error {

	var ret = PCAN_ERROR_UNKNOWN
	var err error = nil

	if p.Config.IsFD {
		// TODO Initialize and test CAN FD channel : ret, err = InitializeFD(p.Channel, )
		return errors.New("pcan FD not implemented")
	} else {
		ret, err = Initialize(p.Handle, p.Bitrate, p.HWType, p.IOPort, p.Interrupt)
		err = evalRetval(ret, err)
	}

	if err != nil {
		return err
	}

	// prepare WaitForSingleObject implementation when waiting for CAN messages (currently only windows support)
	p.recvEvent = 0
	if hasEvents {
		modkernel32, errLoad := syscall.LoadLibrary("kernel32.dll")
		procCreateEvent, errOpen := syscall.GetProcAddress(modkernel32, "CreateEventW")
		if errLoad == nil && errOpen == nil && procCreateEvent != 0 {
			r0, _, errno := syscall.Syscall(procCreateEvent, 0, 0, 0, 0)
			if errno == 0 && r0 != 0 && syscall.Handle(r0) != syscall.InvalidHandle {
				p.recvEvent = syscall.Handle(r0)
				retVal, errVal := SetValue(p.Handle, PCAN_RECEIVE_EVENT, uint32(r0))
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

// Recv Returns message from PCANStandardBus
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
				default:
					return nil, errWait
				}
			} else {
				// timeout handling
				if time.Now().UnixNano()/int64(time.Millisecond) > endTime {
					return nil, err
				}
				time.Sleep(1 * time.Millisecond)
			}
		}
	}

	return msg, err
}

// recvSingleMessage Reads single message from PCAN CAN gocan.
func (p *pcanBus) recvSingleMessage() (TPCANStatus, *gocan.Message, error) {

	var newMsg gocan.Message
	var msgType gocan.MessageType
	var ret = PCAN_ERROR_UNKNOWN
	var msg TPCANMessage
	var msgFD TPCANMessageFD
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
		if msg.MsgType == PCAN_MESSAGE_STATUS && p.Config.LogStatusFrames {
			log.Warning(getFormattedError(TPCANStatus(msg.Data[PositionStateInDataStatusFrame])))
		}
		if err != nil || ret == PCAN_ERROR_QRCVEMPTY || (msg.MsgType == PCAN_MESSAGE_STATUS && !p.Config.RecvStatusFrames) {
			return ret, nil, err
		}
		rxDLC = msgFD.DLC
		rxMsgType = msgFD.MsgType
		rxTimeStamp = bootTimeEpoch + uint64(timestampFD)/(1000.0*1000.0)
		rxData = msgFD.Data[:getLengthFromDLC(rxDLC)] // only return the suggested message length, even if full message is held in buffer with up to 64 byte
	} else {
		ret, msg, timestamp, err = Read(p.Handle)
		if msg.MsgType == PCAN_MESSAGE_STATUS && p.Config.LogStatusFrames {
			log.Warning(getFormattedError(TPCANStatus(msg.Data[PositionStateInDataStatusFrame])))
		}
		if err != nil || ret == PCAN_ERROR_QRCVEMPTY || (msg.MsgType == PCAN_MESSAGE_STATUS && !p.Config.RecvStatusFrames) {
			return ret, nil, err
		}

		rxDLC = msg.DLC
		rxMsgType = msg.MsgType
		rxTimeStamp = bootTimeEpoch + ((uint64(timestamp.Micros) + 1000*uint64(timestamp.Millis) + uint64(0x100000000)*1000*uint64(timestamp.MillisOverflow)) / (1000.0 * 1000.0))
		rxData = msg.Data[:getLengthFromDLC(rxDLC)] // only return the suggested message length, even if full message is held in buffer with 8 byte
	}

	// determine message frame type
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

// Send Sends message over PCAN channel
func (p *pcanBus) Send(msg *gocan.Message) error {

	var ret = PCAN_ERROR_UNKNOWN
	var err error = nil

	// CAN FD copy to CAN FD message and send
	if p.Config.IsFD {
		var pcanMsg = TPCANMessageFD{
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
		var pcanMsg = TPCANMessage{
			ID:      TPCANMsgID(msg.ID),
			MsgType: msgType,
			DLC:     getDLCFromLength(len(msg.Data)),
		}
		copy(pcanMsg.Data[:], msg.Data)

		ret, err = Write(p.Handle, pcanMsg)
	}

	return evalRetval(ret, err)
}

// StatusIsOkay Convenient function to check PCANStandardBus for PCAN_ERROR_OK Status, other bus errors are ignored
func (p *pcanBus) StatusIsOkay() (bool, error) {
	ret, err := GetStatus(p.Handle)
	return ret == PCAN_ERROR_OK, err
}

// Status Returns Status of PCANStandardBus channel
func (p *pcanBus) Status() (uint32, error) {
	state, err := GetStatus(p.Handle)
	return uint32(state), evalRetval(state, err)
}

// State Returns State of PCANStandardBus channel
func (p *pcanBus) State() gocan.BusState {
	return p.Config.BusState
}

// ReadBuffer Reads from device buffer until it has no more messages stored with an optional message limit
func (p *pcanBus) ReadBuffer(limit uint16) ([]gocan.Message, error) {

	var ret = PCAN_ERROR_UNKNOWN
	var msg *gocan.Message
	var err error = nil
	var msgs []gocan.Message

	// read until buffer empty is returned
	for {
		ret, msg, err = p.recvSingleMessage()
		if ret == PCAN_ERROR_QRCVEMPTY || len(msgs) >= int(limit) {
			return msgs, err
		}
		if msg != nil {
			msgs = append(msgs, *msg)
		}
	}
}

// GetValue Retrieves a TPCANParameter value from channel or device
func (p *pcanBus) GetValue(param TPCANParameter) (uint32, error) {
	state, val, err := GetValue(p.Handle, param)
	return val, evalRetval(state, err)
}

// SetValue Configures a TPCANParameter from channel or device
func (p *pcanBus) SetValue(param TPCANParameter, value uint32) error {
	state, err := SetValue(p.Handle, param, value)
	return evalRetval(state, err)
}

// FilterMessages Apply message filter to PCANStandardBus channel
func (p *pcanBus) SetFilter(fromID gocan.MessageID, toID gocan.MessageID, mode uint8) error {
	state, err := FilterMessages(p.Handle, TPCANMsgID(fromID), TPCANMsgID(toID), TPCANMode(mode))
	return evalRetval(state, err)
}

// Reset Resets PCANStandardBus in order to gain PCAN_ERROR_OK Status
func (p *pcanBus) Reset() error {
	state, err := Reset(p.Handle)
	return evalRetval(state, err)
}

// Shutdown Shuts channel down and closes connection
func (p *pcanBus) Shutdown() error {

	state, err := Uninitialize(p.Handle)
	if p.recvEvent != 0 { // close handle
		_ = syscall.CloseHandle(p.recvEvent)
	}
	return evalRetval(state, err)
}

// FlashLED Turn on or off flashing of the device's LED for physical identification purposes
func (p *pcanBus) FlashLED(state bool) error {
	val := 0
	if state {
		val = 1
	}
	return p.SetValue(PCAN_CHANNEL_IDENTIFYING, uint32(val))
}

// getFormattedError Uses the API function to get string description for any TPCANStatus
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

// evalRetval: function handles the API return bus state and error and evaluates final error
// status: PCAN Status returned from api
// err: error returned from api and syscall
func evalRetval(status TPCANStatus, err error) error {
	if err != nil {
		return err
	}
	if status != PCAN_ERROR_OK && (status != PCAN_ERROR_BUSLIGHT && status != PCAN_ERROR_BUSHEAVY) {
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
		return MaxLengthDataCANFDMessage
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
