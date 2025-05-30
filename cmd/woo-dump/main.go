package main

import (
	"flag"
	"fmt"
	"github.com/rigidsh/woo-client/pkg/protocol"
	"io"
	"log"
	"os"
	"time"
	"tinygo.org/x/bluetooth"
)

func main() {
	addressFlag := flag.String("address", "", "Woo device address")
	reconnectFlag := flag.Bool("reconnect", true, "auto reconnect")
	dumpPathFlag := flag.String("dump", "", "dump file path")

	flag.Parse()

	adapter := bluetooth.DefaultAdapter
	err := adapter.Enable()
	if err != nil {
		log.Fatal("Cannot enable bluetooth:", err)
		return
	}

	var address bluetooth.Address

	if *addressFlag == "" {
		address, err = FindWooDevice(5*time.Second, adapter)
		if err != nil {
			log.Fatal("Cannot find Woo device: ", err)
			return
		}
	} else {
		address.Set(*addressFlag)
	}

	var dumpWriter io.Writer

	if *dumpPathFlag != "" {
		log.Printf("Dumping to %s", *dumpPathFlag)
		dumpWriter, err = os.OpenFile(*dumpPathFlag, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal("Cannot open dump file: ", err)
			return
		}
	}

	decoder := protocol.NewPackageDecoder(protocol.NewBufferPackageWriter(func(checksum bool, data []byte) {
		log.Printf("New package(checksum %t): % x", checksum, data)
		if checksum && dumpWriter != nil {
			_, _ = dumpWriter.Write([]byte(fmt.Sprintf("% x\n", data)))
		}
	}))

	err = ConnectToDevice(adapter, *reconnectFlag, address, func(device bluetooth.Device) {
		err = SubscribeWooEvents(device, func(data []byte) {
			_, err := decoder.Write(data)
			if err != nil {
				log.Fatal("Error on decoding message: ", err)
			}
		})
		if err != nil {
			log.Fatal("Cannot subscribe to Woo events: ", err)
			return
		}
	})
	if err != nil {
		log.Fatal("Cannot connect to Woo device: ", err)
		return
	}

	<-make(chan struct{})
}
