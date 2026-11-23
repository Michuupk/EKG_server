package main

func ProcessECGPacket(packet ECGPacket) bool {
	for _, val := range packet.Packet {
		if val > 180.0 { // Przykładowy próg wartości sygnału
			return true
		}
	}
	return false
}

func MonitorBPM(packetChan <-chan ECGPacket, resultChan chan<- bool) {
	for packet := range packetChan {
		hasAnomaly := ProcessECGPacket(packet)
		resultChan <- hasAnomaly
	}
}

func main() {

}
