package main

import "errors"

// BPMCodedConcept reprezentuje kodowanie pomiaru według standardu FHIR
type BPMCodedConcept struct {
	Text string `json:"text"`
}

// BPMQuantity reprezentuje wartość i jednostkę pomiaru
type BPMQuantity struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

// FHIRBPMObservation to struktura FHIR tworzona na serwerze po przeliczeniu sygnału EKG
type FHIRBPMObservation struct {
	ResourceType  string          `json:"resourceType"`
	Status        string          `json:"status"`
	Code          BPMCodedConcept `json:"code"`
	ValueQuantity BPMQuantity     `json:"valueQuantity"`
}

// ECGPacket to koperta telemetryczna odbierana bezpośrednio z ESP32
type ECGPacket struct {
	DeviceID       string             `json:"device_id"`
	PatientID      string             `json:"patient_id"`
	SequenceNumber uint64             `json:"sequence_number"`
	Timestamp      int64              `json:"timestamp"`
	IsBuffered     bool               `json:"is_buffered"`
	Samples        []float64          `json:"samples"`
	Observation    FHIRBPMObservation `json:"observation,omitempty"` // Tworzone na serwerze
}

// Validate weryfikuje kompletność surowych danych przychodzących z mikrokontrolera
func (p ECGPacket) Validate() error {
	if p.DeviceID == "" {
		return errors.New("brak device_id")
	}
	if p.PatientID == "" {
		return errors.New("brak patient_id")
	}
	if p.Timestamp <= 0 {
		return errors.New("nieprawidłowy timestamp")
	}
	if len(p.Samples) == 0 {
		return errors.New("wymagana jest co najmniej 1 próbka w tablicy samples")
	}
	return nil
}
