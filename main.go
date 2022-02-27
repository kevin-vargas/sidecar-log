package main

import (
	"fmt"
	"os"
	"time"
)

const POD_ID = "HOSTNAME"

func main() {
	fmt.Println("Running Sidecar logger")
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			go func() {
				podId := os.Getenv(POD_ID)
				if len(podId) > 0 {
					fmt.Println("POD ID")
					fmt.Println(podId)
				} else {
					fmt.Println("NO HOSTNAME ENV")
				}
			}()
		case <-quit:
			ticker.Stop()
			return
		}
	}

}
