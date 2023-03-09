package main

import (
	"fmt"
	"github.com/hypebeast/go-osc/osc"
)

func main() {
	addr := "127.0.0.1:8765"
	d := osc.NewStandardDispatcher()
	err := d.AddMsgHandler("/message/address", func(msg *osc.Message) {
		osc.PrintMessage(msg)
	})
	if err != nil {
		fmt.Println("d.AddMShHandler()")
	}

	server := &osc.Server{
		Addr:       addr,
		Dispatcher: d,
	}

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}

	client := osc.NewClient("localhost", 8765)
	msg := osc.NewMessage("/osc/address")
	msg.Append(int32(111))
	msg.Append(true)
	msg.Append("hello")
	err = client.Send(msg)
	if err != nil {
		fmt.Println("client.Sned")
	}

}
