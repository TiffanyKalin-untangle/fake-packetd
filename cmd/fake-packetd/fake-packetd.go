package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/untangle/packetd/services/logger"
	"github.com/TiffanyKalin-untangle/fake-packetd/services/fake-restd"
)

var shutdownFlag bool

func main() {
	startServices()

	handleSignals()

	for !getShutdown() {
		select {
		case <-time.After(2 * time.Second):
			logger.Debug("fake-packetd is running...\n")
		}
	}

	stopServices()
	logger.Info("Shutdown fake-packetd...\n")
}

func startServices() {
	setIsShutdown(false)
	logger.Startup()
	fake_restd.Startup()
	logger.Info("fake-packetd starting up...\n")
}

func stopServices() {
	logger.Shutdown()
}

func setIsShutdown(flag bool) {
	shutdownFlag = flag
}

func getShutdown() bool {
	return shutdownFlag
}

func handleSignals() {
	// Add SIGINT & SIGTERM handler (exit)
	termch := make(chan os.Signal, 1)
	signal.Notify(termch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-termch
		go func() {
			logger.Info("Received signal [%v]. Setting shutdown flag\n", sig)
			setIsShutdown(true)
		}()
	}()
}