package zmqd

import (
	"sync"
	"syscall"
	"time"

	zmq "github.com/pebbe/zmq4"
	"github.com/untangle/packetd/services/logger"
)

var isShutdown = make(chan struct{})
var wg sync.WaitGroup

func Startup() {
	logger.Info("Starting zmq service...\n")
	socketServer()

}

func Shutdown() {
	close(isShutdown)
	wg.Wait()
}

func socketServer() {
	zmqSocket, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		logger.Warn("Failed to create zmq socket...", err)
	}

	zmqSocket.Bind("tcp://*:5555")
	wg.Add(1)
	go func(waitgroup *sync.WaitGroup, socket *zmq.Socket) {
		defer socket.Close()
		defer waitgroup.Done()
		tick := time.Tick(500 * time.Millisecond)
		for {
			select {
			case <-isShutdown:
				logger.Info("Shutdown is seen\n")
				return
			case <-tick:
				logger.Info("Listening for requests\n")
				msg, err := socket.Recv(zmq.DONTWAIT)
				if err != nil {
					if zmq.AsErrno(err) != zmq.AsErrno(syscall.EAGAIN) {
						logger.Warn("Error on receive ", err, "\n")
					}
					continue
				}
				logger.Info("Received ", string(msg), "\n")

				// Process message
				time.Sleep(time.Second)

				socket.Send(msg, 0)
				logger.Info("Sent ", msg, "\n")
			}
		} 
	}(&wg, zmqSocket)
}