package tracker

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"regexp"
	"time"
)

const (
	pid             uint64 = 0x41727101980
	timeoutInterval int    = 3
	retryCount      int    = 0
)

func IsUdp(link string) bool {
	match, _ := regexp.MatchString("^udp", link)
	return match
}

func ProcessUdp(link string, responseChan chan []byte) {
	rand.Seed(time.Now().UnixNano())

	handleTracker(link, responseChan)
}

func handleTracker(url string, responseChan chan []byte) {
	conn := getConnection(url)

	tries := 0
	timeoutSeconds := timeoutInterval
	recv := make([]byte, 16)

	for {
		go sendConnectionRequest(conn)
		deadline := time.Now().Add(time.Duration(time.Duration(timeoutSeconds) * time.Second))
		conn.SetReadDeadline(deadline)
		_, err := conn.Read(recv)

		if err != nil {
			fmt.Println("timeout")
			if tries < retryCount {
				tries++
				timeoutSeconds += timeoutInterval
				continue
			} else {
				responseChan <- nil
			}
		} else {
			responseChan <- recv
		}

		break
	}
}

func getConnection(addr string) *net.UDPConn {
	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func sendConnectionRequest(conn *net.UDPConn) {
	time.Sleep(1 * time.Second)
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[0:], pid)
	binary.BigEndian.PutUint32(buf[8:], 0)
	binary.BigEndian.PutUint32(buf[12:], rand.Uint32())

	_, err := conn.Write(buf)

	if err != nil {
		log.Fatal(err)
	}
}
