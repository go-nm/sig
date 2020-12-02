package sig

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

var stopSignals = []os.Signal{
	syscall.SIGHUP,
	syscall.SIGINT,
	syscall.SIGQUIT,
	syscall.SIGTERM,
}

// StopSignalSync func
func StopSignalSync(stopFn func()) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, stopSignals...)

	sig := <-stopChan
	log.Printf("Received OS signal %s\n", sig)
	stopFn()
}

// StopSignalE func
func StopSignalE(startFn func() error, stopFn func() error) error {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, stopSignals...)

	errChan := make(chan error, 1)
	go func() {
		errChan <- startFn()
	}()

	select {
	case sig := <-stopChan:
		log.Printf("Received OS signal %s\n", sig)
		return stopFn()

	case err := <-errChan:
		return err
	}
}
