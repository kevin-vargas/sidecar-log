package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/kevin-vargas/sidecar-log/k3s"
	"github.com/kevin-vargas/sidecar-log/pubsub"

	"github.com/joho/godotenv"
)

const POD_ID = "HOSTNAME"

type logPubSub struct {
	pubsub.MQTTI
	topic string
}

//TODO: concurrent log publish, io reader io write not a lot of copys.

func makeEntryHandle(m *logPubSub) func(entry []byte) {
	return func(entry []byte) {
		//TODO: error handle
		m.Publish(m.topic, string(entry))
	}
}

//TODO: error handling
func readLog(m *logPubSub, log []byte) {
	entryHandle := makeEntryHandle(m)
	buf := new(bytes.Buffer)
	for _, elem := range log {
		if elem == 10 {
			entryHandle(buf.Bytes())
			buf.Reset()
		} else {
			buf.WriteByte(elem)
		}
	}
	// flush
	if buf.Len() > 0 {
		entryHandle(buf.Bytes())
	}
}
func main() {
	fmt.Println("Running Sidecar logger")
	godotenv.Load(".env")
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	clientK3S := k3s.New()
	clientMQTT := pubsub.New()
	clientLogger := &logPubSub{
		clientMQTT,
		"log",
	}
	for {
		select {
		case <-ticker.C:
			go func() {
				logsbytes, err := clientK3S.GetLogs()
				if err != nil {
					fmt.Println(err)
				} else {
					go readLog(clientLogger, logsbytes)
				}
			}()
		case <-quit:
			ticker.Stop()
			return
		}
	}

}
