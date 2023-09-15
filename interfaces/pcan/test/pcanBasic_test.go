package test

import (
	"fmt"
	"testing"
	"time"
	"unsafe"

	"github.com/morgadow/gocan/interfaces/pcan"
)

// For this tests, connect a PCAN USB Bus and periodically send four messages:
// - 11bit 500 ms: 0x123 - {1,2,3,4,5,6,7,8}
// - 11bit 500 ms: 0x321 - {8,7,6,5,4,3,2,1}
// - 29bit 500 ms: 0x0000123 - {1,2,3,4,5,6,7,8}
// - 29bit 500 ms: 0x0000321 - {8,7,6,5,4,3,2,1}

func init() {
	err := pcan.LoadAPI()
	if err != nil {
		panic(err)
	}
}

func auxInitBasic(channel pcan.TPCANHandle) {
	state, err := pcan.InitializeBasic(channel, pcan.PCAN_BAUD_500K)
	if err != nil || state != pcan.PCAN_ERROR_OK {
		fmt.Println(fmt.Errorf("Helper function got error or invalid stats: %v - %w", state, err))
		panic(err)
	}
	state, err = pcan.Uninitialize(pcan.PCAN_USBBUS1)
	fmt.Println(state)
	fmt.Println(err)
}

func auxInitFD(channel pcan.TPCANHandle) {
	state, err := pcan.InitializeFD(channel, "")
	if err != nil || state != pcan.PCAN_ERROR_OK {
		fmt.Println(fmt.Errorf("Helper function got error or invalid stats: %v - %w", state, err))
	}
}

func auxReadBasic(channel pcan.TPCANHandle, timeout uint32) (pcan.TPCANStatus, pcan.TPCANMsg, pcan.TPCANTimestamp, error) {
	var state pcan.TPCANStatus
	var msg pcan.TPCANMsg
	var timestamp pcan.TPCANTimestamp
	var err error

	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	endTime := startTime + int64(timeout) // timeout is [ms]

	for {
		state, msg, timestamp, err = pcan.Read(pcan.PCAN_USBBUS1)
		if state != pcan.PCAN_ERROR_QRCVEMPTY {
			if msg.MsgType == pcan.PCAN_MESSAGE_STANDARD || msg.MsgType == pcan.PCAN_MESSAGE_EXTENDED {
				fmt.Printf("Got Message: %v - ID: %x, Data: %v, DLC: %v, Type: %v\n", timestamp, msg.ID, msg.DLC, msg.Data, msg.MsgType)
				break
			}
		}
		if time.Now().UnixNano()/int64(time.Millisecond) > endTime {
			return pcan.PCAN_ERROR_QRCVEMPTY, msg, timestamp, nil
		}
		time.Sleep(1 * time.Millisecond)
	}
	return state, msg, timestamp, err
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

	state, err := pcan.InitializeBasic(pcan.PCAN_USBBUS1, pcan.PCAN_BAUD_500K)

	if state != pcan.PCAN_ERROR_OK && state != pcan.PCAN_ERROR_BUSLIGHT && state != pcan.PCAN_ERROR_BUSHEAVY {

		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestInitialize(t *testing.T) {

	state, err := pcan.Initialize(pcan.PCAN_USBBUS1, pcan.PCAN_BAUD_500K, pcan.PCAN_TYPE_ISA, 0x02A0, 11)

	if state != pcan.PCAN_ERROR_OK && state != pcan.PCAN_ERROR_BUSLIGHT && state != pcan.PCAN_ERROR_BUSHEAVY {

		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestUnitialize(t *testing.T) {

	auxInitBasic(pcan.PCAN_USBBUS1)
	state, err := pcan.Uninitialize(pcan.PCAN_USBBUS1)

	if state != pcan.PCAN_ERROR_OK && state != pcan.PCAN_ERROR_BUSLIGHT && state != pcan.PCAN_ERROR_BUSHEAVY {

		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestReset(t *testing.T) {

	auxInitBasic(pcan.PCAN_USBBUS1)
	state, err := pcan.Reset(pcan.PCAN_USBBUS1)

	if state != pcan.PCAN_ERROR_OK && state != pcan.PCAN_ERROR_BUSLIGHT && state != pcan.PCAN_ERROR_BUSHEAVY {

		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}
func TestGetStatus(t *testing.T) {

	auxInitBasic(pcan.PCAN_USBBUS1)
	state, err := pcan.GetStatus(pcan.PCAN_USBBUS1)

	if state != pcan.PCAN_ERROR_OK && state != pcan.PCAN_ERROR_BUSLIGHT && state != pcan.PCAN_ERROR_BUSHEAVY {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestRead(t *testing.T) {

	auxInitBasic(pcan.PCAN_USBBUS1)

	state, msg, timestamp, err := auxReadBasic(pcan.PCAN_USBBUS1, 5000)

	if msg.ID != 0x123 && msg.ID != 0x321 && msg.ID != 0x00000123 && msg.ID != 0x00000321 {
		t.Errorf("wrong message id: %v", msg.ID)
	}
	if msg.Data != [8]byte{8, 7, 6, 5, 4, 3, 2, 1} && msg.Data != [8]byte{1, 2, 3, 4, 5, 6, 7, 8} {
		t.Errorf("wrong message data: %v", msg.Data)
	}
	if (timestamp.Micros == 0 || timestamp.Micros == 0xFFFF) && (timestamp.Millis == 0 || timestamp.Millis == 0xFFFFFFFF) && (uint32(timestamp.MillisOverflow) == 0 || timestamp.MillisOverflow == 0xFFFF) {
		t.Errorf("wrong message timestamp: %v", timestamp)

	}

	if state != pcan.PCAN_ERROR_OK && state != pcan.PCAN_ERROR_BUSLIGHT && state != pcan.PCAN_ERROR_BUSHEAVY {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestReadFD(t *testing.T) {

	auxInitFD(pcan.PCAN_USBBUS1)
	state, msg, timestamp, err := pcan.ReadFD(pcan.PCAN_USBBUS1)

	fmt.Printf("Got Message: %v - ID: %v, Data: %v, DLC: %v, Type: %v", timestamp, msg.ID, msg.DLC, msg.Data, msg.MsgType)

	if state != pcan.PCAN_ERROR_OK && state != pcan.PCAN_ERROR_BUSLIGHT && state != pcan.PCAN_ERROR_BUSHEAVY {

		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestWrite(t *testing.T) {

	auxInitBasic(pcan.PCAN_USBBUS1)
	msg := pcan.TPCANMsg{ID: 0x123, Data: [8]byte{0, 1, 2, 3, 4, 5, 6, 7}}
	state, err := pcan.Write(pcan.PCAN_USBBUS1, msg)

	if state != pcan.PCAN_ERROR_OK && state != pcan.PCAN_ERROR_BUSLIGHT && state != pcan.PCAN_ERROR_BUSHEAVY {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestWriteFD(t *testing.T) {

	auxInitFD(pcan.PCAN_USBBUS1)
	msg := pcan.TPCANMsgFD{ID: 0x123, Data: [64]byte{}}
	state, err := pcan.WriteFD(pcan.PCAN_USBBUS1, msg)

	if state != pcan.PCAN_ERROR_OK && state != pcan.PCAN_ERROR_BUSLIGHT && state != pcan.PCAN_ERROR_BUSHEAVY {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestFilterMessages(t *testing.T) {

	auxInitBasic(pcan.PCAN_USBBUS1)

	// check for 29 bit messages
	state, err := pcan.FilterMessages(pcan.PCAN_USBBUS1, 0x100, 0x200, pcan.PCAN_MODE_STANDARD)
	fmt.Println("Set to: pcan.PCAN_USBBUS1, 0x100, 0x200, pcan.PCAN_MODE_STANDARD")
	if state != pcan.PCAN_ERROR_OK && state != pcan.PCAN_ERROR_BUSLIGHT && state != pcan.PCAN_ERROR_BUSHEAVY {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	// test standard messages and msg id filter
	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	endTime := startTime + int64(2000) // 2 second timeout

	for {
		_, msg, _, _ := auxReadBasic(pcan.PCAN_USBBUS1, 5000)

		if msg.MsgType == pcan.PCAN_MESSAGE_EXTENDED {
			t.Errorf("expected only STANDARD messages, got: %v", msg.MsgType)
		}
		if msg.ID < 0x100 || msg.ID > 0x200 {
			t.Errorf("expected msg ID from 0x100 to 0x200, got: %x", msg.ID)
		}

		if time.Now().UnixNano()/int64(time.Millisecond) > endTime {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}

	state, err = pcan.FilterMessages(pcan.PCAN_USBBUS1, 0x200, 0x400, pcan.PCAN_MODE_EXTENDED)
	fmt.Println("Set to: pcan.PCAN_USBBUS1, 0x200, 0x400, pcan.PCAN_MODE_EXTENDED ")
	if state != pcan.PCAN_ERROR_OK && state != pcan.PCAN_ERROR_BUSLIGHT && state != pcan.PCAN_ERROR_BUSHEAVY {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	// test extended messages and msg id filter
	startTime = time.Now().UnixNano() / int64(time.Millisecond)
	endTime = startTime + int64(2000) // 2 second timeout

	for {
		_, msg, _, _ := auxReadBasic(pcan.PCAN_USBBUS1, 5000)

		if msg.MsgType == pcan.PCAN_MESSAGE_STANDARD {
			t.Errorf("expected only EXTENDED messages, got: %v", msg.MsgType)
		}
		if msg.ID < 0x200 || msg.ID > 0x400 {
			t.Errorf("expected msg ID from 0x200 to 0x400, got: %x", msg.ID)
		}

		if time.Now().UnixNano()/int64(time.Millisecond) > endTime {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}

	// TODO maybe integrate this in PCANBus level instead: state, err = pcan.SetValue(pcan.PCAN_USBBUS1, pcan.PCAN_MESSAGE_FILTER, uint32(pcan.PCAN_FILTER_CUSTOM))
}

func TestResetFilter(t *testing.T) {
	t.Errorf("test not implemented")
}

func TestGetParameter(t *testing.T) {
	auxInitBasic(pcan.PCAN_USBBUS1)
	state, val, err := pcan.GetParameter(pcan.PCAN_USBBUS1, pcan.PCAN_CHANNEL_CONDITION)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if val != pcan.PCAN_CHANNEL_OCCUPIED {
		t.Errorf("got wrong condition: %v", val)
	}
}

func TestSetParameter(t *testing.T) {
	auxInitBasic(pcan.PCAN_USBBUS1)

	state, err := pcan.SetParameter(pcan.PCAN_USBBUS1, pcan.PCAN_CHANNEL_IDENTIFYING, pcan.PCAN_PARAMETER_ON)
	if state != pcan.PCAN_ERROR_OK {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestGetValue(t *testing.T) {
	var channelCount uint32
	state, err := pcan.GetValue(pcan.PCAN_NONEBUS, pcan.PCAN_ATTACHED_CHANNELS_COUNT, unsafe.Pointer(&channelCount), uint32(unsafe.Sizeof(channelCount)))
	if state != pcan.PCAN_ERROR_OK && state != pcan.PCAN_ERROR_BUSLIGHT && state != pcan.PCAN_ERROR_BUSHEAVY {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if channelCount != 1 {
		t.Errorf("got invalid parameter value back: %v", channelCount)
	}
}

func TestSetValue(t *testing.T) {
	auxInitBasic(pcan.PCAN_USBBUS1)
	state, err := pcan.SetParameter(pcan.PCAN_USBBUS1, pcan.PCAN_LISTEN_ONLY, pcan.PCAN_PARAMETER_ON)
	if state != pcan.PCAN_ERROR_OK && state != pcan.PCAN_ERROR_BUSLIGHT && state != pcan.PCAN_ERROR_BUSHEAVY {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	ret, msg, ts, err := auxReadBasic(pcan.PCAN_USBBUS1, 5000)
	if msg.ID != 0x0 || msg.Data != [pcan.LENGTH_DATA_CAN_MESSAGE]byte{0, 0, 0, 0, 0, 0, 0, 0} {
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

func TestLookupChannel(t *testing.T) {

	var deviceType = "PCAN_USB"
	var deviceID = ""
	var controllerNumber = ""
	var iPAddress = ""

	pcan.LoadAPI()

	state, handle, err := pcan.LookUpChannel(deviceType, deviceID, controllerNumber, iPAddress)
	if state != pcan.PCAN_ERROR_OK && state != pcan.PCAN_ERROR_BUSLIGHT && state != pcan.PCAN_ERROR_BUSHEAVY {
		t.Errorf("got non okay status code: %x", state)
	}
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if handle != pcan.PCAN_USBBUS1 {
		t.Errorf("got invalid handle: %v", handle)
	}
}

func TestAttachedChannels(t *testing.T) {
	pcan.LoadAPI()
	channels, err := pcan.AttachedChannels()
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if len(channels) < 1 {
		t.Errorf("got invalid channels length: %v", len(channels))
	}

	if channels[0].Channel != pcan.PCAN_USBBUS1 {
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
	unavailable := (channels[0].ChannelCondition & uint32(pcan.PCAN_CHANNEL_UNAVAILABLE)) == uint32(pcan.PCAN_CHANNEL_UNAVAILABLE)
	available := (channels[0].ChannelCondition & uint32(pcan.PCAN_CHANNEL_AVAILABLE)) == uint32(pcan.PCAN_CHANNEL_AVAILABLE)
	occupied := (channels[0].ChannelCondition & uint32(pcan.PCAN_CHANNEL_OCCUPIED)) == uint32(pcan.PCAN_CHANNEL_OCCUPIED)
	pcanview := (channels[0].ChannelCondition & uint32(pcan.PCAN_CHANNEL_PCANVIEW)) == uint32(pcan.PCAN_CHANNEL_PCANVIEW)

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

func TestAttachedChannelsCount(t *testing.T) {
	pcan.LoadAPI()
	count, err := pcan.AttachedChannelsCount()
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if count != 1 {
		t.Errorf("got wrong channel count: %v", count)
	}
}
