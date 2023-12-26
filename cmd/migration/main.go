package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	diContainer := BuildContainer()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("\nOS Signal Received. Running Teardown..")
		diContainer.Invoke(Teardown)
		os.Exit(1)
	}()

	defer diContainer.Invoke(Teardown)

	err := diContainer.Invoke(Init)
	if err != nil {
		log.Panic("Invoke container Init error:", err)
	}
}
