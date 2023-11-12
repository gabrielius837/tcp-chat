package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/google/uuid"
)

type serverCtx struct {
	mutex sync.RWMutex
	users map[uuid.UUID]net.Conn
}

func newServerCtx() *serverCtx {
	ctx := &serverCtx{
		users: make(map[uuid.UUID]net.Conn),
	}
	return ctx
}

const (
	PROTOCOL = "tcp"
	ADDRESS  = "127.0.0.1:5555"
	PASSWORD = "minecraft\n"
)

func writeMsg(conn net.Conn, msg string) error {
	bytes := []byte(msg)
	_, err := conn.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func readMsg(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return msg, nil
}

func fatal(err error) {
	if err == nil {
		return
	}

	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func (ctx *serverCtx) userConnect(uuid uuid.UUID, conn net.Conn) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	ctx.users[uuid] = conn
	for _, userConn := range ctx.users {
		writeMsg(userConn, fmt.Sprintf("%s has connected\n", uuid.String()))
	}
}

func (ctx *serverCtx) userDisconnect(uuid uuid.UUID) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	delete(ctx.users, uuid)
	for _, userConn := range ctx.users {
		writeMsg(userConn, fmt.Sprintf("%s has disconnected\n", uuid.String()))
	}
}

func (ctx *serverCtx) broadcastMessage(uuid uuid.UUID, msg string) {
	ctx.mutex.RLock()
	defer ctx.mutex.RUnlock()
	for _, userConn := range ctx.users {
		writeMsg(userConn, fmt.Sprintf("%s: %s", uuid.String(), msg))
	}
}

func handleRequest(ctx *serverCtx, uuid uuid.UUID, conn net.Conn) {
	defer conn.Close()

	writeMsg(conn, "enter password:\n")
	password, err := readMsg(conn)
	if password != PASSWORD || err != nil {
		return
	}

	writeMsg(conn, fmt.Sprintf("welcome %s\n", uuid.String()))
	defer ctx.userDisconnect(uuid)
	for {
		msg, err := readMsg(conn)
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Fprintln(os.Stderr, err)
			}
			break
		}
		ctx.broadcastMessage(uuid, msg)
	}
}

func main() {
	cert, err := tls.LoadX509KeyPair("certificates/certificate.crt", "certificates/private.key")
	fatal(err)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := tls.Listen(PROTOCOL, ADDRESS, tlsConfig)
	fatal(err)
	defer listener.Close()

	ctx := newServerCtx()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		uuid, _ := uuid.NewRandom()
		ctx.userConnect(uuid, conn)

		go handleRequest(ctx, uuid, conn)
	}
}
