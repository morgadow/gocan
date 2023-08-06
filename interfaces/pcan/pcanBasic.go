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

var pcanAPIHandle syscall.Handle = 0 // process handle for PCAN driver
var pHandleInitialize uintptr = 0
var pHandleInitializeFD uintptr = 0
var pHandleUninitialize uintptr = 0
var pHandleReset uintptr = 0
var pHandleGetStatus uintptr = 0
var pHandleRead uintptr = 0
var pHandleReadFD uintptr = 0
var pHandleWrite uintptr = 0
var pHandleWriteFD uintptr = 0
var pHandleFilterMessages uintptr = 0
var pHandleGetValue uintptr = 0
var pHandleSetValue uintptr = 0
var pHandleGetErrorText uintptr = 0
var pHandleLookUpChannel uintptr = 0
var apiLoaded bool

// Loads PCAN API (.ddl) file
func LoadAPI() error {
	var err error = nil
	var dll = ""

	// evaluate operating system and architecture and select driver file
	switch runtime.GOOS {
	case "windows":
		dll = driverWin
	case "darwin":
		dll = driverMac
	default:
		dll = driverLinux
	}

	pcanAPIHandle, err = syscall.LoadLibrary(dll)
	if err != nil {
		return err
	}

	pHandleInitialize, _ = syscall.GetProcAddress(pcanAPIHandle, "CAN_Initialize")
	pHandleInitializeFD, _ = syscall.GetProcAddress(pcanAPIHandle, "CAN_InitializeFD")
	pHandleUninitialize, _ = syscall.GetProcAddress(pcanAPIHandle, "CAN_Uninitialize")
	pHandleReset, _ = syscall.GetProcAddress(pcanAPIHandle, "CAN_Reset")
	pHandleGetStatus, _ = syscall.GetProcAddress(pcanAPIHandle, "CAN_GetStatus")
	pHandleRead, _ = syscall.GetProcAddress(pcanAPIHandle, "CAN_Read")
	pHandleReadFD, _ = syscall.GetProcAddress(pcanAPIHandle, "CAN_ReadFD")
	pHandleWrite, _ = syscall.GetProcAddress(pcanAPIHandle, "CAN_Write")
	pHandleWriteFD, _ = syscall.GetProcAddress(pcanAPIHandle, "CAN_WriteFD")
	pHandleFilterMessages, _ = syscall.GetProcAddress(pcanAPIHandle, "CAN_FilterMessages")
	pHandleGetValue, _ = syscall.GetProcAddress(pcanAPIHandle, "CAN_GetValue")
	pHandleSetValue, _ = syscall.GetProcAddress(pcanAPIHandle, "CAN_SetValue")
	pHandleGetErrorText, _ = syscall.GetProcAddress(pcanAPIHandle, "CAN_GetErrorText")
	pHandleLookUpChannel, _ = syscall.GetProcAddress(pcanAPIHandle, "CAN_LookUpChannel")

	apiLoaded = pHandleInitialize > 0 && pHandleInitializeFD > 0 && pHandleReset > 0 && pHandleGetStatus > 0 &&
		pHandleRead > 0 && pHandleReadFD > 0 && pHandleWrite > 0 && pHandleWriteFD > 0 && pHandleFilterMessages > 0 && pHandleGetValue > 0 &&
		pHandleSetValue > 0 && pHandleGetErrorText > 0 && pHandleLookUpChannel > 0 && pHandleUninitialize > 0

	if !apiLoaded {
		return errors.New("could not load pointers to pcan functions")
	}
	return nil
}

// Unloads PCAN API (.ddl) file
func UnloadAPI() error {

	// reset pointers
	pHandleInitialize = 0
	pHandleInitializeFD = 0
	pHandleUninitialize = 0
	pHandleReset = 0
	pHandleGetStatus = 0
	pHandleRead = 0
	pHandleReadFD = 0
	pHandleWrite = 0
	pHandleWriteFD = 0
	pHandleFilterMessages = 0
	pHandleGetValue = 0
	pHandleSetValue = 0
	pHandleGetErrorText = 0
	pHandleLookUpChannel = 0
	pHandleUninitialize = 0
	apiLoaded = false

	err := syscall.FreeLibrary(pcanAPIHandle)
	return err
}

// Initializes a PCAN Channel
// Channel: The handle of a PCAN Channel
// baudRate: The speed for the communication (BTR0BTR1 code)
// hwType: Non-PnP: The type of hardware and operation mode
// ioPort: Non-PnP: The I/O address for the parallel port
// interrupt: Non-PnP: Interrupt number of the parallel por
func Initialize(channel TPCANHandle, baudRate TPCANBaudrate, hwType TPCANType, ioPort uint32, interrupt uint16) (TPCANStatus, error) {
	ret, _, errCall := syscall.SyscallN(pHandleInitialize, uintptr(channel), uintptr(baudRate), uintptr(hwType), uintptr(ioPort), uintptr(interrupt))
	return TPCANStatus(ret), sysCallErr(errCall)
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
	ret, _, errCall := syscall.SyscallN(pHandleInitializeFD, uintptr(channel), uintptr(unsafe.Pointer(&bitRateFD)))
	return TPCANStatus(ret), sysCallErr(errCall)
}

// Uninitializes PCAN Channels initialized by CAN_Initialize
// Channel: The handle of a PCAN Channel
func Uninitialize(channel TPCANHandle) (TPCANStatus, error) {
	ret, _, errCall := syscall.SyscallN(pHandleUninitialize, uintptr(channel))
	return TPCANStatus(ret), sysCallErr(errCall)
}

// Resets the receive and transmit queues of the PCAN Channel
// Channel: The handle of a PCAN Channel
func Reset(channel TPCANHandle) (TPCANStatus, error) {
	ret, _, errCall := syscall.SyscallN(pHandleReset, uintptr(channel))
	return TPCANStatus(ret), sysCallErr(errCall)
}

// Gets the current status of a PCAN Channel
// Channel: The handle of a PCAN Channel
func GetStatus(channel TPCANHandle) (TPCANStatus, error) {
	ret, _, errCall := syscall.SyscallN(pHandleGetStatus, uintptr(channel))
	return TPCANStatus(ret), sysCallErr(errCall)
}

// Reads a CAN message from the receive queue of a PCAN Channel
// Channel: The handle of a PCAN Channel
func Read(channel TPCANHandle) (TPCANStatus, TPCANMsg, TPCANTimestamp, error) {
	var msg TPCANMsg
	var timeStamp = TPCANTimestamp{}

	ret, _, errCall := syscall.SyscallN(pHandleRead, uintptr(channel), uintptr(unsafe.Pointer(&msg)), uintptr(unsafe.Pointer(&timeStamp)))
	return TPCANStatus(ret), msg, timeStamp, sysCallErr(errCall)
}

// Reads a CAN message from the receive queue of a FD capable PCAN Channel
// Channel: The handle of a PCAN Channel
func ReadFD(channel TPCANHandle) (TPCANStatus, TPCANMsgFD, TPCANTimestampFD, error) {
	var msgFD TPCANMsgFD
	var timeStampFD TPCANTimestampFD

	ret, _, errCall := syscall.SyscallN(pHandleReadFD, uintptr(channel), uintptr(unsafe.Pointer(&msgFD)), uintptr(unsafe.Pointer(&timeStampFD)))
	return TPCANStatus(ret), msgFD, timeStampFD, sysCallErr(errCall)
}

// Transmits a CAN message
// Channel: The handle of a PCAN Channel
// msg: A Message struct with the message to be sent
func Write(channel TPCANHandle, msg TPCANMsg) (TPCANStatus, error) {
	ret, _, errCall := syscall.SyscallN(pHandleWrite, uintptr(channel), uintptr(unsafe.Pointer(&msg)))
	return TPCANStatus(ret), sysCallErr(errCall)
}

// Transmits a CAN message over a FD capable PCAN Channel
// Channel: The handle of a PCAN Channel
// msgFD A MessageFD struct with the message to be sent
func WriteFD(channel TPCANHandle, msgFD TPCANMsgFD) (TPCANStatus, error) {
	ret, _, errCall := syscall.SyscallN(pHandleWriteFD, uintptr(channel), uintptr(unsafe.Pointer(&msgFD)))
	return TPCANStatus(ret), sysCallErr(errCall)
}

// Configures the reception filter
// Channel: The handle of a PCAN Channel
// fromID: The lowest CAN ID to be received
// toID: The highest CAN ID to be received
// mode: Message type, Standard (11-bit identifier) or Extended (29-bit identifier)
func FilterMessages(channel TPCANHandle, fromID TPCANMsgID, toID TPCANMsgID, mode TPCANMode) (TPCANStatus, error) {
	ret, _, errCall := syscall.SyscallN(pHandleFilterMessages, uintptr(channel), uintptr(fromID), uintptr(toID), uintptr(mode))
	return TPCANStatus(ret), sysCallErr(errCall)
}

// the filter applied to handle
func ResetFilter(channel TPCANHandle) (TPCANStatus, error) {
	return PCAN_ERROR_OK, nil
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
	ret, err := SetValue(channel, param, unsafe.Pointer(&val), unsafe.Sizeof(val))
	return TPCANStatus(ret), err
}

// Retrieves a PCAN Channel value
// Channel: The handle of a PCAN Channel
// param: The TPCANParameter parameter to get
// Note: Parameters can be present or not according with the kind
// Note: Parameters can be present or not according with the kind of Hardware (PCAN Channel) being used.
// If a parameter is not available, a PCAN_ERROR_ILLPARAMTYPE error will be returned
func GetValue(channel TPCANHandle, param TPCANParameter, buffer unsafe.Pointer, bufferSize uint32) (TPCANStatus, error) { // TODO change buffersize to uintptr
	ret, _, errCall := syscall.SyscallN(pHandleGetValue, uintptr(channel), uintptr(param), uintptr(buffer), uintptr(bufferSize))
	return TPCANStatus(ret), sysCallErr(errCall)
}

// Configures a PCAN Channel value.
// Channel: The handle of a PCAN Channel
// param: The TPCANParameter parameter to set
// value: Value of parameter
// Note: Parameters can be present or not according with the kind of Hardware (PCAN Channel) being used.
// If a parameter is not available, a PCAN_ERROR_ILLPARAMTYPE error will be returned
func SetValue(channel TPCANHandle, param TPCANParameter, buffer unsafe.Pointer, bufferSize uintptr) (TPCANStatus, error) {
	ret, _, errCall := syscall.SyscallN(pHandleSetValue, uintptr(channel), uintptr(param), uintptr(buffer), bufferSize)
	return TPCANStatus(ret), sysCallErr(errCall)
}

// Returns a descriptive text of a given TPCANStatus error code, in any desired language
// err: A TPCANStatus error code
// language: Indicates a 'Primary language ID'
func GetErrorText(status TPCANStatus, language TPCANLanguage) (TPCANStatus, [256]byte, error) {
	var buffer [256]byte

	ret, _, errCall := syscall.SyscallN(pHandleGetErrorText, uintptr(status), uintptr(language), uintptr(unsafe.Pointer(&buffer)))
	return TPCANStatus(ret), buffer, sysCallErr(errCall)
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

	ret, _, errCall := syscall.SyscallN(pHandleLookUpChannel, uintptr(unsafe.Pointer(&sParameters)), uintptr(unsafe.Pointer(&foundChannel)))
	return TPCANStatus(ret), foundChannel, sysCallErr(errCall)
}

// helper function to handle syscall return value
func sysCallErr(errCall syscall.Errno) error {
	var err error
	if errCall != 0 {
		err = errors.New(errCall.Error())
	}
	return err
}
