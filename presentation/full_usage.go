// Generated with Gemini-3.1-PRO
package main

import (
	"fmt"
	"time"
)

type CriticalAlert struct {
	PatientID int
	BPM       int
	Message   string
}

// Goroutine B: Moduł wywoływania alarmu (Odbiorca)
func EmergencyHandler(alertChan chan CriticalAlert) {
	for alert := range alertChan {
		// Odbiór natychmiast odblokowuje Goroutine A
		fmt.Printf("[ALARM CRITICAL] Pacjent %d | BPM: %d | Status: %s\n", 
			alert.PatientID, alert.BPM, alert.Message)
		// Rozpoczęcie procedury powiadamiania lekarza...
	}
}

// Goroutine A: Analizator sygnału EKG (Nadawca)
func ReadTelemetry(alertChan chan CriticalAlert) {
	for {
		bpm := fetchCurrentBPM() // Odczyt z czujnika

		if bpm > 180 || bpm < 40 {
			alert := CriticalAlert{PatientID: 402, BPM: bpm, Message: "bpm OutOfRange"}
			
			// SYNHRONICZNE PRZEKAZANIE (RENDEZVOUS):
			// Goroutine A ZATRZYMUJE pracę i czeka, aż EmergencyHandler 
			// fizycznie odczyta ten obiekt z kanału.
			alertChan <- alert 
			
			logAlertDispatched() // Wykonuje się DOPIERO PO odbiorze przez Goroutine B
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	// Tworzenie niebuforowanego kanału w pamięci RAM
	alertChan := make(chan CriticalAlert)

	go EmergencyHandler(alertChan) // Uruchomienie odbiorcy w tle
	go ReadTelemetry(alertChan)    // Uruchomienie analizatora w tle

	select {} // Blokada głównej funkcji main (utrzymanie programu)
}

func fetchCurrentBPM() int { return 195 }
func logAlertDispatched()  { fmt.Println("[LOG] Alarm przekazany.") }