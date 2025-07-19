package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var n, k int64
    var l, v, b float64
    if _, err := fmt.Fscan(reader, &n, &l, &v, &b, &k); err != nil {
        return
    }
    g := n / k
    if n%k > 0 {
        g++
    }
    fg := float64(g)
    numerator := l*v + (2*b*fg - b) * l
    denominator := (2*b*fg - b) * v + b*b
    ans := numerator / denominator
    fmt.Fprintf(writer, "%.10f", ans)
}
