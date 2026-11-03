// Generated with Gemini-3.1-PRO
package main

import "fmt"

func main() {
    ch := make(chan string, 2) // 1. Tworzenie bufora na 2 elementy

    ch <- "Kawa"               // 2. Pierwsza paczka
    ch <- "Herbata"            // 3. Druga paczka

    fmt.Println(<-ch)          // 4. Pierwszy odczyt
    fmt.Println(<-ch)          // 5. Drugi odczyt
}