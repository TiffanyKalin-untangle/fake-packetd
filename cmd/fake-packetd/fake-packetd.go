package main

import (
	"github.com/untangle/packetd/services/logger"
)

func main() {
	logger.Startup()
	logger.Info("fake-packetd staring up...\n")
}