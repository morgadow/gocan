package test

import (
	"fmt"
	"testing"

	"github.com/morgadow/gocan"
	"github.com/morgadow/gocan/interfaces/pcan"
)

func auxInitBus(channel string) (gocan.Bus, error) {
	cfg := gocan.Config{BusType: "pcan", Channel: channel, BaudRate: 500000, BusState: gocan.ACTIVE, IsFD: false, RecvStatusFrames: true, LogStatusFrames: true}
	pbus, err := pcan.NewPCANBus(&cfg)
	return pbus, err
}

func TestChannelCondition(t *testing.T) {
	pbus, err := auxInitBus("PCAN_USBBUS1")
	if err != nil {
		t.Errorf("error while creating bus: %v", err)
	}

	cond, err := pbus.ChannelCondition()
	fmt.Println(cond)

}
