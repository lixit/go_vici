package main

import "github.com/strongswan/govici/vici"

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

	c, err := vici.MarshalMessage(&conn)
	if err != nil {
		return err
	}

	m := vici.NewMessage()
	if err := m.Set(conn.Name, c); err != nil {
		return err
	}

	_, err = session.CommandRequest("load-conn", m)

	return err
}

func initiate(ike, child string) error {

	m := vici.NewMessage()

	if err := m.Set("child", child); err != nil {
		return err
	}

	if err := m.Set("ike", ike); err != nil {
		return err
	}

	ms, err := session.StreamedCommandRequest("initiate", "control-log", m)
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
