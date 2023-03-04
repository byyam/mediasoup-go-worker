package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/byyam/mediasoup-go-worker/pkg/udschannel"
)

func settings() []string {
	var options = make([]string, 0)
	options = append(options, fmt.Sprintf("--logLevel=%s", "wrn"))

	return options
}

func main() {
	opts := settings()
	controller, err := udschannel.NewController("../controlled/controlled", opts, 2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("controller start pid:%d\n", controller.Child.Pid)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer func() {
			defer wg.Done()
			ticker.Stop()
			fmt.Printf("stop routine task\n")
		}()

		for {
			select {
			case exitMsg := <-controller.ChildExit:
				fmt.Printf("controlled exit:%v\n", exitMsg)
				// todo: clear resources and exit
				return

			case <-ticker.C:
				fmt.Printf("ticker dida\n")
			}
		}
	}()
	wg.Wait()
	fmt.Printf("controller exit\n")
}
