package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

func writeToWs(ws *websocket.Conn, writeCode int, data []byte) (err error) {
	ws.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.WriteMessage(writeCode, data)
}

//from io.Reader to websocket's message
func pumpStdout(ws *websocket.Conn, r io.Reader, done chan struct{}) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		fmt.Println(s.Text())
		if err := writeToWs(ws, websocket.TextMessage, s.Bytes()); err != nil {
			ws.Close()
			break
		}
	}
	if s.Err() != nil {
		fmt.Println("scan:", s.Err())
	}
	close(done)

	writeToWs(ws, websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	ws.Close()
}

func internalError(ws *websocket.Conn, msg string, err error) {
	fmt.Println(msg, err)
	ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade:", err)
		return
	}

	defer ws.Close()

	cmd := exec.Command(ExecFname, CmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating StdoutPipe for Cmd", err)
	}

	stdoutDone := make(chan struct{})
	go pumpStdout(ws, cmdReader, stdoutDone)

	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting Cmd", err)
	}

	select { //If no stdout, waits for 5 seconds and terminates proc
	case <-stdoutDone:
		fmt.Println("Got message from stdout")
		//	case <-time.After(time.Second * 5):
		//		fmt.Println("Killing proc")
		//		// A bigger bonk on the head.
		//		if err := proc.Signal(os.Kill); err != nil {
		//			fmt.Println("term:", err)
		//		}
		<-stdoutDone
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error waiting for Cmd", err)
	}
}
