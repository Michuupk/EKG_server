package main

func MonitorBPM(bpmChan <-chan int, resultChan chan<- bool) {
	for bpm := range bpmChan {
		if bpm < 40 || bpm > 180 {
			resultChan <- true
		} else {
			resultChan <- false
		}
	}
}

func main() {

}
