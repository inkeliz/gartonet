package main

import (
	"time"

	"github.com/inkeliz/gartonet"
)

func main() {
	client, err := gartonet.NewClientString("192.168.0.60")
	if err != nil {
		panic(err)
	}

	packet := gartonet.NewPacket(0, 0)
	for range time.Tick(time.Second / 40) {

		// DMX commands:
		for i := 0; i < len(packet.DMX); i++ {
			packet.DMX[i] = byte(i)
		}

		// Sending:
		if err := client.Send(packet); err != nil {
			panic(err)
		}
	}
}
