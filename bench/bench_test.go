package main

import (
	"testing"

	"github.com/inkeliz/gartonet"
)

func BenchmarkGartonet(b *testing.B) {
	client, err := gartonet.NewClientString("192.168.0.60")
	if err != nil {
		panic(err)
	}

	b.ReportAllocs()

	packet := gartonet.NewPacket(0, 0)
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(packet.DMX); i++ {
			packet.DMX[i] = byte(i)
		}
		if err := client.Send(packet); err != nil {
			panic(err)
		}
	}
}
