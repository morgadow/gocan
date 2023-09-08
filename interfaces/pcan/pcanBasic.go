package pcan

import (
	"errors"
	"runtime"
	"syscall"
	"unsafe"
)

const driverWin = "PCANBasic.dll"     // PCAN windows driver dll name, which can be imported once the PCAN driver is installed
const driverMac = "libPCBUSB.dylib"   // PCAN max driver file name, which can be imported once the PCAN driver is installed
const driverLinux = "libpcanbasic.so" // PCAN linux driver file name, which can be imported once the PCAN driver is installed

var ErrAPINotLoadedOrFound = errors.New("pcan api not loaded or installed, please load api over pcan.LoadAPI")

var pcanAPIHandle *syscall.DLL = nil // procedure handle for PCAN driver
var pHandleInitialize *syscall.Proc = nil
var pHandleInitializeFD *syscall.Proc = nil
var pHandleUninitialize *syscall.Proc = nil
var pHandleReset *syscall.Proc = nil
var pHandleGetStatus *syscall.Proc = nil
var pHandleRead *syscall.Proc = nil
var pHandleReadFD *syscall.Proc = nil
var pHandleWrite *syscall.Proc = nil
var pHandleWriteFD *syscall.Proc = nil
var pHandleFilterMessages *syscall.Proc = nil
var pHandleGetValue *syscall.Proc = nil
var pHandleSetValue *syscall.Proc = nil
var pHandleGetErrorText *syscall.Proc = nil
var pHandleLookUpChannel *syscall.Proc = nil
var apiLoaded bool

// Loads PCAN API (.ddl) file
func LoadAPI() error {
	var err error = nil
	var driver = ""

	// evaluate operating system and architecture and select driver file
	switch runtime.GOOS {
	case "windows":
		driver = driverWin
	case "darwin":
		driver = driverMac
	default:
		driver = driverLinux
	}

	pcanAPIHandle, err = syscall.LoadDLL(driver)
	if err != nil || pcanAPIHandle == nil {
		return err
	}

	pHandleInitialize, _ = pcanAPIHandle.FindProc("CAN_Initialize")
	pHandleInitializeFD, _ = pcanAPIHandle.FindProc("CAN_InitializeFD")
	pHandleUninitialize, _ = pcanAPIHandle.FindProc("CAN_Uninitialize")
	pHandleReset, _ = pcanAPIHandle.FindProc("CAN_Reset")
	pHandleGetStatus, _ = pcanAPIHandle.FindProc("CAN_GetStatus")
	pHandleRead, _ = pcanAPIHandle.FindProc("CAN_Read")
	pHandleReadFD, _ = pcanAPIHandle.FindProc("CAN_ReadFD")
	pHandleWrite, _ = pcanAPIHandle.FindProc("CAN_Write")
	pHandleWriteFD, _ = pcanAPIHandle.FindProc("CAN_WriteFD")
	pHandleFilterMessages, _ = pcanAPIHandle.FindProc("CAN_FilterMessages")
	pHandleGetValue, _ = pcanAPIHandle.FindProc("CAN_GetValue")
	pHandleSetValue, _ = pcanAPIHandle.FindProc("CAN_SetValue")
	pHandleGetErrorText, _ = pcanAPIHandle.FindProc("CAN_GetErrorText")
	pHandleLookUpChannel, _ = pcanAPIHandle.FindProc("CAN_LookUpChannel")

	apiLoaded = pHandleInitialize != nil && pHandleInitializeFD != nil && pHandleReset != nil && pHandleGetStatus != nil &&
		pHandleRead != nil && pHandleReadFD != nil && pHandleWrite != nil && pHandleWriteFD != nil && pHandleFilterMessages != nil && pHandleGetValue != nil &&
		pHandleSetValue != nil && pHandleGetErrorText != nil && pHandleLookUpChannel != nil && pHandleUninitialize != nil

	if !apiLoaded {
		return errors.New("could not load pointers to pcan functions")
	}
	return nil
}

// Unloads PCAN API (.ddl) file
func UnloadAPI() error {

	// reset pointers
	pHandleInitialize = nil
	pHandleInitializeFD = nil
	pHandleUninitialize = nil
	pHandleReset = nil
	pHandleGetStatus = nil
	pHandleRead = nil
	pHandleReadFD = nil
	pHandleWrite = nil
	pHandleWriteFD = nil
	pHandleFilterMessages = nil
	pHandleGetValue = nil
	pHandleSetValue = nil
	pHandleGetErrorText = nil
	pHandleLookUpChannel = nil
	pHandleUninitialize = nil
	apiLoaded = false

	err := pcanAPIHandle.Release()
	return err
}

// Initializes a PCAN Channel
// Channel: The handle of a PCAN Channel
// baudRate: The speed for the communication (BTR0BTR1 code)
// hwType: Non-PnP: The type of hardware and operation mode
// ioPort: Non-PnP: The I/O address for the parallel port
// interrupt: Non-PnP: Interrupt number of the parallel por
func Initialize(channel TPCANHandle, baudRate TPCANBaudrate, hwType TPCANType, ioPort uint32, interrupt uint16) (TPCANStatus, error) {
	r1, _, errno := pHandleInitialize.Call(uintptr(channel), uintptr(baudRate), uintptr(hwType), uintptr(ioPort), uintptr(interrupt))
	return TPCANStatus(r1), sysCallErr(errno)
}

// Initializes a FD capable PCAN Channel
// Channel: The handle of a PCAN Channel
// bitRateFD: The speed for the communication (FD bit rate string)
// Note:
// Bit rate string must follow the following construction rules:
//   - parameter and values must be separated by '='
//   - Couples of Parameter/value must be separated by ','
//   - Following Parameter must be filled out: f_clock, data_brp, data_sjw, data_tseg1, data_tseg2,
//     nom_brp, nom_sjw, nom_tseg1, nom_tseg2.
//   - Following Parameters are optional (not used yet): data_ssp_offset, nom_sam
//   - Example: f_clock=80000000,nom_brp=10,nom_tseg1=5,nom_tseg2=2,nom_sjw=1,data_brp=4,data_tseg1=7,data_tseg2=2,data_sjw=1
func InitializeFD(channel TPCANHandle, bitRateFD TPCANBitrateFD) (TPCANStatus, error) {
	ret, _, errno := pHandleInitializeFD.Call(uintptr(channel), uintptr(unsafe.Pointer(&bitRateFD)))
	return TPCANStatus(ret), sysCallErr(errno)
}

// Uninitializes PCAN Channels initialized by CAN_Initialize
// Channel: The handle of a PCAN Channel
func Uninitialize(channel TPCANHandle) (TPCANStatus, error) {
	ret, _, errno := pHandleUninitialize.Call(uintptr(channel))
	return TPCANStatus(ret), sysCallErr(errno)
}

// Resets the receive and transmit queues of the PCAN Channel
// Channel: The handle of a PCAN Channel
func Reset(channel TPCANHandle) (TPCANStatus, error) {
	ret, _, errno := pHandleReset.Call(uintptr(channel))
	return TPCANStatus(ret), sysCallErr(errno)
}

// Gets the current status of a PCAN Channel
// Channel: The handle of a PCAN Channel
func GetStatus(channel TPCANHandle) (TPCANStatus, error) {
	ret, _, errno := pHandleGetStatus.Call(uintptr(channel))
	return TPCANStatus(ret), sysCallErr(errno)
}

// Reads a CAN message from the receive queue of a PCAN Channel
// Channel: The handle of a PCAN Channel
func Read(channel TPCANHandle) (TPCANStatus, TPCANMsg, TPCANTimestamp, error) {
	var msg TPCANMsg
	var timeStamp = TPCANTimestamp{}

	ret, _, errno := pHandleRead.Call(uintptr(channel), uintptr(unsafe.Pointer(&msg)), uintptr(unsafe.Pointer(&timeStamp)))
	return TPCANStatus(ret), msg, timeStamp, sysCallErr(errno)
}

// Reads a CAN message from the receive queue of a FD capable PCAN Channel
// Channel: The handle of a PCAN Channel
func ReadFD(channel TPCANHandle) (TPCANStatus, TPCANMsgFD, TPCANTimestampFD, error) {
	var msgFD TPCANMsgFD
	var timeStampFD TPCANTimestampFD

	ret, _, errno := pHandleReadFD.Call(uintptr(channel), uintptr(unsafe.Pointer(&msgFD)), uintptr(unsafe.Pointer(&timeStampFD)))
	return TPCANStatus(ret), msgFD, timeStampFD, sysCallErr(errno)
}

// Transmits a CAN message
// Channel: The handle of a PCAN Channel
// msg: A Message struct with the message to be sent
func Write(channel TPCANHandle, msg TPCANMsg) (TPCANStatus, error) {
	ret, _, errno := pHandleWrite.Call(uintptr(channel), uintptr(unsafe.Pointer(&msg)))
	return TPCANStatus(ret), sysCallErr(errno)
}

// Transmits a CAN message over a FD capable PCAN Channel
// Channel: The handle of a PCAN Channel
// msgFD A MessageFD struct with the message to be sent
func WriteFD(channel TPCANHandle, msgFD TPCANMsgFD) (TPCANStatus, error) {
	ret, _, errno := pHandleWriteFD.Call(uintptr(channel), uintptr(unsafe.Pointer(&msgFD)))
	return TPCANStatus(ret), sysCallErr(errno)
}

// Configures the reception filter
// Channel: The handle of a PCAN Channel
// fromID: The lowest CAN ID to be received
// toID: The highest CAN ID to be received
// mode: Message type, Standard (11-bit identifier) or Extended (29-bit identifier)
func FilterMessages(channel TPCANHandle, fromID TPCANMsgID, toID TPCANMsgID, mode TPCANMode) (TPCANStatus, error) {
	ret, _, errno := pHandleFilterMessages.Call(uintptr(channel), uintptr(fromID), uintptr(toID), uintptr(mode))
	return TPCANStatus(ret), sysCallErr(errno)
}

// the filter applied to handle
func ResetFilter(channel TPCANHandle) (TPCANStatus, error) {
	return PCAN_ERROR_OK, nil // TODO this function is empty!?
}

// Retrieves a PCAN Channel value using a defined parameter value type
// Channel: The handle of a PCAN Channel
// param: The TPCANParameter parameter to get
// Note: Parameters can be present or not according with the kind of Hardware (PCAN Channel) being used.
// If a parameter is not available, a PCAN_ERROR_ILLPARAMTYPE error will be returned
func GetParameter(channel TPCANHandle, param TPCANParameter) (TPCANStatus, TPCANParameterValue, error) {
	var val TPCANParameterValue
	ret, err := GetValue(channel, param, unsafe.Pointer(&val), uint32(unsafe.Sizeof(val)))
	return TPCANStatus(ret), val, err
}

// Configures a PCAN Channel value using a defined parameter value type
// Channel: The handle of a PCAN Channel
// param: The TPCANParameter parameter to set
// value: Value of parameter
// Note: Parameters can be present or not according with the kind of Hardware (PCAN Channel) being used.
// If a parameter is not available, a PCAN_ERROR_ILLPARAMTYPE error will be returned
func SetParameter(channel TPCANHandle, param TPCANParameter, val TPCANParameterValue) (TPCANStatus, error) {
	ret, errno := SetValue(channel, param, unsafe.Pointer(&val), unsafe.Sizeof(val))
	return TPCANStatus(ret), sysCallErr(errno)
}

// Retrieves a PCAN Channel value
// Channel: The handle of a PCAN Channel
// param: The TPCANParameter parameter to get
// Note: Parameters can be present or not according with the kind
// Note: Parameters can be present or not according with the kind of Hardware (PCAN Channel) being used.
// If a parameter is not available, a PCAN_ERROR_ILLPARAMTYPE error will be returned
func GetValue(channel TPCANHandle, param TPCANParameter, buffer unsafe.Pointer, bufferSize uint32) (TPCANStatus, error) { // TODO change buffersize to uintptr
	ret, _, errno := pHandleGetValue.Call(uintptr(channel), uintptr(param), uintptr(buffer), uintptr(bufferSize))
	return TPCANStatus(ret), sysCallErr(errno)
}

// Configures a PCAN Channel value.
// Channel: The handle of a PCAN Channel
// param: The TPCANParameter parameter to set
// value: Value of parameter
// Note: Parameters can be present or not according with the kind of Hardware (PCAN Channel) being used.
// If a parameter is not available, a PCAN_ERROR_ILLPARAMTYPE error will be returned
func SetValue(channel TPCANHandle, param TPCANParameter, buffer unsafe.Pointer, bufferSize uintptr) (TPCANStatus, error) {
	ret, _, errno := pHandleSetValue.Call(uintptr(channel), uintptr(param), uintptr(buffer), bufferSize)
	return TPCANStatus(ret), sysCallErr(errno)
}

// Returns a descriptive text of a given TPCANStatus error code, in any desired language
// err: A TPCANStatus error code
// language: Indicates a 'Primary language ID'
func GetErrorText(status TPCANStatus, language TPCANLanguage) (TPCANStatus, [MAX_LENGHT_STRING_BUFFER]byte, error) {
	var buffer [MAX_LENGHT_STRING_BUFFER]byte

	ret, _, errno := pHandleGetErrorText.Call(uintptr(status), uintptr(language), uintptr(unsafe.Pointer(&buffer)))
	return TPCANStatus(ret), buffer, sysCallErr(errno)
}

// Finds a PCAN-Basic Channel that matches with the given parameters
// parameters: A comma separated string contained pairs of parameter-name/value to be matched within a PCAN-Basic Channel
// foundChannels: Buffer for returning the PCAN-Basic Channel when found
func LookUpChannel(deviceType string, deviceID string, controllerNumber string, ipAdress string) (TPCANStatus, TPCANHandle, error) {

	var sParameters string = ""
	var foundChannel TPCANHandle

	// merge search parameter
	if deviceType != "" {
		sParameters += string(LOOKUP_DEVICE_TYPE) + "=" + deviceType
	}

	if deviceID != "" {
		if sParameters != "" {
			sParameters += ", "
		}
		sParameters += string(LOOKUP_DEVICE_ID) + "=" + deviceID
	}
	if controllerNumber != "" {
		if sParameters != "" {
			sParameters += ", "
		}
		sParameters += string(LOOKUP_CONTROLLER_NUMBER) + "=" + controllerNumber
	}
	if ipAdress != "" {
		if sParameters != "" {
			sParameters += ", "
		}
		sParameters += string(LOOKUP_IP_ADDRESS) + "=" + ipAdress
	}

	ret, _, errno := pHandleLookUpChannel.Call(uintptr(unsafe.Pointer(&sParameters)), uintptr(unsafe.Pointer(&foundChannel)))
	return TPCANStatus(ret), foundChannel, sysCallErr(errno)
}

// helper function to handle syscall return value
func sysCallErr(err error) error {
	errno := err.(syscall.Errno)
	if errno != 0 {
		return errors.New(errno.Error())
	}
	return nil
}
