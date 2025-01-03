package main

import (
	"fmt"

	"github.com/morgadow/gocan"
	"github.com/morgadow/gocan/factory"
)

func main() {

	// list all available channels
	channels := factory.ListAllChannels()
	for i, x := range channels {
		fmt.Printf("Available channels for bustype '%v': %v", i, x)
	}

	// connect to new channel
	cfg := gocan.Config{BusType: "pcan", Channel: "PCAN_USBBUS1", BaudRate: 500000}
	canbus, err := factory.CreateBus(&cfg)
	if err != nil {
		fmt.Printf("Error creating can connection: %e", err)
		return
	}

	// send a extended message
	msg := gocan.Message{ID: 0x123, Data: []byte{0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55}, DLC: 8}
	err = canbus.Send(&msg)
	if err != nil {
		fmt.Printf("Error sending message: %e", err)
		return
	}

	// send a standard message
	msg = gocan.Message{ID: 0x12345, Data: []byte{0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55}, DLC: 8, IsExtended: true}
	err = canbus.Send(&msg)
	if err != nil {
		fmt.Printf("Error sending message: %e", err)
		return
	}

	// read a message with timeout (only prints some if another device is sending)
	rxmsg, err := canbus.Recv(500)
	if err != nil {
		fmt.Printf("Error reading message: %e", err)
		return
	}
	fmt.Printf("Received message with ID %v and data: %v", rxmsg.ID, rxmsg.Data)
}
