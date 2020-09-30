package main

import (
	"context"
	"fmt"

	"github.com/strongswan/govici/vici"
)

func events() {
	if err := session.Subscribe("ike-updown", "log"); err != nil {
		fmt.Println(err)
		return
	}

	name := "rw1"

	for {
		e, err := session.NextEvent(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		switch e.Name {
		case "ike-updown":
			state := e.Message.Get(name).(*vici.Message).Get("state")
			fmt.Printf("IKE SA state changed (name=%s): %s\n", name, state)
		case "log":
			fmt.Println(e.Message.Get("msg"))
		}
	}
}
