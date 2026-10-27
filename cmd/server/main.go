package main

import (
	"fmt"
	"log"
	"net/http"

	// Ufamy i korzystamy bezpośrednio z zewnętrznej biblioteki
	"github.com/gorilla/websocket"
)

// Konfiguracja "Upgradera" (zmiana HTTP na WebSocket)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Zezwalamy na połączenia z dowolnego źródła
	},
}

// Funkcja obsługująca pojedyncze połączenie WebSocket
func handleECGConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade połączenia
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Błąd podczas upgrade'u połączenia:", err)
		return
	}
	defer ws.Close()

	log.Printf("Nowy klient podłączony z adresu: %s", r.RemoteAddr)

	// Główna pętla odbierająca paczki
	for {
		messageType, payload, err := ws.ReadMessage()
		if err != nil {
			log.Println("Klient rozłączył się lub wystąpił błąd odczytu:", err)
			break
		}

		fmt.Printf("Otrzymano paczkę EKG (%d bajtów): %s\n", len(payload), string(payload))

		// Odsyłamy potwierdzenie do klienta
		ackMessage := []byte("Serwer: Paczka EKG odebrana.")
		err = ws.WriteMessage(messageType, ackMessage)
		if err != nil {
			log.Println("Błąd wysyłania potwierdzenia:", err)
			break
		}
	}
}

func main() {
	// Przypinamy endpoint "/ws" do naszej funkcji
	http.HandleFunc("/ws", handleECGConnection)

	fmt.Println("Serwer EKG nasłuchuje na porcie :8080...")

	// Start serwera HTTP (pod spodem automatycznie używa Goroutines dla nowych żądań)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Krytyczny błąd serwera: ", err)
	}
}
