package main

import (
	"errors"
	"log"
	"strings"
	"time"
	"tinygo.org/x/bluetooth"
)

func FindWooDevice(timeout time.Duration, adapter *bluetooth.Adapter) (bluetooth.Address, error) {
	resultChan := make(chan bluetooth.Address, 1)
	go func() {
		log.Println("Start scan...")
		err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
			if strings.HasPrefix(result.LocalName(), "WOO-") {
				log.Printf("Found Woo Device %s(%s)", result.LocalName(), result.Address)
				adapter.StopScan()
				log.Println("Stop scan")
				resultChan <- result.Address
			}
		})
		if err != nil {
			log.Fatalf("Scan error: %s", err)
		}
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case <-time.After(timeout):
		adapter.StopScan()
		log.Println("Stop scan")
		return bluetooth.Address{}, errors.New("timeout")
	}
}

func ConnectToDevice(adapter *bluetooth.Adapter, reconnect bool, address bluetooth.Address, callback func(device bluetooth.Device)) error {
	adapter.SetConnectHandler(func(device bluetooth.Device, connected bool) {
		if device.Address != address {
			return
		}

		if connected {
			log.Println("Connected")
			callback(device)
		} else if reconnect {
			log.Println("Reconnecting...")
			go func() {
				_, err := adapter.Connect(address, bluetooth.ConnectionParams{})
				if err != nil {
					log.Fatalf("Cannot reconnect: %s", err)
				}
			}()
		}

	})
	log.Printf("Connecting to %s", address)
	_, err := adapter.Connect(address, bluetooth.ConnectionParams{})
	if err != nil {
		return err
	}

	return nil
}

func SubscribeWooEvents(device bluetooth.Device, callback func(data []byte)) error {
	wooServiceId, err := bluetooth.ParseUUID("00000001-0000-0000-0000-000000000080")

	srvcs, err := device.DiscoverServices([]bluetooth.UUID{wooServiceId})
	if err != nil {
		log.Fatalf("Cannot discover services: %s", err)
	}

	wooService := srvcs[0]

	wooEventsCharacteristicId, _ := bluetooth.ParseUUID("00000003-0000-0000-0000-000000000080")
	chars, err := wooService.DiscoverCharacteristics([]bluetooth.UUID{wooEventsCharacteristicId})
	if err != nil {
		log.Fatalf("Cannot discover characteristics: %s", err)
		return err
	}

	for _, char := range chars {
		err = char.EnableNotifications(func(event []byte) {
			callback(event)
		})
		if err != nil {
			log.Fatalf("Cannot enable notifications: %s", err)
			return err
		}
	}

	return nil
}
