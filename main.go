package main

import (
	"fmt"
	"os"
)

const POD_ID = "HOSTNAME"

func main() {
	podId := os.Getenv(POD_ID)
	if len(podId) > 0 {
		fmt.Println("POD ID")
		fmt.Println(podId)
	} else {
		fmt.Println("NO HOSTNAME ENV")
	}
}
