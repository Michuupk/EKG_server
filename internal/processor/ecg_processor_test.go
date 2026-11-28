package processor

import (
	"context"
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
				if tt.expectedAnomalies[i] != point.IsAnomaly {
					t.Errorf("Błędne wykrycie anomalii")
				}
			}
		})
	}
}

func TestDispatchPoints(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		points          []domain.ECGDataPoint
		storageChan     chan<- domain.ECGDataPoint
		alertChan       chan<- domain.ECGDataPoint
		expectedStorage int
		expectedAlerts  int
	}{
		{
			name: "Paczka z 2 prawidłowymi i 1 anomalią",
			points: []domain.ECGDataPoint{
				{Value: 512, IsAnomaly: false},
				{Value: -1, IsAnomaly: true},
				{Value: 490, IsAnomaly: false},
			},
			expectedStorage: 3,
			expectedAlerts:  1,
		},
		{
			name: "Paczka bez anomalii",
			points: []domain.ECGDataPoint{
				{Value: 512, IsAnomaly: false},
				{Value: 530, IsAnomaly: false},
			},
			expectedStorage: 2,
			expectedAlerts:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageChan := make(chan domain.ECGDataPoint, 10)
			alertChan := make(chan domain.ECGDataPoint, 10)

			err := DispatchPoints(ctx, tt.points, storageChan, alertChan)
			if err != nil {
				t.Fatalf("Oczekiwano braku błędu, otrzymano: %v", err)
			}

			if len(storageChan) != tt.expectedStorage {
				t.Errorf("storageChan: oczekiwano %d punktów, otrzymano %d", tt.expectedStorage, len(storageChan))
			}

			if len(alertChan) != tt.expectedAlerts {
				t.Errorf("alertChan: oczekiwano %d alertów, otrzymano %d", tt.expectedAlerts, len(alertChan))
			}
		})
	}
}

func TestStartProcessor(t *testing.T) {
	tests := []struct {
		name            string
		inputPackets    []domain.ECGPacket
		expectedStorage int
		expectedAlerts  int
	}{
		{
			name: "Mieszane paczki: 1 poprawna (3 punkty w tym 1 anomalia) i 1 błędna walidacja",
			inputPackets: []domain.ECGPacket{
				{
					DeviceID:  "ESP32-01",
					PatientID: "PAT-01",
					Seq:       1,
					Packet:    []int{500, -1, 600}, // 3 punkty, 1 anomalia (-1)
				},
				{
					DeviceID:  "", // Brak DeviceID -> Validate() zwróci błąd!
					PatientID: "PAT-01",
					Seq:       2,
					Packet:    []int{500, 500},
				},
			},
			expectedStorage: 3, // Tylko 3 punkty z pierwszej (poprawnej) paczki
			expectedAlerts:  1, // 1 anomalia z pierwszej paczki
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			inChan := make(chan domain.ECGPacket, 10)
			storageChan := make(chan domain.ECGDataPoint, 10)
			alertChan := make(chan domain.ECGDataPoint, 10)

			for _, pkt := range tt.inputPackets {
				inChan <- pkt
			}
			close(inChan)

			err := StartProcessor(ctx, inChan, storageChan, alertChan)
			if err != nil {
				t.Fatalf("Oczekiwano braku błędu procesora, otrzymano: %v", err)
			}

			if len(storageChan) != tt.expectedStorage {
				t.Errorf("storageChan: oczekiwano %d punktów, otrzymano %d", tt.expectedStorage, len(storageChan))
			}

			if len(alertChan) != tt.expectedAlerts {
				t.Errorf("alertChan: oczekiwano %d alertów, otrzymano %d", tt.expectedAlerts, len(alertChan))
			}
		})
	}
}
