package main

func ProcessECGPacket() {

}

func MonitorBPM(bpmChan <-chan FHIRBPMObservation, resultChan chan<- bool) {
	for obs := range bpmChan {
		if obs.ValueQuantity.Value < 40 || obs.ValueQuantity.Value > 180 {
			resultChan <- true
		} else {
			resultChan <- false
		}
	}
}

func main() {

}
