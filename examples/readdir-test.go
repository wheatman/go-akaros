package main

import (
	"os"
	"fmt"
)

func main() {
	file, _ := os.Open(".")
	names, _ := file.Readdirnames(100)
	fmt.Printf("%v\n", names)
}
