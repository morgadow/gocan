package pcan

import (
	"errors"
	"runtime"
	"syscall"
	"unsafe"
)

const LengthDataCANMessage = 8       // maximum amount of bytes in an PCAN CAN message
const MaxLengthDataCANFDMessage = 64 // maximum amount of bytes in can CAN FD message

const DriverWin = "PCANBasic.dll"     // PCAN windows driver dll name, which can be imported once the PCAN driver is installed
const DriverMac = "libPCBUSB.dylib"   // PCAN max driver file name, which can be imported once the PCAN driver is installed
const DriverLinux = "libpcanbasic.so" // PCAN linux driver file name, which can be imported once the PCAN driver is installed

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
var APILoaded bool

type TPCANHandle uint16      // Represents a PCAN hardware Channel handle
type TPCANStatus uint32      // Represents a PCAN status/error code
type TPCANParameter uint8    // Represents a PCAN parameter to be read or set
type TPCANDevice uint8       // Represents a PCAN device
type TPCANMessageType uint8  // Represents the type of a PCAN message
type TPCANType uint8         // Represents the type of PCAN hardware to be initialized
type TPCANMode uint8         // Represents a PCAN filter mode
type TPCANBaudrate uint16    // Represents a PCAN Baud rate register value
type TPCANBitrateFD string   // Represents a PCAN-FD bit rate string
type TPCANTimestampFD uint64 // Represents a timestamp of a received PCAN FD message
type TPCANMsgID uint32       // 11/29-bit message identifier
type TPCANLanguage uint16

const ( // Languages to select for error texts
	LanguageNeutral TPCANLanguage = 0x00
	LanguageGerman  TPCANLanguage = 0x07
	LanguageEnglish TPCANLanguage = 0x09
	LanguageItalian TPCANLanguage = 0x10
	LanguageSpanish TPCANLanguage = 0x0A
	LanguageFrench  TPCANLanguage = 0x0C
)

const ( // Currently defined and supported PCAN channels
	PCAN_NONEBUS = TPCANHandle(0x00) // Undefined/default value for a PCAN bus

	PCAN_ISABUS1 = TPCANHandle(0x21) // PCAN-ISA interface, Channel 1
	PCAN_ISABUS2 = TPCANHandle(0x22) // PCAN-ISA interface, Channel 2
	PCAN_ISABUS3 = TPCANHandle(0x23) // PCAN-ISA interface, Channel 3
	PCAN_ISABUS4 = TPCANHandle(0x24) // PCAN-ISA interface, Channel 4
	PCAN_ISABUS5 = TPCANHandle(0x25) // PCAN-ISA interface, Channel 5
	PCAN_ISABUS6 = TPCANHandle(0x26) // PCAN-ISA interface, Channel 6
	PCAN_ISABUS7 = TPCANHandle(0x27) // PCAN-ISA interface, Channel 7
	PCAN_ISABUS8 = TPCANHandle(0x28) // PCAN-ISA interface, Channel 8

	PCAN_DNGBUS1 = TPCANHandle(0x31) // PCAN-Dongle/LPT interface, Channel 1

	PCAN_PCIBUS1  = TPCANHandle(0x41)  // PCAN-PCI interface, Channel 1
	PCAN_PCIBUS2  = TPCANHandle(0x42)  // PCAN-PCI interface, Channel 2
	PCAN_PCIBUS3  = TPCANHandle(0x43)  // PCAN-PCI interface, Channel 3
	PCAN_PCIBUS4  = TPCANHandle(0x44)  // PCAN-PCI interface, Channel 4
	PCAN_PCIBUS5  = TPCANHandle(0x45)  // PCAN-PCI interface, Channel 5
	PCAN_PCIBUS6  = TPCANHandle(0x46)  // PCAN-PCI interface, Channel 6
	PCAN_PCIBUS7  = TPCANHandle(0x47)  // PCAN-PCI interface, Channel 7
	PCAN_PCIBUS8  = TPCANHandle(0x48)  // PCAN-PCI interface, Channel 8
	PCAN_PCIBUS9  = TPCANHandle(0x409) // PCAN-PCI interface, Channel 9
	PCAN_PCIBUS10 = TPCANHandle(0x40A) // PCAN-PCI interface, Channel 10
	PCAN_PCIBUS11 = TPCANHandle(0x40B) // PCAN-PCI interface, Channel 11
	PCAN_PCIBUS12 = TPCANHandle(0x40C) // PCAN-PCI interface, Channel 12
	PCAN_PCIBUS13 = TPCANHandle(0x40D) // PCAN-PCI interface, Channel 13
	PCAN_PCIBUS14 = TPCANHandle(0x40E) // PCAN-PCI interface, Channel 14
	PCAN_PCIBUS15 = TPCANHandle(0x40F) // PCAN-PCI interface, Channel 15
	PCAN_PCIBUS16 = TPCANHandle(0x410) // PCAN-PCI interface, Channel 16

	PCAN_USBBUS1  = TPCANHandle(0x51)  // PCAN-USB interface, Channel 1
	PCAN_USBBUS2  = TPCANHandle(0x52)  // PCAN-USB interface, Channel 2
	PCAN_USBBUS3  = TPCANHandle(0x53)  // PCAN-USB interface, Channel 3
	PCAN_USBBUS4  = TPCANHandle(0x54)  // PCAN-USB interface, Channel 4
	PCAN_USBBUS5  = TPCANHandle(0x55)  // PCAN-USB interface, Channel 5
	PCAN_USBBUS6  = TPCANHandle(0x56)  // PCAN-USB interface, Channel 6
	PCAN_USBBUS7  = TPCANHandle(0x57)  // PCAN-USB interface, Channel 7
	PCAN_USBBUS8  = TPCANHandle(0x58)  // PCAN-USB interface, Channel 8
	PCAN_USBBUS9  = TPCANHandle(0x509) // PCAN-USB interface, Channel 9
	PCAN_USBBUS10 = TPCANHandle(0x50A) // PCAN-USB interface, Channel 10
	PCAN_USBBUS11 = TPCANHandle(0x50B) // PCAN-USB interface, Channel 11
	PCAN_USBBUS12 = TPCANHandle(0x50C) // PCAN-USB interface, Channel 12
	PCAN_USBBUS13 = TPCANHandle(0x50D) // PCAN-USB interface, Channel 13
	PCAN_USBBUS14 = TPCANHandle(0x50E) // PCAN-USB interface, Channel 14
	PCAN_USBBUS15 = TPCANHandle(0x50F) // PCAN-USB interface, Channel 15
	PCAN_USBBUS16 = TPCANHandle(0x510) // PCAN-USB interface, Channel 16

	PCAN_PCCBUS1 = TPCANHandle(0x61) // PCAN-PC Card interface, Channel 1
	PCAN_PCCBUS2 = TPCANHandle(0x62) // PCAN-PC Card interface, Channel 2

	PCAN_LANBUS1  = TPCANHandle(0x801) // PCAN-LAN interface, Channel 1
	PCAN_LANBUS2  = TPCANHandle(0x802) // PCAN-LAN interface, Channel 2
	PCAN_LANBUS3  = TPCANHandle(0x803) // PCAN-LAN interface, Channel 3
	PCAN_LANBUS4  = TPCANHandle(0x804) // PCAN-LAN interface, Channel 4
	PCAN_LANBUS5  = TPCANHandle(0x805) // PCAN-LAN interface, Channel 5
	PCAN_LANBUS6  = TPCANHandle(0x806) // PCAN-LAN interface, Channel 6
	PCAN_LANBUS7  = TPCANHandle(0x807) // PCAN-LAN interface, Channel 7
	PCAN_LANBUS8  = TPCANHandle(0x808) // PCAN-LAN interface, Channel 8
	PCAN_LANBUS9  = TPCANHandle(0x809) // PCAN-LAN interface, Channel 9
	PCAN_LANBUS10 = TPCANHandle(0x80A) // PCAN-LAN interface, Channel 10
	PCAN_LANBUS11 = TPCANHandle(0x80B) // PCAN-LAN interface, Channel 11
	PCAN_LANBUS12 = TPCANHandle(0x80C) // PCAN-LAN interface, Channel 12
	PCAN_LANBUS13 = TPCANHandle(0x80D) // PCAN-LAN interface, Channel 13
	PCAN_LANBUS14 = TPCANHandle(0x80E) // PCAN-LAN interface, Channel 14
	PCAN_LANBUS15 = TPCANHandle(0x80F) // PCAN-LAN interface, Channel 15
	PCAN_LANBUS16 = TPCANHandle(0x810) // PCAN-LAN interface, Channel 16
)

const ( // Represent the PCAN error and status codes
	PCAN_ERROR_OK           = TPCANStatus(0x00000)                                                                                          // No error
	PCAN_ERROR_XMTFULL      = TPCANStatus(0x00001)                                                                                          // Transmit buffer in CAN controller is full
	PCAN_ERROR_OVERRUN      = TPCANStatus(0x00002)                                                                                          // CAN controller was read too late
	PCAN_ERROR_BUSLIGHT     = TPCANStatus(0x00004)                                                                                          // BusIf error: an error counter reached the 'light' limit
	PCAN_ERROR_BUSHEAVY     = TPCANStatus(0x00008)                                                                                          // BusIf error: an error counter reached the 'heavy' limit
	PCAN_ERROR_BUSWARNING   = PCAN_ERROR_BUSHEAVY                                                                                           // BusIf error: an error counter reached the 'warning' limit
	PCAN_ERROR_BUSPASSIVE   = TPCANStatus(0x40000)                                                                                          // BusIf error: the CAN controller is error passive
	PCAN_ERROR_BUSOFF       = TPCANStatus(0x00010)                                                                                          // BusIf error: the CAN controller is in bus-off state
	PCAN_ERROR_ANYBUSERR    = PCAN_ERROR_BUSWARNING | PCAN_ERROR_BUSLIGHT | PCAN_ERROR_BUSHEAVY | PCAN_ERROR_BUSOFF | PCAN_ERROR_BUSPASSIVE // Mask for all bus errors
	PCAN_ERROR_QRCVEMPTY    = TPCANStatus(0x00020)                                                                                          // Receive queue is empty
	PCAN_ERROR_QOVERRUN     = TPCANStatus(0x00040)                                                                                          // Receive queue was read too late
	PCAN_ERROR_QXMTFULL     = TPCANStatus(0x00080)                                                                                          // Transmit queue is full
	PCAN_ERROR_REGTEST      = TPCANStatus(0x00100)                                                                                          // Test of the CAN controller hardware registers failed (no hardware found)
	PCAN_ERROR_NODRIVER     = TPCANStatus(0x00200)                                                                                          // Driver not loaded
	PCAN_ERROR_HWINUSE      = TPCANStatus(0x00400)                                                                                          // Hardware already in use by a Net
	PCAN_ERROR_NETINUSE     = TPCANStatus(0x00800)                                                                                          // A Client is already connected to the Net
	PCAN_ERROR_ILLHW        = TPCANStatus(0x01400)                                                                                          // Hardware handle is invalid
	PCAN_ERROR_ILLNET       = TPCANStatus(0x01800)                                                                                          // Net handle is invalid
	PCAN_ERROR_ILLCLIENT    = TPCANStatus(0x01C00)                                                                                          // Client handle is invalid
	PCAN_ERROR_ILLHANDLE    = PCAN_ERROR_ILLHW | PCAN_ERROR_ILLNET | PCAN_ERROR_ILLCLIENT                                                   // Mask for all handle errors
	PCAN_ERROR_RESOURCE     = TPCANStatus(0x02000)                                                                                          // Resource (FIFO, Client, timeout) cannot be created
	PCAN_ERROR_ILLPARAMTYPE = TPCANStatus(0x04000)                                                                                          // Invalid parameter
	PCAN_ERROR_ILLPARAMVAL  = TPCANStatus(0x08000)                                                                                          // Invalid parameter value
	PCAN_ERROR_UNKNOWN      = TPCANStatus(0x10000)                                                                                          // Unknown error
	PCAN_ERROR_ILLDATA      = TPCANStatus(0x20000)                                                                                          // Invalid data, function, or action
	PCAN_ERROR_CAUTION      = TPCANStatus(0x2000000)                                                                                        // An operation was successfully carried out, however, irregularities were registered
	PCAN_ERROR_INITIALIZE   = TPCANStatus(0x4000000)                                                                                        // Channel is not initialized [Value was changed from 0x40000 to 0x4000000]
	PCAN_ERROR_ILLOPERATION = TPCANStatus(0x8000000)                                                                                        // Invalid operation [Value was changed from 0x80000 to 0x8000000]
)

const ( // PCAN devices
	PCAN_NONE    = TPCANDevice(0x00) // Undefined, unknown or not selected PCAN device value
	PCAN_PEAKCAN = TPCANDevice(0x01) // PCAN Non-Plug&Play devices. NOT USED WITHIN PCAN-Basic API
	PCAN_ISA     = TPCANDevice(0x02) // PCAN-ISA, PCAN-PC/104, and PCAN-PC/104-Plus
	PCAN_DNG     = TPCANDevice(0x03) // PCAN-Dongle
	PCAN_PCI     = TPCANDevice(0x04) // PCAN-PCI, PCAN-cPCI, PCAN-miniPCI, and PCAN-PCI Express
	PCAN_USB     = TPCANDevice(0x05) // PCAN-USB and PCAN-USB Pro
	PCAN_PCC     = TPCANDevice(0x06) // PCAN-PC Card
	PCAN_VIRTUAL = TPCANDevice(0x07) // PCAN Virtual hardware. NOT USED WITHIN PCAN-Basic API
	PCAN_LAN     = TPCANDevice(0x08) // PCAN Gateway devices
)

const ( // PCAN parameters
	PCAN_DEVICE_NUMBER            = TPCANParameter(0x01) // PCAN-USB device number parameter
	PCAN_5VOLTS_POWER             = TPCANParameter(0x02) // PCAN-PC Card 5-Volt power parameter
	PCAN_RECEIVE_EVENT            = TPCANParameter(0x03) // PCAN receive event handler parameter
	PCAN_MESSAGE_FILTER           = TPCANParameter(0x04) // PCAN message filter parameter
	PCAN_API_VERSION              = TPCANParameter(0x05) // PCAN-Basic API version parameter
	PCAN_CHANNEL_VERSION          = TPCANParameter(0x06) // PCAN device Channel version parameter
	PCAN_BUSOFF_AUTORESET         = TPCANParameter(0x07) // PCAN Reset-On-Busoff parameter
	PCAN_LISTEN_ONLY              = TPCANParameter(0x08) // PCAN Listen-Only parameter
	PCAN_LOG_LOCATION             = TPCANParameter(0x09) // Directory path for log files
	PCAN_LOG_STATUS               = TPCANParameter(0x0A) // Debug-Log activation status
	PCAN_LOG_CONFIGURE            = TPCANParameter(0x0B) // Configuration of the debugged information (LOG_FUNCTION_***)
	PCAN_LOG_TEXT                 = TPCANParameter(0x0C) // Custom insertion of text into the log file
	PCAN_CHANNEL_CONDITION        = TPCANParameter(0x0D) // Availability status of a PCAN-Channel
	PCAN_HARDWARE_NAME            = TPCANParameter(0x0E) // PCAN hardware name parameter
	PCAN_RECEIVE_STATUS           = TPCANParameter(0x0F) // Message reception status of a PCAN-Channel
	PCAN_CONTROLLER_NUMBER        = TPCANParameter(0x10) // CAN-Controller number of a PCAN-Channel
	PCAN_TRACE_LOCATION           = TPCANParameter(0x11) // Directory path for PCAN trace files
	PCAN_TRACE_STATUS             = TPCANParameter(0x12) // CAN tracing activation status
	PCAN_TRACE_SIZE               = TPCANParameter(0x13) // Configuration of the maximum file size of a CAN trace
	PCAN_TRACE_CONFIGURE          = TPCANParameter(0x14) // Configuration of the trace file storing mode (TRACE_FILE_***)
	PCAN_CHANNEL_IDENTIFYING      = TPCANParameter(0x15) // Physical identification of a USB based PCAN-Channel by blinking its associated LED
	PCAN_CHANNEL_FEATURES         = TPCANParameter(0x16) // Capabilities of a PCAN device (FEATURE_***)
	PCAN_BITRATE_ADAPTING         = TPCANParameter(0x17) // Using of an existing bit rate (PCAN-View connected to a Channel)
	PCAN_BITRATE_INFO             = TPCANParameter(0x18) // Configured bit rate as Btr0Btr1 value
	PCAN_BITRATE_INFO_FD          = TPCANParameter(0x19) // Configured bit rate as TPCANBitrateFD string
	PCAN_BUSSPEED_NOMINAL         = TPCANParameter(0x1A) // Configured nominal CAN BusIf speed as Bits per seconds
	PCAN_BUSSPEED_DATA            = TPCANParameter(0x1B) // Configured CAN data speed as Bits per seconds
	PCAN_IP_ADDRESS               = TPCANParameter(0x1C) // Remote address of a LAN Channel as string in IPv4 format
	PCAN_LAN_SERVICE_STATUS       = TPCANParameter(0x1D) // Status of the Virtual PCAN-Gateway Service
	PCAN_ALLOW_STATUS_FRAMES      = TPCANParameter(0x1E) // Status messages reception status within a PCAN-Channel
	PCAN_ALLOW_RTR_FRAMES         = TPCANParameter(0x1F) // RTR messages reception status within a PCAN-Channel
	PCAN_ALLOW_ERROR_FRAMES       = TPCANParameter(0x20) // Error messages reception status within a PCAN-Channel
	PCAN_INTERFRAME_DELAY         = TPCANParameter(0x21) // Delay, in microseconds, between sending frames
	PCAN_ACCEPTANCE_FILTER_11BIT  = TPCANParameter(0x22) // Filter over code and mask patterns for 11-Bit messages
	PCAN_ACCEPTANCE_FILTER_29BIT  = TPCANParameter(0x23) // Filter over code and mask patterns for 29-Bit messages
	PCAN_IO_DIGITAL_CONFIGURATION = TPCANParameter(0x24) // Output mode of 32 digital I/O pin of a PCAN-USB Chip. 1: Output-Active 0 : Output Inactive
	PCAN_IO_DIGITAL_VALUE         = TPCANParameter(0x25) // Value assigned to a 32 digital I/O pins of a PCAN-USB Chip
	PCAN_IO_DIGITAL_SET           = TPCANParameter(0x26) // Value assigned to a 32 digital I/O pins of a PCAN-USB Chip - Multiple digital I/O pins to 1 = High
	PCAN_IO_DIGITAL_CLEAR         = TPCANParameter(0x27) // Clear multiple digital I/O pins to 0
	PCAN_IO_ANALOG_VALUE          = TPCANParameter(0x28) // Get value of a single analog input pin
)

const ( // PCAN parameter values
	PCAN_PARAMETER_OFF       = int(0x00)                                      // The PCAN parameter is not set (inactive)
	PCAN_PARAMETER_ON        = int(0x01)                                      // The PCAN parameter is set (active)
	PCAN_FILTER_CLOSE        = int(0x00)                                      // The PCAN filter is closed. No messages will be received
	PCAN_FILTER_OPEN         = int(0x01)                                      // The PCAN filter is fully opened. All messages will be received
	PCAN_FILTER_CUSTOM       = int(0x02)                                      // The PCAN filter is custom configured. Only registered messages will be received
	PCAN_CHANNEL_UNAVAILABLE = int(0x00)                                      // The PCAN-Channel handle is illegal, or its associated hardware is not available
	PCAN_CHANNEL_AVAILABLE   = int(0x01)                                      // The PCAN-Channel handle is available to be connected (Plug&Play Hardware: it means furthermore that the hardware is plugged-in)
	PCAN_CHANNEL_OCCUPIED    = int(0x02)                                      // The PCAN-Channel handle is valid, and is already being used
	PCAN_CHANNEL_PCANVIEW    = PCAN_CHANNEL_AVAILABLE | PCAN_CHANNEL_OCCUPIED // The PCAN-Channel handle is already being used by a PCAN-View application, but is available to connect

	LOG_FUNCTION_DEFAULT    = int(0x00)   // Logs system exceptions / errors
	LOG_FUNCTION_ENTRY      = int(0x01)   // Logs the entries to the PCAN-Basic API functions
	LOG_FUNCTION_PARAMETERS = int(0x02)   // Logs the parameters passed to the PCAN-Basic API functions
	LOG_FUNCTION_LEAVE      = int(0x04)   // Logs the exits from the PCAN-Basic API functions
	LOG_FUNCTION_WRITE      = int(0x08)   // Logs the CAN messages passed to the CAN_Write function
	LOG_FUNCTION_READ       = int(0x10)   // Logs the CAN messages received within the CAN_Read function
	LOG_FUNCTION_ALL        = int(0xFFFF) // Logs all possible information within the PCAN-Basic API functions

	TRACE_FILE_SINGLE    = int(0x00) // A single file is written until it size reaches PAN_TRACE_SIZE
	TRACE_FILE_SEGMENTED = int(0x01) // Traced data is distributed in several files with size PAN_TRACE_SIZE
	TRACE_FILE_DATE      = int(0x02) // Includes the date into the name of the trace file
	TRACE_FILE_TIME      = int(0x04) // Includes the start time into the name of the trace file
	TRACE_FILE_OVERWRITE = int(0x80) // Causes the overwriting of available traces (same name)

	FEATURE_FD_CAPABLE    = int(0x01) // Device supports flexible data-rate (CAN-FD)
	FEATURE_DELAY_CAPABLE = int(0x02) // Device supports a delay between sending frames (FPGA based USB devices)
	FEATURE_IO_CAPABLE    = int(0x04) // Device supports I/O functionality for electronic circuits (USB-Chip devices)

	SERVICE_STATUS_STOPPED = int(0x01) // The service is not running
	SERVICE_STATUS_RUNNING = int(0x04) // The service is running
)

const ( // PCAN message types
	PCAN_MESSAGE_STANDARD = TPCANMessageType(0x00) // The PCAN message is a CAN Standard Frame (11-bit identifier)
	PCAN_MESSAGE_RTR      = TPCANMessageType(0x01) // The PCAN message is a CAN Remote-Transfer-Request Frame
	PCAN_MESSAGE_EXTENDED = TPCANMessageType(0x02) // The PCAN message is a CAN Extended Frame (29-bit identifier)
	PCAN_MESSAGE_FD       = TPCANMessageType(0x04) // The PCAN message represents a FD frame in terms of CiA Specs
	PCAN_MESSAGE_BRS      = TPCANMessageType(0x08) // The PCAN message represents a FD bit rate switch (CAN data at a higher bit rate)
	PCAN_MESSAGE_ESI      = TPCANMessageType(0x10) // The PCAN message represents a FD error state indicator(CAN FD transmitter was error active)
	PCAN_MESSAGE_ERRFRAME = TPCANMessageType(0x40) // The PCAN message represents an error frame
	PCAN_MESSAGE_STATUS   = TPCANMessageType(0x80) // The PCAN message represents a PCAN status message
)

const ( // Frame Type / Initialization Mode
	PCAN_MODE_STANDARD = TPCANMode(PCAN_MESSAGE_STANDARD)
	PCAN_MODE_EXTENDED = TPCANMode(PCAN_MESSAGE_EXTENDED)
)

const ( // Baud rate codes = BTR0/BTR1 register values for the CAN controller
	PCAN_BAUD_1M   = TPCANBaudrate(0x0014) //   1     MBit/s
	PCAN_BAUD_800K = TPCANBaudrate(0x0016) // 800     kBit/s
	PCAN_BAUD_500K = TPCANBaudrate(0x001C) // 500     kBit/s
	PCAN_BAUD_250K = TPCANBaudrate(0x011C) // 250     kBit/s
	PCAN_BAUD_125K = TPCANBaudrate(0x031C) // 125     kBit/s
	PCAN_BAUD_100K = TPCANBaudrate(0x432F) // 100     kBit/s
	PCAN_BAUD_95K  = TPCANBaudrate(0xC34E) //  95,238 kBit/s
	PCAN_BAUD_83K  = TPCANBaudrate(0x852B) //  83,333 kBit/s
	PCAN_BAUD_50K  = TPCANBaudrate(0x472F) //  50     kBit/s
	PCAN_BAUD_47K  = TPCANBaudrate(0x1414) //  47,619 kBit/s
	PCAN_BAUD_33K  = TPCANBaudrate(0x8B2F) //  33,333 kBit/s
	PCAN_BAUD_20K  = TPCANBaudrate(0x532F) //  20     kBit/s
	PCAN_BAUD_10K  = TPCANBaudrate(0x672F) //  10     kBit/s
	PCAN_BAUD_5K   = TPCANBaudrate(0x7F7F) //   5     kBit/s
)

const (
	/*
		# Represents the configuration for a CAN bit rate
		# Note:
		#    * Each parameter and its value must be separated with a '='.
		#    * Each pair of parameter/value must be separated using ','.
		#
		# Example:
		#    f_clock=80000000,nom_brp=10,nom_tseg1=5,nom_tseg2=2,nom_sjw=1,data_brp=4,data_tseg1=7,data_tseg2=2,data_sjw=1
	*/
	PCAN_BR_CLOCK       = TPCANBitrateFD("f_clock")
	PCAN_BR_CLOCK_MHZ   = TPCANBitrateFD("f_clock_mhz")
	PCAN_BR_NOM_BRP     = TPCANBitrateFD("nom_brp")
	PCAN_BR_NOM_TSEG1   = TPCANBitrateFD("nom_tseg1")
	PCAN_BR_NOM_TSEG2   = TPCANBitrateFD("nom_tseg2")
	PCAN_BR_NOM_SJW     = TPCANBitrateFD("nom_sjw")
	PCAN_BR_NOM_SAMPLE  = TPCANBitrateFD("nom_sam")
	PCAN_BR_DATA_BRP    = TPCANBitrateFD("data_brp")
	PCAN_BR_DATA_TSEG1  = TPCANBitrateFD("data_tseg1")
	PCAN_BR_DATA_TSEG2  = TPCANBitrateFD("data_tseg2")
	PCAN_BR_DATA_SJW    = TPCANBitrateFD("data_sjw")
	PCAN_BR_DATA_SAMPLE = TPCANBitrateFD("data_ssp_offset")
)

const ( // Supported No-Plug-And-Play Hardware types
	PCAN_TYPE_ISA         = TPCANType(0x01) // PCAN-ISA 82C200
	PCAN_TYPE_ISA_SJA     = TPCANType(0x09) // PCAN-ISA SJA1000
	PCAN_TYPE_ISA_PHYTEC  = TPCANType(0x04) // PHYTEC ISA
	PCAN_TYPE_DNG         = TPCANType(0x02) // PCAN-Dongle 82C200
	PCAN_TYPE_DNG_EPP     = TPCANType(0x03) // PCAN-Dongle EPP 82C200
	PCAN_TYPE_DNG_SJA     = TPCANType(0x05) // PCAN-Dongle SJA1000
	PCAN_TYPE_DNG_SJA_EPP = TPCANType(0x06) // PCAN-Dongle EPP SJA1000
)

// TPCANMessage Represents a PCAN message
type TPCANMessage struct {
	ID      TPCANMsgID                 // 11/29-bit message identifier
	MsgType TPCANMessageType           // Type of the message
	DLC     uint8                      // Data Length Code of the message (0..8)
	Data    [LengthDataCANMessage]byte // Data of the message (DATA[0]..DATA[7])
}

// TPCANTimestamp Represents a timestamp of a received PCAN message
// Total Microseconds = micros + 1000 * millis + 0x100000000 * 1000 * millis_overflow
type TPCANTimestamp struct {
	Millis         uint32 // Base-value: milliseconds: 0.. 2^32-1
	MillisOverflow uint16 // Roll-arounds of millis
	Micros         uint16 // Microseconds: 0..999
}

// TPCANMessageFD Represents a PCAN CAN FD message
type TPCANMessageFD struct {
	ID      TPCANMsgID
	MsgType TPCANMessageType
	DLC     uint8
	Data    []byte
}

// LoadAPI Loads PCAN API (.ddl) file
func LoadAPI() error {
	var err error = nil
	var dll = ""

	// evaluate operating system and architecture and select driver file
	switch runtime.GOOS {
	case "windows":
		dll = DriverWin
	case "darwin":
		dll = DriverMac
	default:
		dll = DriverLinux
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

	APILoaded = pHandleInitialize > 0 && pHandleInitializeFD > 0 && pHandleReset > 0 && pHandleGetStatus > 0 &&
		pHandleRead > 0 && pHandleReadFD > 0 && pHandleWrite > 0 && pHandleWriteFD > 0 && pHandleFilterMessages > 0 && pHandleGetValue > 0 &&
		pHandleSetValue > 0 && pHandleGetErrorText > 0 && pHandleLookUpChannel > 0 && pHandleUninitialize > 0

	if !APILoaded {
		return errors.New("could not load pointers to pcan functions")
	}
	return nil
}

// UnloadAPI Unloads PCAN API (.ddl) file
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
	APILoaded = false

	err := syscall.FreeLibrary(pcanAPIHandle)
	return err
}

// Initialize Initializes a PCAN Channel
// Channel: The handle of a PCAN Channel
// baudRate: The speed for the communication (BTR0BTR1 code)
// hwType: Non-PnP: The type of hardware and operation mode
// ioPort: Non-PnP: The I/O address for the parallel port
// interrupt: Non-PnP: Interrupt number of the parallel por
func Initialize(channel TPCANHandle, baudRate TPCANBaudrate, hwType TPCANType, ioPort uint32, interrupt uint16) (TPCANStatus, error) {

	var nargs uintptr = 0
	var err error = nil

	if !APILoaded {
		return PCAN_ERROR_UNKNOWN, ErrAPINotLoadedOrFound
	}

	ret, _, errCall := syscall.Syscall6(pHandleInitialize, nargs, uintptr(channel), uintptr(baudRate), uintptr(hwType), uintptr(ioPort), uintptr(interrupt), 0)
	if errCall != 0 {
		err = errors.New(errCall.Error())
	}
	return TPCANStatus(ret), err
}

// InitializeFD Initializes a FD capable PCAN Channel
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

	var nargs uintptr = 0
	var err error = nil

	if !APILoaded {
		return PCAN_ERROR_UNKNOWN, ErrAPINotLoadedOrFound
	}

	ret, _, errCall := syscall.Syscall(pHandleInitializeFD, nargs, uintptr(channel), uintptr(unsafe.Pointer(&bitRateFD)), 0)
	if errCall != 0 {
		err = errors.New(errCall.Error())
	}
	return TPCANStatus(ret), err
}

// Uninitialize Uninitializes PCAN Channels initialized by CAN_Initialize
// Channel: The handle of a PCAN Channel
func Uninitialize(channel TPCANHandle) (TPCANStatus, error) {

	var nargs uintptr = 0
	var err error = nil

	if !APILoaded {
		return PCAN_ERROR_UNKNOWN, ErrAPINotLoadedOrFound
	}

	ret, _, errCall := syscall.Syscall(pHandleUninitialize, nargs, uintptr(channel), 0, 0)
	if errCall != 0 {
		err = errors.New(errCall.Error())
	}
	return TPCANStatus(ret), err
}

// Reset Resets the receive and transmit queues of the PCAN Channel
// Channel: The handle of a PCAN Channel
func Reset(channel TPCANHandle) (TPCANStatus, error) {

	var nargs uintptr = 0
	var err error = nil

	if !APILoaded {
		return PCAN_ERROR_UNKNOWN, ErrAPINotLoadedOrFound
	}

	ret, _, errCall := syscall.Syscall(pHandleReset, nargs, uintptr(channel), 0, 0)
	if errCall != 0 {
		err = errors.New(errCall.Error())
	}
	return TPCANStatus(ret), err
}

// GetStatus Gets the current status of a PCAN Channel
// Channel: The handle of a PCAN Channel
func GetStatus(channel TPCANHandle) (TPCANStatus, error) {

	var nargs uintptr = 0
	var err error = nil

	if !APILoaded {
		return PCAN_ERROR_UNKNOWN, ErrAPINotLoadedOrFound
	}

	ret, _, errCall := syscall.Syscall(pHandleGetStatus, nargs, uintptr(channel), 0, 0)
	if errCall != 0 {
		err = errors.New(errCall.Error())
	}
	return TPCANStatus(ret), err
}

// Read Reads a CAN message from the receive queue of a PCAN Channel
// Channel: The handle of a PCAN Channel
func Read(channel TPCANHandle) (TPCANStatus, TPCANMessage, TPCANTimestamp, error) {

	var msg TPCANMessage
	var timeStamp = TPCANTimestamp{}
	var nargs uintptr = 0
	var err error = nil

	if !APILoaded {
		return PCAN_ERROR_UNKNOWN, msg, timeStamp, ErrAPINotLoadedOrFound
	}

	ret, _, errCall := syscall.Syscall(pHandleRead, nargs, uintptr(channel), uintptr(unsafe.Pointer(&msg)), uintptr(unsafe.Pointer(&timeStamp)))
	if errCall != 0 {
		err = errors.New(errCall.Error())
	}
	return TPCANStatus(ret), msg, timeStamp, err
}

// ReadFD Reads a CAN message from the receive queue of a FD capable PCAN Channel
// Channel: The handle of a PCAN Channel
func ReadFD(channel TPCANHandle) (TPCANStatus, TPCANMessageFD, TPCANTimestampFD, error) {

	var msgFD TPCANMessageFD
	var timeStampFD TPCANTimestampFD
	var nargs uintptr = 0
	var err error = nil

	if !APILoaded {
		return PCAN_ERROR_UNKNOWN, msgFD, timeStampFD, ErrAPINotLoadedOrFound
	}

	ret, _, errCall := syscall.Syscall(pHandleReadFD, nargs, uintptr(channel), uintptr(unsafe.Pointer(&msgFD)), uintptr(unsafe.Pointer(&timeStampFD)))
	if errCall != 0 {
		err = errors.New(errCall.Error())
	}
	return TPCANStatus(ret), msgFD, timeStampFD, err
}

// Write Transmits a CAN message
// Channel: The handle of a PCAN Channel
// msg: A Message struct with the message to be sent
func Write(channel TPCANHandle, msg TPCANMessage) (TPCANStatus, error) {

	var nargs uintptr = 0
	var err error = nil

	if !APILoaded {
		return PCAN_ERROR_UNKNOWN, ErrAPINotLoadedOrFound
	}

	ret, _, errCall := syscall.Syscall(pHandleWrite, nargs, uintptr(channel), uintptr(unsafe.Pointer(&msg)), 0)
	if errCall != 0 {
		err = errors.New(errCall.Error())
	}
	return TPCANStatus(ret), err
}

// WriteFD Transmits a CAN message over a FD capable PCAN Channel
// Channel: The handle of a PCAN Channel
// msgFD A MessageFD struct with the message to be sent
func WriteFD(channel TPCANHandle, msgFD TPCANMessageFD) (TPCANStatus, error) {

	var nargs uintptr = 0
	var err error = nil

	if !APILoaded {
		return PCAN_ERROR_UNKNOWN, ErrAPINotLoadedOrFound
	}

	ret, _, errCall := syscall.Syscall(pHandleWriteFD, nargs, uintptr(channel), uintptr(unsafe.Pointer(&msgFD)), 0)
	if errCall != 0 {
		err = errors.New(errCall.Error())
	}
	return TPCANStatus(ret), err
}

// FilterMessages Configures the reception filter
// Channel: The handle of a PCAN Channel
// fromID: The lowest CAN ID to be received
// toID: The highest CAN ID to be received
// mode: Message type, Standard (11-bit identifier) or Extended (29-bit identifier)
func FilterMessages(channel TPCANHandle, fromID TPCANMsgID, toID TPCANMsgID, mode TPCANMode) (TPCANStatus, error) {

	var nargs uintptr = 0
	var err error = nil

	if !APILoaded {
		return PCAN_ERROR_UNKNOWN, ErrAPINotLoadedOrFound
	}

	ret, _, errCall := syscall.Syscall6(pHandleFilterMessages, nargs, uintptr(channel), uintptr(fromID), uintptr(toID), uintptr(mode), 0, 0)
	if errCall != 0 {
		err = errors.New(errCall.Error())
	}
	return TPCANStatus(ret), err
}

// GetValue Retrieves a PCAN Channel value
// Channel: The handle of a PCAN Channel
// param: The TPCANParameter parameter to get
// Note: Parameters can be present or not according with the kind
//
//	of Hardware (PCAN Channel) being used. If a parameter is not available,
//	a PCAN_ERROR_ILLPARAMTYPE error will be returned
func GetValue(channel TPCANHandle, param TPCANParameter) (TPCANStatus, uint32, error) {

	var nargs uintptr = 0
	var err error = nil
	var value uint32 = 0

	if !APILoaded {
		return PCAN_ERROR_UNKNOWN, value, ErrAPINotLoadedOrFound
	}

	ret, _, errCall := syscall.Syscall6(pHandleGetValue, nargs, uintptr(channel), uintptr(param), uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(value), 0, 0)
	if errCall != 0 {
		err = errors.New(errCall.Error())
	}
	return TPCANStatus(ret), value, err
}

// SetValue Configures a PCAN Channel value.
// Channel: The handle of a PCAN Channel
// param: The TPCANParameter parameter to set
// value: Value of parameter
// Note: Parameters can be present or not according with the kind
//
//	of Hardware (PCAN Channel) being used. If a parameter is not available,
//	a PCAN_ERROR_ILLPARAMTYPE error will be returned.
func SetValue(channel TPCANHandle, param TPCANParameter, value uint32) (TPCANStatus, error) { // todo can this always be uint32? I think a string value is also possible?!

	var nargs uintptr = 0
	var err error = nil

	if !APILoaded {
		return PCAN_ERROR_UNKNOWN, ErrAPINotLoadedOrFound
	}

	ret, _, errCall := syscall.Syscall6(pHandleSetValue, nargs, uintptr(channel), uintptr(param), uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(value), 0, 0)
	if errCall != 0 {
		err = errors.New(errCall.Error())
	}
	return TPCANStatus(ret), err
}

// GetErrorText Returns a descriptive text of a given TPCANStatus error code, in any desired language
// err: A TPCANStatus error code
// language: Indicates a 'Primary language ID'
func GetErrorText(status TPCANStatus, language TPCANLanguage) (TPCANStatus, [256]byte, error) {
	var nargs uintptr = 0
	var err error = nil
	var buffer [256]byte

	if !APILoaded {
		return PCAN_ERROR_UNKNOWN, buffer, ErrAPINotLoadedOrFound
	}

	ret, _, errCall := syscall.Syscall(pHandleGetErrorText, nargs, uintptr(status), uintptr(language), uintptr(unsafe.Pointer(&buffer)))
	if errCall != 0 {
		err = errors.New(errCall.Error())
	}
	return TPCANStatus(ret), buffer, err
}

// LookUpChannel Finds a PCAN-Basic Channel that matches with the given parameters
// parameters: A comma separated string contained pairs of parameter-name/value to be matched within a PCAN-Basic Channel
// foundChannels: Buffer for returning the PCAN-Basic Channel when found
func LookUpChannel(parameters string, foundChannel []TPCANHandle) (TPCANStatus, error) {

	var nargs uintptr = 0
	var err error = nil

	if !APILoaded {
		return PCAN_ERROR_UNKNOWN, ErrAPINotLoadedOrFound
	}

	ret, _, errCall := syscall.Syscall(pHandleLookUpChannel, nargs, uintptr(unsafe.Pointer(&parameters)), uintptr(unsafe.Pointer(&foundChannel)), 0)
	if errCall != 0 {
		err = errors.New(errCall.Error())
	}
	return TPCANStatus(ret), err
}

// UninitializeAllChannels Uninitializes all PCAN Channels initialized by CAN_Initialize
func UninitializeAllChannels() (TPCANStatus, error) {
	return Uninitialize(PCAN_NONEBUS)
}
