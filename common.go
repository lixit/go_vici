package main

import (
	"fmt"

	"github.com/strongswan/govici/vici"
)

// print version info
func version() {
	m, err := session.CommandRequest("version", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, k := range m.Keys() {
		fmt.Printf("%v: %v\n", k, m.Get(k))
	}
}

// Get Statistics info
func stats() {
	m, err := session.CommandRequest("stats", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// for _, k := range m.Keys() {
	// 	fmt.Printf("%v: %v\n", k, m.Get(k))
	// }

	fmt.Printf("%v: %v\n", "running", m.Get("uptime").(*vici.Message).Get("running"))
	fmt.Printf("%v: %v\n", "since", m.Get("uptime").(*vici.Message).Get("since"))
	fmt.Printf("%v: %v\n", "plugins", m.Get("plugins"))

}

func reloadSettings() {
	//reload-settings
	m, err := session.CommandRequest("reload-settings", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	if m.Get("success") == "yes" {
		fmt.Printf("%v: %v\n", "success", "yes")
	} else {
		fmt.Printf("%v: %v\n", "err", m.Get("errmsg"))
	}

}
