package main

import "testing"

func TestECGPacket_Validate(t *testing.T) {
	tests := []struct {
		name    string
		packet  ECGPacket
		wantErr bool
	}{
		{
			name: "Poprawna paczka danych",
			packet: ECGPacket{
				DeviceID:   "DEV-ESP32-01",
				Seq:        100,
				Packet:     []float64{0.12, 0.15, 0.85, -0.05},
				PatientID:  "PACJENT-999",
				IsBuffered: false,
			},
			wantErr: false,
		},
		{
			name: "Brak device_id",
			packet: ECGPacket{
				DeviceID:  "",
				Seq:       101,
				Packet:    []float64{0.12},
				PatientID: "PACJENT-999",
			},
			wantErr: true,
		},
		{
			name: "Brak patient_id",
			packet: ECGPacket{
				DeviceID:  "DEV-ESP32-01",
				Seq:       102,
				Packet:    []float64{0.12},
				PatientID: "",
			},
			wantErr: true,
		},
		{
			name: "Pusta tablica próbek EKG",
			packet: ECGPacket{
				DeviceID:  "DEV-ESP32-01",
				Seq:       103,
				Packet:    []float64{}, // len(Packet) == 0
				PatientID: "PACJENT-999",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.packet.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}