# gocan

A Golang CAN Bus interface supporting different CAN device manufactures.
Supports Windows, Linux and MacOS.

```golang
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
 TraceStart(filePath string, maxFileSize int) error            // Starts recording a trace on given path
 TraceStop() error                                             // Stops recording currently running trace
}
```

## Interfaces

### PEAK Systems

```golang

 // Create CAN bus connection with configuration
 pcanBus, err := factory.CreateBus(&config)
 if err != nil {
  fmt.Printf(err.Error())
  return
 }

 // Receive message over bus
 rxMsg, err := pcanBus.Recv(60)
 if err != nil {
  fmt.Printf(err.Error())
 }
 if rxMsg != nil {
  fmt.Printf("\nMsg ID: %v, Msg DLC: %v, Msg Data: %v", rxMsg.ID, rxMsg.DLC, rxMsg.Data)
 }

 // Send a message over the bus
 txMsg := bus.Message{
  ID:   0x123,
  DLC:  8,
  Data: []uint8{1, 2, 3, 4, 5, 6, 7, 8}}

 err = pcanBus.Send(&txMsg)
 if err != nil {
  fmt.Printf(err.Error())
 }

 // Check if bus is okay
 ok, err := pcanBus.StatusIsOkay()
 if err != nil {
  fmt.Printf(err.Error())
 }
 fmt.Printf("\nBus okay: %v", ok)

 // Get bus status
 status, err := pcanBus.Status()
 if err != nil {
  fmt.Printf(err.Error())
 }
 fmt.Printf("\nStatus: %v", status)

 // Reset CAN in case of an invalid bus state
 err = pcanBus.Reset()
 if err != nil {
  fmt.Printf(err.Error())
 }

 // Apply filter to channel
 err = pcanBus.SetFilter(0x000100, 0x000200, uint8(pcan.PCAN_MODE_EXTENDED))
 if err != nil {
  fmt.Printf(err.Error())
 }

 // Shutdown channel again
 err = pcanBus.Shutdown()
 if err != nil {
  fmt.Printf(err.Error())
 }

```

## Changelog

- v1.0.0:
  - initial working version implementing project structure for sending and receiving messages over a CAN bus connection
  - fully implemented *pcan* can interface

- v1.0.1:
  - message DLC now automatically evaluated when sending a message, before this the DLC must be set manually to an non zero value
  - updated README with an changelog section

- v1.1.0:
  - internal restructure for PCAN interface
  - added tracing related functions to gocan bus interface and implemented functionality in PCAN interface
    - added TraceStart(filePath string, maxFileSize int) error:
    - added TraceStop() error
  - changed license from MIT to GPLv3
  - added factory function for detecting available handles
  - changed config struct for more receive filtering options and implemented direct dll support for this settings in PCAN
  - changed internal call convention for PCAN driver
  - changed pcan handle initialization to be only for plug n play devices on gocan interface, old variant still usable for direct pcan interface call

## Known Issues

This section lists all known-issues, missing features and open bugs.

### Interfaces

#### PCAN

- missing documentation examples for new functions
- Invalid buffer size error in LookupChannel function
- error FILE_NOT_FOUND when calling the Shutdown or Uninitialize function: problem probably located in .dll file itself
- FilterMessages function not working correctly; does not apply any filter to the PCAN channel
- Missing implementation of any further filter option as message masks
- Missing implementation of CANFD functionality due to missing test hardware
- Evaluation of channel condition propably incorrect as every connection is marked as unavailable
- Setting parameter as the PCAN_READ_ONLY does not have an impact, reading of message is still possible
- Recorded trace files maintain empty even there is traffic on selected CAN bus
