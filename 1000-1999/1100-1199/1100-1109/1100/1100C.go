package main

import (
    "bufio"
    "fmt"
    "math"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    var n int
    var r float64
    if _, err := fmt.Fscan(reader, &n, &r); err != nil {
        return
    }
    angle := math.Pi / float64(n)
    s := math.Sin(angle)
    R := r * s / (1 - s)
    fmt.Printf("%.10f\n", R)
}
