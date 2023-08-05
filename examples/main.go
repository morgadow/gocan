package main

import (
	"fmt"

	"github.com/morgadow/gocan/factory"
)

func main() {

	// list all available channels
	channels := factory.ListAllChannels()
	fmt.Println(channels)
}
