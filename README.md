# gocan

A Golang CAN Bus interface supporting different CAN device manufactures.
Supports Windows, Linux and MacOS.

## Supported Interfaces

- PEAK Systems PCAN

## Notes

- Error frames are only received on Recv() call as message if the config *"RecvErrorFrames"* is set
- Error frames are logged if the config *"LogErrorFrames"* is set

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

## Open TODOs

1. PCAN

- CAN FD still not implemented
- Implement GetAllAvailableHandles function: Maybe there is a convinient API function ?
- SetValue and GetValue methods are not tested

2. Other interfaces as vector, kvasor and others
