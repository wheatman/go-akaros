package main

import "fmt"
import "time"

func main() {
	t := time.Date(2013, time.July, 19, 4, 38, 0, 0, time.UTC)
    fmt.Printf("Hello from Go!\n")
    fmt.Printf("First Contact: %v\n", t.Local())
}

