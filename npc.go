package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

type NPC struct {
	port string
	rx   chan []byte
	tx   chan byte
}

func (n *NPC) writeByte() {
	for b := range n.rx {
		cpu.dmem[0xc0] = b[0]
	}
}

func (n *NPC) readByte() {
	for {
		n.tx <- cpu.dmem[0xc6]
	}
}

func receive_data(c net.Conn, out chan []byte) {
	defer close(out)
	fmt.Fprintf(c, "Greetings, Professor Falken\n")
	reader := bufio.NewReader(c)
	for {
		fmt.Fprintf(c, "$> ")
		b, err := reader.ReadBytes(byte(10))
		if err != nil {
			fmt.Fprintf(c, "Sorry, Charlie! %v\n", err)
			break
		}
		if string(b) == "quit\r\n" {
			c.Close()
			break
		} else {
			out <- b
		}
	}
}

func send_data(c net.Conn, in <-chan byte) {
	defer c.Close()
	for {
		m := <-in
		io.Copy(c, bytes.NewBufferString(string(m)))
	}
}

func (n *NPC) Server() {
	listener, err := net.Listen("tcp", n.port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			break
		}
		go receive_data(conn, n.rx)
		go n.writeByte()
		go n.readByte()
		go send_data(conn, n.tx)
	}
}
