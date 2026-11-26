package domain

import "testing"

func TestECGPacket_Validate(t *testing.T) {
	tests := []struct {
		name    string
		packet  ECGPacket
		wantErr bool
	}{
		{
			name: "Poprawna paczka z czujnika AD8232",
			packet: ECGPacket{
				DeviceID:   "ESP32-01",
				Seq:        1,
				Packet:     []int{512, 530, 610, 490},
				PatientID:  "550e8400-e29b-41d4-a716-446655440000",
				IsBuffered: false,
			},
			wantErr: false,
		},
		{
			name: "Błąd pomiarowy AD8232 (wartość poza 0-1023) przechodzi walidację struktury",
			packet: ECGPacket{
				DeviceID:   "ESP32-01",
				Seq:        2,
				Packet:     []int{512, -1, 1050},
				PatientID:  "550e8400-e29b-41d4-a716-446655440000",
				IsBuffered: false,
			},
			wantErr: false,
		},
		{
			name: "Brak device_id",
			packet: ECGPacket{
				DeviceID:  "",
				Seq:       3,
				Packet:    []int{512},
				PatientID: "550e8400-e29b-41d4-a716-446655440000",
			},
			wantErr: true,
		},
		{
			name: "Brak patient_id (wymóg RODO)",
			packet: ECGPacket{
				DeviceID:  "ESP32-01",
				Seq:       4,
				Packet:    []int{512},
				PatientID: "",
			},
			wantErr: true,
		},
		{
			name: "Pusta tablica próbek",
			packet: ECGPacket{
				DeviceID:  "ESP32-01",
				Seq:       5,
				Packet:    []int{},
				PatientID: "550e8400-e29b-41d4-a716-446655440000",
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
