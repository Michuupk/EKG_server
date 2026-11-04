package main_test

import (
	main "ekg_server/cmd/server"
	"testing"
)

func TestECGPacket_Validate(t *testing.T) {
	tests := []struct {
		name    string
		packet  main.ECGPacket
		wantErr bool
	}{
		{
			name: "Poprawna paczka z ESP32 (surowe próbki AD8232)",
			packet: main.ECGPacket{
				DeviceID:       "ESP32-DEV-01",
				PatientID:      "PAT-10203",
				SequenceNumber: 1,
				Timestamp:      1700000000,
				Samples:        []float64{0.12, 0.15, 0.28, -0.05},
			},
			wantErr: false,
		},
		{
			name: "Brak device_id",
			packet: main.ECGPacket{
				DeviceID:       "",
				PatientID:      "PAT-10203",
				SequenceNumber: 1,
				Timestamp:      1700000000,
				Samples:        []float64{0.12},
			},
			wantErr: true,
		},
		{
			name: "Brak patient_id",
			packet: main.ECGPacket{
				DeviceID:       "ESP32-DEV-01",
				PatientID:      "",
				SequenceNumber: 1,
				Timestamp:      1700000000,
				Samples:        []float64{0.12},
			},
			wantErr: true,
		},
		{
			name: "Nieprawidłowy timestamp (<= 0)",
			packet: main.ECGPacket{
				DeviceID:       "ESP32-DEV-01",
				PatientID:      "PAT-10203",
				SequenceNumber: 1,
				Timestamp:      0,
				Samples:        []float64{0.12},
			},
			wantErr: true,
		},
		{
			name: "Pusta tablica próbek (brak sygnału)",
			packet: main.ECGPacket{
				DeviceID:       "ESP32-DEV-01",
				PatientID:      "PAT-10203",
				SequenceNumber: 1,
				Timestamp:      1700000000,
				Samples:        []float64{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.packet.Validate()
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}
