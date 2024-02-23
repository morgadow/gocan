package test

import (
	"fmt"
	"math"
	"testing"
	"time"

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

// test if message is received in expected interval
// Note: For this to work, send a message from another device with fixed interval of 100ms
func TestMsgInterval(t *testing.T) {
	pbus, err := auxInitBus("PCAN_USBBUS1")
	if err != nil {
		t.Errorf("error while creating bus: %v", err)
	}

	expID := gocan.MessageID(0x13000350)
	var expInterval int64 = 100 * 1000 // µs
	var measInterval int64 = 0
	startTime := time.Now()
	duration := 10 * time.Second
	timeoutRecv := 100 * time.Millisecond
	timeLastMsg := time.Time{}
	timeFirstSet := false
	measIntervalSet := false
	amountMsgs := 0
	expMsgs := 100

	// loop and receive multiple messages to check interval
	for time.Since(startTime) < duration {
		msg, _ := pbus.Recv(int(timeoutRecv))
		if msg != nil && msg.ID == expID {
			amountMsgs++
			if !timeFirstSet {
				timeLastMsg = time.Now()
				timeFirstSet = true
				continue
			}
			timeSinceLastMsg := time.Since(timeLastMsg).Microseconds()
			if !measIntervalSet {
				measInterval = timeSinceLastMsg
				measIntervalSet = true
				continue
			}
			measInterval = (measInterval*80 + timeSinceLastMsg*20) / 100
			timeLastMsg = time.Now()
		}
	}

	if !near32(amountMsgs, expMsgs, 2) {
		t.Errorf("invalid msg count. Got %v msgs\n", amountMsgs)
	} else {
		fmt.Printf("msg count okay. Got %v msgs\n", amountMsgs)
	}

	if !near64(measInterval, expInterval, 50) {
		t.Errorf("invalid interval. got: %v µs, expected: %v µs. \n", measInterval, expInterval)
	} else {
		fmt.Printf("interval okay. got: %v µs, expected: %v µs.\n", measInterval, expInterval)
	}
}

func near32(value, target, tolerance int) bool {
	return math.Abs(float64(value-target)) <= float64(tolerance)
}
func near64(value, target, tolerance int64) bool {
	return math.Abs(float64(value-target)) <= float64(tolerance)
}
