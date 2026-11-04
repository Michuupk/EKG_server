package main

import "errors"

// BPMCodedConcept reprezentuje zakodowaną informację o rodzaju pomiaru (np. "Heart rate")
type BPMCodedConcept struct {
	Text string `json:"text"`
}

// BPMQuantity reprezentuje zmierzoną wartość wraz z jednostką
type BPMQuantity struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

// FHIRBPMObservation to główna struktura zgodna z FHIR dla pomiaru tętna
type FHIRBPMObservation struct {
	ResourceType  string          `json:"resourceType"`
	Status        string          `json:"status"`
	Code          BPMCodedConcept `json:"code"`
	ValueQuantity BPMQuantity     `json:"valueQuantity"`
}

type ECGPacket struct {
	DeviceID       string             `json:"device_id"`
	PatientID      string             `json:"patient_id"`
	SequenceNumber uint64             `json:"sequence_number"`
	Timestamp      int64              `json:"timestamp"`
	IsBuffered     bool               `json:"is_buffered"`
	Observation    FHIRBPMObservation `json:"observation"`
}

func (p ECGPacket) Validate() error {
	if p.DeviceID == "" {
		return errors.New("brak device_id")
	}
	if p.Timestamp <= 0 {
		return errors.New("nieprawidłowy timestamp")
	}
	if p.Observation.ResourceType != "Observation" {
		return errors.New("brak lub błędny resourceType w danych FHIR")
	}
	if p.Observation.ValueQuantity.Value <= 0 {
		return errors.New("wartość BPM musi być dodatnia")
	}
	return nil
}
