package main

import (
	"context"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/strongswan/govici/vici"
)

func main() {
	session, err := vici.NewSession()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer session.Close()

	// print version info
	m, err := session.CommandRequest("version", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, k := range m.Keys() {
		fmt.Printf("%v: %v\n", k, m.Get(k))
	}

	//list CA certs
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
		fmt.Println(m1)
	}

	// Event listener

	if err := session.Subscribe("ike-updown", "log"); err != nil {
		fmt.Println(err)
		return
	}

	name := "rw"

	for {
		e, err := session.NextEvent(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		switch e.Name {
		case "ike-updown":
			state := e.Message.Get(name).(*vici.Message).Get("state")
			fmt.Println("IKE SA state changed (name=%s): %s\n", name, state)
		case "log":
			fmt.Println(e.Message.Get("msg"))
		}
	}
}

type cert struct {
	Type string `vici:"type"`
	Flag string `vici:"flag"`
	Data string `vici:"data"`
}

func loadX509Cert(path string, cacert bool) error {
	s, err := vici.NewSession()
	if err != nil {
		return err
	}
	defer s.Close()

	flag := "NONE"
	if cacert {
		flag = "CA"
	}

	//Read cert data from the file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(data)

	cert := cert{
		Type: "X509",
		Flag: flag,
		Data: string(block.Bytes),
	}

	m, err := vici.MarshalMessage(&cert)
	if err != nil {
		return err
	}

	_, err = s.CommandRequest("load-cert", m)

	return err
}

type connection struct {
	Name string // This field will NOT be marshaled!

	LocalAddrs []string            `vici:"local_addrs"`
	Local      *localOpts          `vici:"local"`
	Remote     *remoteOpts         `vici:"remote"`
	Children   map[string]*childSA `vici:"children"`
	Version    int                 `vici:"version"`
	Proposals  []string            `vici:"proposals"`
}

type localOpts struct {
	Auth  string   `vici:"auth"`
	Certs []string `vici:"certs"`
	ID    string   `vici:"id"`
}

type remoteOpts struct {
	Auth string `vici:"auth"`
}

type childSA struct {
	LocalTrafficSelectors []string `vici:"local_ts"`
	Updown                string   `vici:"updown"`
	ESPProposals          []string `vici:"esp_proposals"`
}

func loadConn(conn connection) error {
	s, err := vici.NewSession()
	if err != nil {
		return err
	}
	defer s.Close()

	c, err := vici.MarshalMessage(&conn)
	if err != nil {
		return err
	}

	m := vici.NewMessage()
	if err := m.Set(conn.Name, c); err != nil {
		return err
	}

	_, err = s.CommandRequest("load-conn", m)

	return err
}

func initiate(ike, child string) error {
	s, err := vici.NewSession()
	if err != nil {
		return err
	}
	defer s.Close()

	m := vici.NewMessage()
	if err := m.Set("child", child); err != nil {
		return err
	}

	if err := m.Set("ike", ike); err != nil {
		return err
	}

	ms, err := s.StreamedCommandRequest("initiate", "control-log", m)
	if err != nil {
		return err
	}

	for _, msg := range ms.Messages() {
		if err := msg.Err(); err != nil {
			return err
		}
	}

	return nil
}
