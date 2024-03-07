package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

// Реализовать простейший telnet-клиент.

// Примеры вызовов:
// go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

// Требования:
// Программа должна подключаться к указанному хосту (ip или доменное имя + порт) по протоколу TCP.
// После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
// Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s)
// При нажатии Ctrl+D программа должна закрывать сокет и завершаться.
// Если сокет закрывается со стороны сервера, программа должна также завершаться.
// При подключении к несуществующему сервер, программа должна завершаться через timeout

type telnetParameters struct {
	Timeout time.Duration
}

func main() {
	parameters := parceCmdArgs()

	if flag.NArg() != 2 {
		fmt.Println("Use: go-telnet [--timeout=10s] host port")
		os.Exit(1)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)

	address := fmt.Sprintf("%s:%s", host, port)

	connect, err := net.DialTimeout("tcp", address, parameters.Timeout)
	if err != nil {
		fmt.Printf("Connection error: %s: %v\n", address, err)
		os.Exit(1)
	}
	defer connect.Close()

	fmt.Println("Connected to", address)

	go readFromServer(connect)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintf(connect, "%s\n", line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error when reading data: %v\n", err)
		os.Exit(1)
	}
}

func readFromServer(connect net.Conn) {
	scanner := bufio.NewScanner(connect)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
	fmt.Println("The connection was interrupted")
	os.Exit(0)
}

func parceCmdArgs() telnetParameters {
	parameters := telnetParameters{}

	timeout := flag.Duration("timeout", 10*time.Second, "Connection timeout")
	flag.Parse()

	parameters.Timeout = *timeout

	return parameters
}
