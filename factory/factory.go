package factory

import (
	"errors"

	"github.com/morgadow/gocan/bus"
	"github.com/morgadow/gocan/interfaces/pcan"
)

// CreateBus Creates and initializes a connection to a CANBus
func CreateBus(config *bus.Config) (bus.Bus, error) {

	var newBus bus.Bus = nil
	var err error = nil

	// create the selected can bus connection
	switch config.BusType {
	case "pcan":
		newBus, err = pcan.NewPCANBus(config)

	default:
		return nil, errors.New("invalid interface selected or interface not implemented")
	}

	return newBus, err
}
