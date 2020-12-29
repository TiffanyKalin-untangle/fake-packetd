package dispatch

import (
	"net"
	"sync"
	"time"
)

type Test struct {
	ConntrackID       uint32
	ConnMark          uint32
	SessionID         int64
	Family            uint8
	CreationTime      time.Time
	LastUpdateTime    time.Time
	LastActivityTime  time.Time
	ClientSideTuple   Tuple
	ServerSideTuple   Tuple
	TimeoutSeconds    uint32
	TimestampStart    uint64
	TimestampStop     uint64
	TCPState          uint
	EventCount        uint64
	ClientBytes       uint64
	ServerBytes       uint64
	TotalBytes        uint64
	ClientPackets     uint64
	ServerPackets     uint64
	TotalPackets      uint64
	ClientBytesDiff   uint64  // the ClientBytes diff since last update
	ServerBytesDiff   uint64  // the ServerBytes diff since last update
	TotalBytesDiff    uint64  // the TotalBytes diff since last update
	ClientPacketsDiff uint64  // the ClientPackets diff since last update
	ServerPacketsDiff uint64  // the ServerPackets diff since last update
	TotalPacketsDiff  uint64  // the TotalPackets diff since last update
	ClientByteRate    float32 // the Client byte rate site the last update
	ServerByteRate    float32 // the Server byte rate site the last update
	TotalByteRate     float32 // the Total byte rate site the last update
	ClientPacketRate  float32 // the Client packet rate site the last update
	ServerPacketRate  float32 // the Server packet rate site the last update
	TotalPacketRate   float32 // the Total packet rate site the last update
	Guardian          sync.RWMutex
}

func Startup() {
}

func Shutdown() {

}

func GetConntrackTable() ([]map[string]interface{}) {
	var conntrackTable []map[string]interface{}
	for i := 1; i <= 10; i++ {
		conntrack := new(Test)
		conntrack.ConntrackID = 10
		conntrack.ConnMark = 10
		conntrack.CreationTime = time.Now()
		conntrack.Family = 2
		conntrack.LastActivityTime = time.Now()
		conntrack.EventCount = 1
		conntrack.ClientSideTuple.Protocol = 2
		conntrack.ClientSideTuple.ClientAddress = make(net.IP, 4)
		conntrack.ClientSideTuple.ClientPort = 10
		conntrack.ClientSideTuple.ServerAddress = make(net.IP, 4)
		conntrack.ClientSideTuple.ServerPort = 10
		conntrack.ServerSideTuple.Protocol = 10
		conntrack.ServerSideTuple.ClientAddress = make(net.IP, 4)
		conntrack.ServerSideTuple.ClientPort = 10
		conntrack.ServerSideTuple.ServerAddress = make(net.IP, 4)
		conntrack.ServerSideTuple.ServerPort = 10
		conntrack.ClientBytes = 10
		conntrack.ServerBytes = 10
		conntrack.TotalBytes = 10
		conntrack.ClientPackets = 10
		conntrack.ServerPackets = 10
		conntrack.TotalPackets = 10
		conntrack.TimeoutSeconds = 10
		conntrack.TimestampStart = 10
		conntrack.TimestampStop = 10
		conntrack.TCPState = 10
		conntrackMap := parseConntrack(conntrack)

		conntrackTable = append(conntrackTable, conntrackMap)
	}

	return conntrackTable
}

// parse a line of /proc/net/nf_conntrack and return the info in a map
// parse a line of /proc/net/nf_conntrack and return the info in a map
func parseConntrack(ct *Test) map[string]interface{} {
	m := make(map[string]interface{})

	// ignore 127.0.0.1 traffic
	if ct.ClientSideTuple.ClientAddress != nil && ct.ClientSideTuple.ClientAddress.String() == "127.0.0.1" {
		return nil
	}
	if ct.ClientSideTuple.ServerAddress != nil && ct.ClientSideTuple.ServerAddress.String() == "127.0.0.1" {
		return nil
	}
	if ct.ClientSideTuple.ClientAddress != nil && ct.ClientSideTuple.ClientAddress.String() == "::1" {
		return nil
	}
	if ct.ClientSideTuple.ServerAddress != nil && ct.ClientSideTuple.ServerAddress.String() == "::1" {
		return nil
	}

	m["conntrack_id"] = uint(ct.ConntrackID)
	m["session_id"] = uint(ct.SessionID)
	m["family"] = uint(ct.Family)
	m["ip_protocol"] = uint(ct.ClientSideTuple.Protocol)

	m["timeout_seconds"] = uint(ct.TimeoutSeconds)
	m["tcp_state"] = uint(ct.TCPState)

	m["client_address"] = ct.ClientSideTuple.ClientAddress.String()
	m["client_port"] = uint(ct.ClientSideTuple.ClientPort)
	m["server_address"] = ct.ClientSideTuple.ServerAddress.String()
	m["server_port"] = uint(ct.ClientSideTuple.ServerPort)
	m["client_address_new"] = ct.ServerSideTuple.ClientAddress.String()
	m["client_port_new"] = uint(ct.ServerSideTuple.ClientPort)
	m["server_address_new"] = ct.ServerSideTuple.ServerAddress.String()
	m["server_port_new"] = uint(ct.ServerSideTuple.ServerPort)

	m["bytes"] = uint(ct.TotalBytes)
	m["client_bytes"] = uint(ct.ClientBytes)
	m["server_bytes"] = uint(ct.ServerBytes)
	m["packets"] = uint(ct.TotalPackets)
	m["client_packets"] = uint(ct.ClientPackets)
	m["server_packets"] = uint(ct.ServerPackets)

	m["timestamp_start"] = uint(ct.TimestampStart)
	if ct.TimestampStart != 0 {
		m["age_milliseconds"] = uint((uint64(time.Now().UnixNano()) - ct.TimestampStart) / 1000000)
	}

	var mark uint32
	mark = ct.ConnMark
	clientInterfaceID := mark & 0x000000ff
	clientInterfaceType := mark & 0x03000000 >> 24
	serverInterfaceID := mark & 0x0000ff00 >> 8
	serverInterfaceType := mark & 0x0c000000 >> 26
	priority := mark & 0x00ff0000 >> 16
	m["mark"] = uint(mark)
	m["client_interface_id"] = uint(clientInterfaceID)
	m["client_interface_type"] = uint(clientInterfaceType)
	m["server_interface_id"] = uint(serverInterfaceID)
	m["server_interface_type"] = uint(serverInterfaceType)
	m["priority"] = priority

	return m
}