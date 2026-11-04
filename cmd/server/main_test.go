package main

import "testing"

func TestMonitorBPM(t *testing.T) {
	tests := []struct {
		name     string
		bpmValue int
		want     bool
	}{
		{"Prawidłowe tętno spoczynkowe", 75, false},
		{"Wykrycie tachykardii", 190, true},
		{"Wykrycie bradykardii", 35, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bpmChan := make(chan int)
			resultChan := make(chan bool)

			go MonitorBPM(bpmChan, resultChan)

			bpmChan <- tt.bpmValue

			got := <-resultChan
			if got != tt.want {
				t.Errorf("MonitorBPM() zgłosił = %v, oczekiwano = %v", got, tt.want)
			}

			close(bpmChan)
		})
	}
}
