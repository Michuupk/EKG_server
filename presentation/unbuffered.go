// Generated with Gemini-3.1-PRO
package main

import "fmt"

func main() {
    ch := make(chan string) // 1. Tworzenie kanału niebuforowanego

    go func() {
        ch <- "Pałeczka"     // 2. Nadawca wysyła wiadomość
    }()

    msg := <-ch             // 3. Odbiorca odbiera wiadomość
    fmt.Println(msg)
}