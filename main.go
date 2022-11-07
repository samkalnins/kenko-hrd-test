package main

import (
	"encoding/binary"
	"flag"
	"log"
	"time"

	"tinygo.org/x/bluetooth"
)

var (
	flagScan        = flag.Bool("scan", false, "Scan for bluetooth devices")
	flagScanTimeout = flag.Int("scan_timeout", 10, "Seconds to scan for")

	flagStream = flag.String("stream", "", "Connect and stream from a given device UUID")

	adapter = bluetooth.DefaultAdapter

	heartRateServiceUUID        = bluetooth.ServiceUUIDHeartRate
	heartRateCharacteristicUUID = bluetooth.CharacteristicUUIDHeartRateMeasurement
)

func bluetooth_packet_parse(b []byte) {
	var hrm_bpm float64
	var rr_intervals []float64
	num_rr_values := 0

	// http://www.mariam.qa/post/hr-ble/
	hrm_data_flag := b[0]

	is_uint16 := hrm_data_flag&0x1 != 0
	ee_present := hrm_data_flag&0x8 != 0
	rr_present := hrm_data_flag&0x16 != 0

	hrm_bytesize := 1

	if is_uint16 {
		hrm_bytesize = 2
		hrm_bpm = float64(binary.LittleEndian.Uint16(b[1:3]))
	} else {
		hrm_bpm = float64(uint8(b[1]))
	}

	ee_bytesize := 0
	if ee_present {
		ee_bytesize = 2
	}

	if rr_present {
		num_rr_values = (len(b) - 1 - hrm_bytesize - ee_bytesize) / 2

		start_pos := len(b) - num_rr_values*2
		for start_pos < len(b) {
			rr_ms := float64(binary.LittleEndian.Uint16(b[start_pos:start_pos+2])) / 1024 * 1000
			rr_intervals = append(rr_intervals, rr_ms)
			start_pos += 2
		}
	}
	log.Printf("bpm=%0f num_rr_intervals=%d rr_intervals_ms=%v", hrm_bpm, num_rr_values, rr_intervals)
	//fmt.Printf("%d,%0f\n", time.Now().Unix(), hrm_bpm)
}

func streamFromDevice(addr bluetooth.Addresser, name string) {
	var device *bluetooth.Device
	var err error

	log.Printf("Streaming from %s [%s]\n", addr.String(), name)
	for {
		device, err = adapter.Connect(addr, bluetooth.ConnectionParams{ConnectionTimeout: 10000})
		if err == nil {
			log.Printf("Connected to %s [%s]\n", addr.String(), name)
			break
		} else {
			log.Println(err)
		}
		time.Sleep(5 * time.Second)
	}

	srvcs, err := device.DiscoverServices([]bluetooth.UUID{heartRateServiceUUID})
	if err != nil || len(srvcs) == 0 {
		log.Fatalf("Could not find heart rate service for %s", name)
	}

	chars, err := srvcs[0].DiscoverCharacteristics([]bluetooth.UUID{heartRateCharacteristicUUID})
	if err != nil || len(chars) == 0 {
		log.Fatalf("Could not find heart rate characteristic for %s", name)
	}

	chars[0].EnableNotifications(func(buf []byte) {
		bluetooth_packet_parse(buf)
	})
}

func scan(adapter *bluetooth.Adapter) map[bluetooth.Addresser]string {
	found_devices := make(map[bluetooth.Addresser]string)

	start_time := time.Now()
	adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {

		if name, ok := found_devices[result.Address]; !ok {
			found_devices[result.Address] = result.LocalName()
		} else {
			if name == "" && name != result.LocalName() {
				// LocalName isn't always broadcast
				found_devices[result.Address] = result.LocalName()
			}
		}

		if time.Since(start_time) > time.Duration(time.Duration(*flagScanTimeout)*time.Second) {
			adapter.StopScan()
		}
	})
	return found_devices
}

func main() {
	flag.Parse()

	err := adapter.Enable()
	if err != nil {
		log.Fatal("Could not enable bluetooth stack")
	}

	if *flagScan {
		log.Println("Scanning for bluetooth devices:")
		for addr, name := range scan(adapter) {
			log.Printf("Found %s %s\n", addr.String(), name)
		}
	}

	if *flagStream != "" {
		for addr, name := range scan(adapter) {
			if *flagStream == addr.String() {
				streamFromDevice(addr, name)
				select {}
			}
		}
	}
}
