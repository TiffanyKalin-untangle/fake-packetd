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

func GetConntrackTable() ([]map[string]interface{}) {
	var conntrackTable []map[string]interface{}
	for i := 1; i <= 10; i++ {
		conntrack := make(map[string]interface{})
		conntrack["conntrack_id"] = uint32(i * 1000)
		conntrack["timestamp_start"]= uint64(1+i)
		time.Sleep(time.Millisecond)
		conntrack["timestamp_stop"] = uint64(2+i)

		conntrackTable = append(conntrackTable, conntrack)
	}

	return conntrackTable
}