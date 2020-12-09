package fake_restd

import (
	"time"

	zmq "github.com/pebbe/zmq4"
	"github.com/untangle/packetd/services/logger"
)

func Startup() {
	logger.Info("Starting fake-restd...\n")

	socket, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		logger.Warn("Failed to create zmq socket...", err)
	}
	defer socket.Close()

	

	go func {
		for {
			msg, _ := socket.Recv(0)
			logger.Info("Received ", string(msg))

			time.Sleep(time.Second)

			reply := "World"
			socket.Send(reply, 0)
			logger.Info("Sent ", reply)
		}
	}
}

func setupZmqSocket() (soc *zmq.Socket, err error) {
	socket.Bind("tcp://*:5555")
}