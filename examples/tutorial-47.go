package main

import "fmt"
import "math/cmplx"

func Cbrt(x complex128) complex128 {
    z := x
    for _ = range make([]int, 10000) {
        z = z - (cmplx.Pow(z,3) - x)/(3*cmplx.Pow(z,2))
    }
    return z
}

func main() {
    fmt.Println(Cbrt(2))
}
