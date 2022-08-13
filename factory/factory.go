package factory

import (
	"errors"

	"github.com/morgadow/gocan"
	"github.com/morgadow/gocan/interfaces/pcan"
)

// CreateBus Creates and initializes a connection to a CANBus
func CreateBus(config *gocan.Config) (gocan.Bus, error) {

	var newBus gocan.Bus = nil
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
