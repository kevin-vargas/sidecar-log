package main

import (
	"fmt"
	"sidecar/k3s"
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
				client := k3s.New()
				logsbytes, err := client.GetLogs()
				if err != nil {
					fmt.Println("NUEVA ITERACION XDDDDDDDD")
					fmt.Println(err)
				} else {
					fmt.Println(string(logsbytes))
				}
			}()
		case <-quit:
			ticker.Stop()
			return
		}
	}

}
