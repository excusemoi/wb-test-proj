package main

import (
	"bufio"
	"net"
)

func serverLogic() {
	listener, _ := net.Listen("tcp", ":5000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			println("Cant connect")
			conn.Close()
			continue
		}
		println("Connected")

		bufReader := bufio.NewReader(conn)
		println("Start reading")

		go func(conn net.Conn) {
			defer conn.Close()
			for {
				rb, err_ := bufReader.ReadByte()
				if err_ != nil {
					println("Cant read byte")
					break
				}
				print(string(rb))
				conn.Write([]byte("\nReceived"))
			}
		}(conn)

	}

}
