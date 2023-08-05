package factory

import (
	"errors"

	"github.com/morgadow/gocan"
	"github.com/morgadow/gocan/interfaces/pcan"
)

// Creates and initializes a connection to a CANBus
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

// Lists all available channels for all manufactures
func ListChannels() map[string][]string {

	var channels = make(map[string][]string)

	// pcan
	perr := pcan.LoadAPI()
	if perr == nil {
		pcanHandles, perr := pcan.AttachedChannelsNames()
		if perr == nil {
			channels["pcan"] = pcanHandles
		}
	}

	return channels
}
