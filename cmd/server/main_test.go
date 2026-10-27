package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestECGWebSocketConnection(t *testing.T) {
	// 1. Tworzymy wirtualny serwer testowy (httptest)
	// Przekazujemy mu bezpośrednio naszą funkcję z main.go, którą chcemy przetestować.
	server := httptest.NewServer(http.HandlerFunc(handleECGConnection))
	// defer zapewni wyłączenie serwera po zakończeniu testu.
	defer server.Close()

	// 2. Zamiana URL z HTTP na WS, ponieważ httptest.NewServer uruchamia się po "http"
	// strings.Replace podmienia protokół na potrzeby klienta WebSocket.
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// 3. Połączenie testowego klienta z wirtualnym serwerem
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Nie udało się połączyć z WebSocketem: %v", err)
	}
	defer ws.Close()

	closeMessage := websocket.FormatCloseMessage(3000, "Closed by user")
	err = ws.WriteMessage(websocket.CloseMessage, closeMessage)
	if err != nil {
		t.Fatalf("Błąd podczas wysyłania sygnału zamknięcia 3000: %v", err)
	}

	// 4. Definicja paczki testowej
	testPayload := []byte(`{"id": 1, "lead_v1": 1.5}`)

	// Wysyłamy dane z testowego klienta do testowanego serwera
	err = ws.WriteMessage(websocket.TextMessage, testPayload)
	if err != nil {
		t.Fatalf("Błąd podczas wysyłania wiadomości: %v", err)
	}

	// 5. Odczytujemy odpowiedź serwera (nasz serwer odsyła "Serwer: Paczka EKG odebrana.")
	messageType, responsePayload, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("Błąd podczas odczytu odpowiedzi: %v", err)
	}

	// 6. Aserty i weryfikacja
	if messageType != websocket.TextMessage {
		t.Errorf("Oczekiwano wiadomości tekstowej, otrzymano typ: %d", messageType)
	}

	expectedResponse := "Serwer: Paczka EKG odebrana."
	if string(responsePayload) != expectedResponse {
		t.Errorf("Oczekiwano: %s, otrzymano: %s", expectedResponse, string(responsePayload))
	}

	// Wysyłamy ramkę kontrolną do serwera
	err = ws.WriteMessage(websocket.CloseMessage, closeMessage)
	if err != nil {
		t.Fatalf("Błąd podczas wysyłania sygnału zamknięcia: %v", err)
	}
}
