package domain

import (
	"errors"
	"time"
)

var (
	ErrEmptyDeviceID  = errors.New("brak device_id")
	ErrEmptyPatientID = errors.New("brak patient_id")
	ErrEmptyPacket    = errors.New("tablica packet nie może być pusta")
)

type ECGPacket struct {
	DeviceID   string `json:"device_id"`
	Seq        uint64 `json:"seq"`
	Packet     []int  `json:"packet"`
	PatientID  string `json:"patient_id"`
	IsBuffered bool   `json:"is_buffered"`
}

// Validate sprawdza obecność wymaganych identyfikatorów oraz próbek sygnału
func (p ECGPacket) Validate() error {
	if p.DeviceID == "" {
		return ErrEmptyDeviceID
	}
	if p.PatientID == "" {
		return ErrEmptyPatientID
	}
	if len(p.Packet) == 0 {
		return ErrEmptyPacket
	}
	return nil
}

type ECGDataPoint struct {
	Timestamp time.Time
	DeviceID  string
	PatientID string
	Seq       uint64
	Value     int
	IsAnomaly bool // true dla wartości < 0 lub > 1023
}
