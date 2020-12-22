package zmqd

import (
	"sync"
	"syscall"
	"time"

	zmq "github.com/pebbe/zmq4"
	"github.com/untangle/golang-shared/services/logger"
	zreq "github.com/untangle/golang-shared/structs/protocolbuffers/ZMQRequest"
	"google.golang.org/protobuf/proto"
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
				request, err := socket.RecvMessageBytes(zmq.DONTWAIT)
				if err != nil {
					if zmq.AsErrno(err) != zmq.AsErrno(syscall.EAGAIN) {
						logger.Warn("Error on receive ", err, "\n")
					}
					continue
				}
				logger.Info("Received ", request, "\n")

				// Process message
				reply, err := processMessage(request)
				if err != nil {
					logger.Warn("Error on processing ", err, "\n")
					continue
				}

				socket.SendMessage(reply)
				logger.Info("Sent ", reply, "\n")
			}
		} 
	}(&wg, zmqSocket)
}

func processMessage(msgRaw [][]byte) (processedReply []byte, processErr error) {
	reply := &zreq.ZMQRequest{}
	if err := proto.Unmarshal(msgRaw[0], reply); err != nil {
		return nil, err
	}

	time.Sleep(time.Second)

	encodedReply, err := proto.Marshal(reply)
	if err != nil {
		return nil, err
	}

	return encodedReply, nil
}