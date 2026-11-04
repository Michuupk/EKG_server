package main_test

import (
	main "ekg_server/cmd/server"
	"testing"
)

func TestECGPacket_Validate(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		packet  main.ECGPacket
		wantErr bool
	}{
		{
			name: "Brak device_id w pakiecie",
			packet: main.ECGPacket{
				PatientID:      "patient_001",
				SequenceNumber: 1,
				Timestamp:      17000000,
				IsBuffered:     false,
				Observation: main.FHIRBPMObservation{
					ResourceType: "Observation",
					Status:       "final",
					Code: main.BPMCodedConcept{
						Text: "Heart rate",
					},
					ValueQuantity: main.BPMQuantity{
						Value: 75,
						Unit:  "beats/minute",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			gotErr := tt.packet.Validate()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Validate() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Validate() succeeded unexpectedly")
			}
		})
	}
}
