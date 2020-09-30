package main

import (
	"fmt"

	"github.com/strongswan/govici/vici"
)

//list CA certs

func listCACerts() {
	m1 := vici.NewMessage()

	if err := m1.Set("flag", "CA"); err != nil {
		fmt.Println(err)
		return
	}

	ms, err := session.StreamedCommandRequest("list-certs", "list-cert", m1)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, m1 := range ms.Messages() {
		if m1.Err() != nil {
			fmt.Println(err)
			return
		}

		// Process CA cert information
		fmt.Printf("type: %v\n", m1.Get("type"))
		fmt.Printf("flag: %v\n", m1.Get("flag"))
	}
}
