package main

import (
	"os/exec"
)

func main() {
	cmd := exec.Command("/bin/ash")//, "-c", "ash")
    if err := cmd.Run(); err != nil {
        panic(err)
    }
}
