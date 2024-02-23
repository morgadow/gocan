package test

import (
	"fmt"
	"testing"

	"github.com/morgadow/gocan"
	"github.com/morgadow/gocan/interfaces/pcan"
)

func auxInitBus(channel string) (gocan.Bus, error) {
	cfg := gocan.Config{BusType: "pcan", Channel: channel, BaudRate: 500000, BusState: gocan.ACTIVE}
	pbus, err := pcan.NewPCANBus(&cfg)
	return pbus, err
}

func auxUnitBus(pbus gocan.Bus) error {
	err := pbus.Shutdown()
	return err
}

func TestChannelCondition(t *testing.T) {
	pbus, err := auxInitBus("PCAN_USBBUS1")
	if err != nil {
		t.Errorf("error while creating bus: %v", err)
	}

	cond, err := pbus.ChannelCondition()
	if err != nil {
		t.Errorf("error while reading channel condition: %v", err)
	}
	if cond != gocan.Occupied {
		t.Errorf("got wrong channel condition: %v", cond)
	}
	fmt.Println(cond)
}

func TestRecv(t *testing.T) {
	pbus, err := auxInitBus("PCAN_USBBUS1")
	if err != nil {
		t.Errorf("error while creating bus: %v", err)
	}

	msg, err := pbus.Recv(5000)
	if msg == nil || msg.ID == 0 || err != nil {
		t.Errorf("no message: msg: %v, err: %v", msg, err)
		return
	} else {
		fmt.Printf("received message: msg: %v, err: %v\n", msg, err)
	}

	if msg.Type != gocan.DataFrame || msg.IsExtended != true || len(msg.Data) == 0 {
		t.Errorf("invalid message: msg type: %v, msg id type: %v, msg data len %v", msg.Type, msg.IsExtended, msg.Data)
	}
}

func TestTraceStart(t *testing.T) {
	pbus, err := auxInitBus("PCAN_USBBUS1")
	if err != nil {
		t.Errorf("error while creating bus: %v", err)
	}

	tracePath := "C:/workspace/go/src/github.com/morgadow"
	fileSize := 100
	err = pbus.TraceStart(tracePath, uint32(fileSize))
	if err != nil {
		t.Errorf("error while starting trace: %v", err)
	}

	// activly read from queue to fill file
	pbus.Recv(1)
	pbus.Recv(1)
	pbus.Recv(1)
	pbus.Recv(1)
	pbus.Recv(1)

	err = pbus.TraceStop()
	if err != nil {
		t.Errorf("error while stopping trace: %v", err)
	}

	// TODO check there is a file with messages
}

func TestTraceStop(t *testing.T) {
	pbus, err := auxInitBus("PCAN_USBBUS1")
	if err != nil {
		t.Errorf("error while creating bus: %v", err)
	}

	err = pbus.TraceStop()
	if err != nil {
		t.Errorf("error while stopping trace: %v", err)
	}
}
