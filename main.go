package main

import (
	"fmt"

	"github.com/strongswan/govici/vici"
)

var session *vici.Session

// init function can't be called,
// are automatically executed in the order in which they are declared.
func init() {
	var err error
	session, err = vici.NewSession()
	if err != nil {
		fmt.Println(err)
		return
	}
	//	defer session.Close()
}

func main() {

	listCACerts()
}
