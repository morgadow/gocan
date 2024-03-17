package test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/morgadow/gocan/interfaces/pcan"
)

// Note: For this tests, connect a PCAN USB Bus and periodically send four messages:
// - 11bit 500 ms: 0x123 - {1,2,3,4,5,6,7,8}
// - 11bit 500 ms: 0x321 - {8,7,6,5,4,3,2,1}
// - 29bit 500 ms: 0x0000123 - {1,2,3,4,5,6,7,8}
// - 29bit 500 ms: 0x0000321 - {8,7,6,5,4,3,2,1}

var HANDLE_FOR_TESTS = pcan.PCAN_USBBUS1 // this handle is used for all tests

// this function is executed automatically before every test to load pcan API
func init() {
	err := pcan.LoadAPI()
	if err != nil {
		panic(err)
	}
}

func auxInitBasic() {
	state, err := pcan.InitializeBasic(HANDLE_FOR_TESTS, pcan.PCAN_BAUD_500K)
	if err != nil || state != pcan.PCAN_ERROR_OK {
		fmt.Println(fmt.Errorf("Helper function got error or invalid stats on initialize: %v - %w", state, err))
		panic(err)
	}

	// no error frames
	state, err = pcan.SetParameter(HANDLE_FOR_TESTS, pcan.PCAN_ALLOW_ERROR_FRAMES, pcan.PCAN_PARAMETER_OFF)
	if err != nil || state != pcan.PCAN_ERROR_OK {
		fmt.Println(fmt.Errorf("Helper function got error or invalid stats on forbid error frames: %v - %w", state, err))
		panic(err)
	}

	// no status frames
	state, err = pcan.SetParameter(HANDLE_FOR_TESTS, pcan.PCAN_ALLOW_STATUS_FRAMES, pcan.PCAN_PARAMETER_OFF)
	if err != nil || state != pcan.PCAN_ERROR_OK {
		fmt.Println(fmt.Errorf("Helper function got error or invalid stats on forbid status frames: %v - %w", state, err))
		panic(err)
	}

	// allow echo frames so device can read own messages for testing
	state, err = pcan.SetParameter(HANDLE_FOR_TESTS, pcan.PCAN_ALLOW_ECHO_FRAMES, pcan.PCAN_PARAMETER_ON)
	if err != nil || state != pcan.PCAN_ERROR_OK {
		fmt.Println(fmt.Errorf("Helper function got error or invalid stats on activate echo frames: %v - %w", state, err))
		panic(err)
	}
}

func auxInitFD() {
	state, err := pcan.InitializeFD(HANDLE_FOR_TESTS, "")
	if err != nil || state != pcan.PCAN_ERROR_OK {
		fmt.Println(fmt.Errorf("Helper function got error or invalid stats: %v - %w", state, err))
	}
}

// reads until receives a single message with timeout in ms
func auxReadBasic(timeout time.Duration) (pcan.TPCANStatus, *pcan.TPCANMsg, *pcan.TPCANTimestamp, error) {
	var err error = errors.New("no message received")

	start := time.Now()
	for time.Since(start) < timeout {
		state, msg, timestamp, err := pcan.Read(HANDLE_FOR_TESTS) // direct api call
		if state != pcan.PCAN_ERROR_QRCVEMPTY {
			if msg.MsgType == pcan.PCAN_MESSAGE_STANDARD || msg.MsgType == pcan.PCAN_MESSAGE_EXTENDED {
				fmt.Printf("Got Message: %v - ID: %x, Data: %v, DLC: %v, Type: %v\n", timestamp, msg.ID, msg.DLC, msg.Data, msg.MsgType)
				return state, &msg, &timestamp, err
			}
		}
	}
	return pcan.PCAN_ERROR_UNKNOWN, nil, nil, err
}

func auxErrBufToText(buffer [256]byte) string {
	var numBytes = 0
	for i, b := range buffer {
		if b == 0 {
			numBytes = i
			break
		}
	}
	return string(buffer[:numBytes])
}

func TestInitializeBasic(t *testing.T) {
	state, err := pcan.InitializeBasic(HANDLE_FOR_TESTS, pcan.PCAN_BAUD_500K)
	if state != pcan.PCAN_ERROR_OK {

		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestInitialize(t *testing.T) {

	state, err := pcan.Initialize(HANDLE_FOR_TESTS, pcan.PCAN_BAUD_500K, pcan.PCAN_TYPE_ISA, 0x02A0, 11)

	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestUnitialize(t *testing.T) {

	auxInitBasic()
	state, err := pcan.Uninitialize(HANDLE_FOR_TESTS)

	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestReset(t *testing.T) {

	auxInitBasic()
	state, err := pcan.Reset(HANDLE_FOR_TESTS)

	if state != pcan.PCAN_ERROR_OK {

		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestGetStatus(t *testing.T) {
	auxInitBasic()
	state, err := pcan.GetStatus(HANDLE_FOR_TESTS)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

// Note: For this test function another bus member has so send any messages on PCAN_USBBUS1 500 kBaud
func TestRead(t *testing.T) {
	auxInitBasic()
	state, msg, timestamp, err := auxReadBasic(5000 * time.Millisecond)
	if state != pcan.PCAN_ERROR_OK || msg.ID == 0x0 {
		t.Errorf("no message: state: %v, msg: %v, timestamp: %v, err: %v", state, msg, timestamp, err)
	} else {
		fmt.Printf("received message: state: %v, msg: %v, timestamp: %v, err: %v\n", state, msg, timestamp, err)
	}
}

// Note: For this test function another bus member has so send specicicly designed messages on PCAN_USBBUS1 500 kBaud
func TestRead_Specific(t *testing.T) {

	auxInitBasic()

	state, msg, timestamp, err := auxReadBasic(5000 * time.Millisecond)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	if msg.ID != 0x123 && msg.ID != 0x321 && msg.ID != 0x00000123 && msg.ID != 0x00000321 {
		t.Errorf("wrong message id: %v", msg.ID)
	}
	if msg.Data != [8]byte{8, 7, 6, 5, 4, 3, 2, 1} && msg.Data != [8]byte{1, 2, 3, 4, 5, 6, 7, 8} {
		t.Errorf("wrong message data: %v", msg.Data)
	}
	if (timestamp.Micros == 0 || timestamp.Micros == 0xFFFF) && (timestamp.Millis == 0 || timestamp.Millis == 0xFFFFFFFF) && (uint32(timestamp.MillisOverflow) == 0 || timestamp.MillisOverflow == 0xFFFF) {
		t.Errorf("wrong message timestamp: %v", timestamp)
	}
}

// TODO currently no CAN FD really implemented and no FD hardware for testing available
// func TestReadFD(t *testing.T) {
//
// 	auxInitFD(HANDLE_FOR_TESTS)
// 	state, msg, timestamp, err := pcan.ReadFD(HANDLE_FOR_TESTS)
//
// 	fmt.Printf("Got Message: %v - ID: %v, Data: %v, DLC: %v, Type: %v", timestamp, msg.ID, msg.DLC, msg.Data, msg.MsgType)
//
// 	if state != pcan.PCAN_ERROR_OK  {
//
// 		t.Errorf("got non okay status code: %x", state)
// 	}
// 	if err != nil {
// 		t.Errorf("got error: %v", err)
// 	}
// }

func TestWrite(t *testing.T) {

	auxInitBasic()
	msg := pcan.TPCANMsg{ID: 0x123, Data: [8]byte{0, 1, 2, 3, 4, 5, 6, 7}}
	state, err := pcan.Write(HANDLE_FOR_TESTS, msg)

	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

// TODO currently no CAN FD really implemented and no FD hardware for testing available
// func TestWriteFD(t *testing.T) {
//
// 	auxInitFD(HANDLE_FOR_TESTS)
// 	msg := pcan.TPCANMsgFD{ID: 0x123, Data: [64]byte{}}
// 	state, err := pcan.WriteFD(HANDLE_FOR_TESTS, msg)
//
// 	if state != pcan.PCAN_ERROR_OK {
// 		t.Errorf("got non okay status code: %x", state)
// 	}
// 	if err != nil {
// 		t.Errorf("got error: %v", err)
// 	}
// }

// TODO: Function does not check the min and max msg ID
func TestSetFilter(t *testing.T) {
	auxInitBasic()
	var _trans = map[pcan.TPCANFilterValue]string{pcan.PCAN_FILTER_OPEN: "PCAN_FILTER_OPEN", pcan.PCAN_FILTER_CLOSE: "PCAN_FILTER_CLOSE"}

	// check filter is open
	state, val, err := pcan.GetParameter(HANDLE_FOR_TESTS, pcan.PCAN_MESSAGE_FILTER)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if pcan.TPCANFilterValue(val) != pcan.PCAN_FILTER_OPEN {
		t.Errorf("setting was not set correctly to PCAN_FILTER_OPEN: %v", _trans[pcan.TPCANFilterValue(val)])
	}

	// set fitler to 29 bit messages
	state, err = pcan.SetFilter(HANDLE_FOR_TESTS, 0x100, 0x200, pcan.PCAN_MODE_EXTENDED)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	// check filter is set
	state, val, err = pcan.GetParameter(HANDLE_FOR_TESTS, pcan.PCAN_MESSAGE_FILTER)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if pcan.TPCANFilterValue(val) != pcan.PCAN_FILTER_CLOSE {
		t.Errorf("setting was not set correctly to PCAN_FILTER_CLOSE: %v", _trans[pcan.TPCANFilterValue(val)])
	}

	// delete filter
	state, err = pcan.ResetFilter(HANDLE_FOR_TESTS)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	// check filter is open
	state, val, err = pcan.GetParameter(HANDLE_FOR_TESTS, pcan.PCAN_MESSAGE_FILTER)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if pcan.TPCANFilterValue(val) != pcan.PCAN_FILTER_OPEN {
		t.Errorf("setting was not set correctly to PCAN_FILTER_OPEN: %v", _trans[pcan.TPCANFilterValue(val)])
	}

	// set filter for 11 bit messages
	state, err = pcan.SetFilter(HANDLE_FOR_TESTS, 0x100, 0x200, pcan.PCAN_MODE_STANDARD)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	// check filter is set
	state, val, err = pcan.GetParameter(HANDLE_FOR_TESTS, pcan.PCAN_MESSAGE_FILTER)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if pcan.TPCANFilterValue(val) != pcan.PCAN_FILTER_CLOSE {
		t.Errorf("setting was not set correctly to PCAN_FILTER_CLOSE: %v", _trans[pcan.TPCANFilterValue(val)])
	}
}

// Note: For this test to work, another client has to send very specific messages or echo frames are active
func TestSetFilter_Specific(t *testing.T) {

	auxInitBasic()

	// set filter for 29 bit messages
	state, err := pcan.SetFilter(HANDLE_FOR_TESTS, 0x100, 0x200, pcan.PCAN_MODE_STANDARD)
	fmt.Println("Set to: HANDLE_FOR_TESTS, 0x100, 0x200, pcan.PCAN_MODE_STANDARD")
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	// test standard messages and msg id filter 1s
	start := time.Now()
	for time.Since(start) < 1*time.Second {
		_, msg, _, _ := auxReadBasic(5000 * time.Millisecond)
		if msg == nil {
			t.Errorf("test is only working if messages are sent accoring to note top of file")
			continue
		}
		if msg.MsgType == pcan.PCAN_MESSAGE_EXTENDED {
			t.Errorf("expected only STANDARD messages, got: %v", msg.MsgType)
		}
		if msg.ID < 0x100 || msg.ID > 0x200 {
			t.Errorf("expected msg ID from 0x100 to 0x200, got: %x", msg.ID)
		}
	}

	// set filter for 11 bit messages
	state, err = pcan.SetFilter(HANDLE_FOR_TESTS, 0x200, 0x400, pcan.PCAN_MODE_EXTENDED)
	fmt.Println("Set to: HANDLE_FOR_TESTS, 0x200, 0x400, pcan.PCAN_MODE_EXTENDED ")
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	// test extended messages and msg id filter 1s
	start = time.Now()
	for time.Since(start) < 1*time.Second {
		_, msg, _, _ := auxReadBasic(5000 * time.Millisecond)
		if msg == nil {
			t.Errorf("test is only working if messages are sent accoring to note top of file")
			continue
		}
		if msg.MsgType == pcan.PCAN_MESSAGE_STANDARD {
			t.Errorf("expected only EXTENDED messages, got: %v", msg.MsgType)
		}
		if msg.ID < 0x200 || msg.ID > 0x400 {
			t.Errorf("expected msg ID from 0x200 to 0x400, got: %x", msg.ID)
		}
	}
}

func TestResetFilter(t *testing.T) {
	auxInitBasic()
	state, err := pcan.ResetFilter(HANDLE_FOR_TESTS)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestSetErrorFrames(t *testing.T) {
	auxInitBasic()
	var _trans = map[pcan.TPCANParameterValue]string{pcan.PCAN_PARAMETER_OFF: "PCAN_PARAMETER_OFF", pcan.PCAN_PARAMETER_ON: "PCAN_PARAMETER_ON"}

	// deactivate error frames
	state, err := pcan.SetParameter(HANDLE_FOR_TESTS, pcan.PCAN_ALLOW_ERROR_FRAMES, pcan.PCAN_PARAMETER_OFF)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	// check
	state, val, err := pcan.GetParameter(HANDLE_FOR_TESTS, pcan.PCAN_ALLOW_ERROR_FRAMES)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if val != pcan.PCAN_PARAMETER_OFF {
		t.Errorf("setting was not set correctly to PCAN_PARAMETER_OFF: %v", _trans[val])
	}

	// activate error frames
	state, err = pcan.SetParameter(HANDLE_FOR_TESTS, pcan.PCAN_ALLOW_ERROR_FRAMES, pcan.PCAN_PARAMETER_ON)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	// check
	state, val, err = pcan.GetParameter(HANDLE_FOR_TESTS, pcan.PCAN_ALLOW_ERROR_FRAMES)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if val != pcan.PCAN_PARAMETER_ON {
		t.Errorf("setting was not set correctly to PCAN_PARAMETER_ON: %v", _trans[val])
	}
}

func TestSetStatusFrames(t *testing.T) {
	auxInitBasic()
	var _trans = map[pcan.TPCANParameterValue]string{pcan.PCAN_PARAMETER_OFF: "PCAN_PARAMETER_OFF", pcan.PCAN_PARAMETER_ON: "PCAN_PARAMETER_ON"}

	// deactivate status frames
	state, err := pcan.SetParameter(HANDLE_FOR_TESTS, pcan.PCAN_ALLOW_STATUS_FRAMES, pcan.PCAN_PARAMETER_OFF)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	// check
	state, val, err := pcan.GetParameter(HANDLE_FOR_TESTS, pcan.PCAN_ALLOW_STATUS_FRAMES)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if val != pcan.PCAN_PARAMETER_OFF {
		t.Errorf("setting was not set correctly to PCAN_PARAMETER_OFF: %v", _trans[val])
	}

	// activate status frames
	state, err = pcan.SetParameter(HANDLE_FOR_TESTS, pcan.PCAN_ALLOW_STATUS_FRAMES, pcan.PCAN_PARAMETER_ON)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	// check
	state, val, err = pcan.GetParameter(HANDLE_FOR_TESTS, pcan.PCAN_ALLOW_STATUS_FRAMES)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if val != pcan.PCAN_PARAMETER_ON {
		t.Errorf("setting was not set correctly to PCAN_PARAMETER_ON: %v", _trans[val])
	}
}

func TestSetEchoFrames(t *testing.T) {
	auxInitBasic()
	var _trans = map[pcan.TPCANParameterValue]string{pcan.PCAN_PARAMETER_OFF: "PCAN_PARAMETER_OFF", pcan.PCAN_PARAMETER_ON: "PCAN_PARAMETER_ON"}

	// deactivate echo frames
	state, err := pcan.SetParameter(HANDLE_FOR_TESTS, pcan.PCAN_ALLOW_ECHO_FRAMES, pcan.PCAN_PARAMETER_OFF)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	// check
	state, val, err := pcan.GetParameter(HANDLE_FOR_TESTS, pcan.PCAN_ALLOW_ECHO_FRAMES)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if val != pcan.PCAN_PARAMETER_OFF {
		t.Errorf("setting was not set correctly to PCAN_PARAMETER_OFF: %v", _trans[val])
	}

	// activate echo frames
	state, err = pcan.SetParameter(HANDLE_FOR_TESTS, pcan.PCAN_ALLOW_ECHO_FRAMES, pcan.PCAN_PARAMETER_ON)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	// check
	state, val, err = pcan.GetParameter(HANDLE_FOR_TESTS, pcan.PCAN_ALLOW_ECHO_FRAMES)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if val != pcan.PCAN_PARAMETER_ON {
		t.Errorf("setting was not set correctly to PCAN_PARAMETER_ON: %v", _trans[val])
	}
}

func TestGetChannelCondition(t *testing.T) {
	auxInitBasic()
	state, val, err := pcan.GetParameter(HANDLE_FOR_TESTS, pcan.PCAN_CHANNEL_CONDITION)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if pcan.TPCANCHannelCondition(val) != pcan.PCAN_CHANNEL_OCCUPIED {
		t.Errorf("got wrong condition: %v", val)
	}
}

func TestSetIdentifying(t *testing.T) {
	auxInitBasic()

	state, err := pcan.SetParameter(HANDLE_FOR_TESTS, pcan.PCAN_CHANNEL_IDENTIFYING, pcan.PCAN_PARAMETER_ON)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestReadOnly(t *testing.T) {
	auxInitBasic()
	state, err := pcan.SetParameter(HANDLE_FOR_TESTS, pcan.PCAN_LISTEN_ONLY, pcan.PCAN_PARAMETER_ON)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	ret, msg, ts, err := auxReadBasic(5000 * time.Millisecond)
	if msg != nil {
		t.Errorf("still got a message, expected to be read only: %v", err)
		fmt.Println("state: ", ret)
		fmt.Println("msg: ", msg)
		fmt.Println("timestamp: ", ts)
		fmt.Println("error: ", err)
	}
}

func TestGetErrorText(t *testing.T) {

	pcan.LoadAPI()

	// check retval == PCAN_ERROR_OK
	state, _, err := pcan.GetErrorText(pcan.PCAN_ERROR_OK, pcan.LanguageGerman)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	// check different values
	inputs := []pcan.TPCANStatus{
		pcan.PCAN_ERROR_OK,
		pcan.PCAN_ERROR_BUSHEAVY,
		pcan.PCAN_ERROR_BUSLIGHT,
		pcan.PCAN_ERROR_ILLHANDLE,
		pcan.PCAN_ERROR_ILLOPERATION,
	}
	expOutputsGerman := []string{
		"Kein Fehler",
		"Bus-Fehler: Ein Fehlerz\xe4hler hat die \"heavy\"/\"warning\"\" Obergrenze erreicht bzw. \xfcberschritten",
		"Bus-Fehler: Ein Fehlerz\xe4hler hat die \"light\" Obergrenze erreicht bzw. \xfcberschritten",
		"Der Wert eines Handles (PCAN-Channel, PCAN-Hardware, PCAN-Net, PCAN-Client) ist ung\xfcltig",
		"Ein Vorgang ist aufgrund der aktuellen Konfiguration nicht zul\xe4ssig",
	}
	expOutputsEnglish := []string{
		"No Error",
		"Bus error: an error counter reached the 'heavy'/'warning' limit",
		"Bus error: an error counter reached the 'light' limit",
		"The value of a handle (PCAN-Channel, PCAN-Hardware, PCAN-Net, PCAN-Client) is invalid",
		"An operation is not allowed due to the current configuration",
	}

	for i, elem := range inputs {
		_, buf, _ := pcan.GetErrorText(pcan.TPCANStatus(elem), pcan.LanguageGerman)
		text := auxErrBufToText(buf)
		if text != expOutputsGerman[i] {
			t.Errorf("expected: %v, got: %v", expOutputsGerman[i], text)
		}
	}

	for i, elem := range inputs {
		_, buf, _ := pcan.GetErrorText(pcan.TPCANStatus(elem), pcan.LanguageEnglish)
		text := auxErrBufToText(buf)
		if text != expOutputsEnglish[i] {
			t.Errorf("expected: %v, got: %v", expOutputsEnglish[i], text)
		}
	}
}

// TODO This functionality is not really implemented
// func TestLookupChannel(t *testing.T) {
//
// 	var deviceType = "PCAN_USB"
// 	var deviceID = ""
// 	var controllerNumber = ""
// 	var iPAddress = ""
//
// 	pcan.LoadAPI()
//
// 	state, handle, err := pcan.LookUpChannel(deviceType, deviceID, controllerNumber, iPAddress)
// 	if state != pcan.PCAN_ERROR_OK {
// 		t.Errorf("got non okay status code: %x", state)
// 	}
// 	if err != nil {
// 		t.Errorf("got error: %v", err)
// 	}
// 	if handle != HANDLE_FOR_TESTS {
// 		t.Errorf("got invalid handle: %v", handle)
// 	}
// }

// NOTE: Connect only one channel for this test to work
func TestAttachedChannelsCount(t *testing.T) {
	pcan.LoadAPI()
	count, err := pcan.AttachedChannelsCount()
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if count != 1 {
		t.Errorf("Invalid channel count (this test needs exactly 1 connected device): %v", count)
	}
}

func TestAttachedChannels_Extended(t *testing.T) {
	pcan.LoadAPI()
	channels, err := pcan.AttachedChannels_Extended()
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if len(channels) < 1 {
		t.Errorf("got invalid channels length: %v", len(channels))
	}

	if channels[0].Channel != HANDLE_FOR_TESTS {
		t.Errorf("got invalid channel entry channel: %v", channels[0].Channel)
	}
	if channels[0].DeviceType != pcan.PCAN_USB {
		t.Errorf("got invalid channel entry device type: %v", channels[0].DeviceType)
	}

	// TODO other information missing to test: features, condition, deviceID, controllerNumber
	//
	// Example for features reading:
	// if ((features & PCANBasic.FEATURE_FD_CAPABLE) == PCANBasic.FEATURE_FD_CAPABLE)
	// {
	// 	Console.WriteLine("Channel 0x{0:X}", channelsToCheck[i]);
	// }

	// channel conditions
	unavailable := (channels[0].ChannelCondition & pcan.PCAN_CHANNEL_UNAVAILABLE) == pcan.PCAN_CHANNEL_UNAVAILABLE
	available := (channels[0].ChannelCondition & pcan.PCAN_CHANNEL_AVAILABLE) == pcan.PCAN_CHANNEL_AVAILABLE
	occupied := (channels[0].ChannelCondition & pcan.PCAN_CHANNEL_OCCUPIED) == pcan.PCAN_CHANNEL_OCCUPIED
	pcanview := (channels[0].ChannelCondition & pcan.PCAN_CHANNEL_PCANVIEW) == pcan.PCAN_CHANNEL_PCANVIEW

	fmt.Println("unavailable", unavailable)
	fmt.Println("available", available)
	fmt.Println("occupied", occupied)
	fmt.Println("pcanview", pcanview)
}

func TestAttachedChannelsName(t *testing.T) {
	pcan.LoadAPI()
	channels, err := pcan.AttachedChannelsNames()
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if len(channels) < 1 {
		t.Errorf("got invalid channels length: %v", len(channels))
	}
	if channels[0] != "PCAN_USBBUS1" {
		t.Errorf("got invalid channels name: %v", channels[0])
	}
}
