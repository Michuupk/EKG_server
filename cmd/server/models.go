package main

import "errors"

type ECGPacket struct {
	DeviceID   string    `json:"device_id"`
	Seq        uint64    `json:"seq"`
	Packet     []float64 `json:"packet"`
	PatientID  string    `json:"patient_id"`
	IsBuffered bool      `json:"is_buffered"`
}

// Validate sprawdza obecność wymaganych identyfikatorów oraz próbek sygnału
func (p ECGPacket) Validate() error {
	if p.DeviceID == "" {
		return errors.New("brak device_id")
	}
	if p.PatientID == "" {
		return errors.New("brak patient_id")
	}
	if len(p.Packet) == 0 {
		return errors.New("tablica packet nie może być pusta")
	}
	return nil
}