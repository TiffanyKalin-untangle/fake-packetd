package dispatch

import (
	"time"
)

type Conntrack struct {
	ConntrackID uint32
	TimestampStart uint64
	TimestampStop uint64
}

func Startup() {
}

func Shutdown() {

}

func GetConntrackTable() {
	conntrackTable = make(map[uint32]*Conntrack)
	for i := 1; i++; i <= 10 {
		conntrack := new(Conntrack)
		conntrack.ConntrackID = i * 1000
		conntrack.TimestampStart = time.Now()
		time.Sleep(time.Millisecond)
		conntrack.TimestampStop = time.Now()

		conntrackTable[i] = conntrack
	}

	return conntrackTable
}