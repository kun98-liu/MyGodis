package tcp

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestListenAndServe(t *testing.T) {
	closeChan := make(chan struct{})

	listener, err := net.Listen("tcp", ":0")

	if err != nil {
		t.Error(err)
		return
	}
	//start server
	addr := listener.Addr().String()
	go ListenAndServe(listener, MakeEchoHandler(), closeChan)

	//client start
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("A new conn set up")

	for i := 0; i < 10; i++ {
		val := strconv.Itoa(rand.Int())
		_, err := conn.Write([]byte(val + "\n"))

		if err != nil {
			t.Error(err)
			return
		}

		bufReader := bufio.NewReader(conn)

		line, _, err := bufReader.ReadLine()
		if err != nil {
			t.Error(err)
			return
		}

		if string(line) != val {
			t.Error("EchoServer didn't echo the right msg")
			return
		} else {
			fmt.Println("val = " + val + " | echo : " + string(line))
		}
	}

	conn.Close()
	fmt.Println("Client close the connection")
	time.Sleep(2 * time.Second)

	closeChan <- struct{}{}

	time.Sleep(2 * time.Second)
}
