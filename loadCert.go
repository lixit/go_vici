package main

import (
	"encoding/pem"
	"io/ioutil"

	"github.com/strongswan/govici/vici"
)

type cert struct {
	Type string `vici:"type"`
	Flag string `vici:"flag"`
	Data string `vici:"data"`
}

func loadX509Cert(path string, cacert bool) error {

	flag := "NONE"
	if cacert {
		flag = "CA"
	}

	// Read cert data from the file
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

	_, err = session.CommandRequest("load-cert", m)

	return err
}
