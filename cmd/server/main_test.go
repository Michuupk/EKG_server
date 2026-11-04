package main

import (
	"testing"
)

func TestMonitorBPM(t *testing.T) {
	tests := []struct {
		name        string
		observation FHIRBPMObservation
		want        bool
	}{
		{
			"Prawidłowe tętno FHIR",
			FHIRBPMObservation{
				ResourceType: "Observation",
				Status:       "final",
				Code: BPMCodedConcept{
					Text: "Heart rate",
				},
				ValueQuantity: BPMQuantity{
					Value: 75,
					Unit:  "beats/minute",
				},
			},
			false},
		{
			"Tachykardia FHIR",
			FHIRBPMObservation{
				ResourceType: "Observation",
				Status:       "final",
				Code: BPMCodedConcept{
					Text: "Heart rate",
				},
				ValueQuantity: BPMQuantity{
					Value: 190,
					Unit:  "beats/minute",
				},
			},
			true},
		{
			"Bradykardia FHIR",
			FHIRBPMObservation{
				ResourceType: "Observation",
				Status:       "final",
				Code: BPMCodedConcept{
					Text: "Heart rate",
				},
				ValueQuantity: BPMQuantity{
					Value: 35,
					Unit:  "beats/minute",
				},
			},
			true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bpmChan := make(chan FHIRBPMObservation)
			resultChan := make(chan bool)

			go MonitorBPM(bpmChan, resultChan)

			bpmChan <- tt.observation

			got := <-resultChan
			if got != tt.want {
				t.Errorf("MonitorBPM() zgłosił = %v, oczekiwano = %v", got, tt.want)
			}

			close(bpmChan)
		})
	}
}
