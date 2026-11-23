package main

import (
	"testing"
)

func TestProcessPacket(t *testing.T) {
	tests := []struct {
		name    string
		packet  ECGPacket // Używamy gotowej struktury z models.go
		wantErr bool
	}{
		{
			"Prawidłowe wartości",
			ECGPacket{
				DeviceID:   "DEV-ESP32-01",
				Seq:        100,
				Packet:     []float64{0.12, 0.15, 0.85, -0.05},
				PatientID:  "PACJENT-999",
				IsBuffered: false,
			},
			false,
		},
		{
			"Sygnał BPM poza wartością",
			ECGPacket{
				DeviceID:   "DEV-ESP32-01",
				Seq:        100,
				Packet:     []float64{10000.12, 0.15, 0.85, -0.05},
				PatientID:  "PACJENT-999",
				IsBuffered: false,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bpmChan := make(chan ECGPacket)
			resultChan := make(chan bool)

			go MonitorBPM(bpmChan, resultChan)

			bpmChan <- tt.packet

			got := <-resultChan
			if got != tt.wantErr {
				t.Errorf("MonitorBPM() zgłosił = %v, oczekiwano = %v", got, tt.wantErr)
			}

			close(bpmChan)
		})
	}
}
