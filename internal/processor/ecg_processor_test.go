package processor

import (
	"ekg_server/internal/domain"
	"testing"
	"time"
)

func TestUnpackPacket(t *testing.T) {
	baseTime := time.Date(2026, 11, 25, 14, 0, 0, 0, time.UTC)

	tests := []struct {
		name              string
		inputPacket       domain.ECGPacket
		expectedLength    int
		expectedAnomalies []bool // np. [false, true, false]
	}{
		{
			name: "Prawidłowa paczka",
			inputPacket: domain.ECGPacket{
				DeviceID:  "ESP32-01",
				Seq:       1,
				Packet:    []int{512, 530, 610, 490},
				PatientID: "550e8400-e29b-41d4-a716-446655440000",
			},
			expectedLength:    4,
			expectedAnomalies: []bool{false, false, false, false},
		},
		{
			name: "Paczka z odklejoną elektrodą (wartości -1 oraz 1050)",
			inputPacket: domain.ECGPacket{
				DeviceID:  "ESP32-01",
				Seq:       2,
				Packet:    []int{512, -1, 1050, 490},
				PatientID: "550e8400-e29b-41d4-a716-446655440000",
			},
			expectedLength:    4,
			expectedAnomalies: []bool{false, true, true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UnpackPacket(tt.inputPacket, baseTime)
			if len(result) != tt.expectedLength {
				t.Errorf("Funkcja zwróciła nieodpowiednią długość paczki")
			}
			for i, point := range result {
				if tt.expectedAnomalies[i] == point.IsAnomaly {
					t.Errorf("Błędne wykrycie anomalii")
				}
			}
		})
	}
}
