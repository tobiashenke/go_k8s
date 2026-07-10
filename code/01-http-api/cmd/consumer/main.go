package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
)

func main() {
	const natsUrl = "nats://localhost:4222"
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		log.Fatal(err)
	}
	_, err = nc.Subscribe("items.created", func(msg *nats.Msg) { log.Printf("received: %s", msg.Data) })
	if err != nil {
		log.Fatal(err)
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
