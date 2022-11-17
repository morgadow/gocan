# gocan

A Golang CAN Bus interface supporting different CAN device manufactures.
Supports Windows, Linux and MacOS.

```golang
// Bus Interface for all main CANBus functionality. Lower device interfaces may support more functionality
type Bus interface {
 Send(*Message) error
 Recv(timeout int) (*Message, error)
 StatusIsOkay() (bool, error)
 Status() (uint32, error)
 State() BusState
 ReadBuffer(limit uint16) ([]Message, error)
 SetFilter(fromID MessageID, toID MessageID, mode uint8) error
 Reset() error
 Shutdown() error
}
```

## Supported Interfaces

- PEAK Systems PCAN

### Examples

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

## Open TODOs

- option to create listener for CAN bus connections with different triggers for messages like message id masks, ...
- PCAN: missing CANFD implementation
- PCAN: CAN Epoch not set to correct PCAN timestamp, currently set to boot time
- PCAN: SetValue and GetValue methods are not tested and does not support all types
- Bus: implement GetAllAvailableHandles function. Maybe there is a convinient API function for devices?
- every other manufacturer
