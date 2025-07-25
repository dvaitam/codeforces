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

    // TODO: Implement solution for fish population problem (1023G).
    // Placeholder: not implemented
    var n int
    if _, err := fmt.Fscan(reader, &n); err != nil {
        return
    }
    // skip edges
    for i := 0; i < n-1; i++ {
        var u, v, l int
        fmt.Fscan(reader, &u, &v, &l)
    }
    var k int
    fmt.Fscan(reader, &k)
    // skip observations
    totalFish := 0
    for i := 0; i < k; i++ {
        var d, f, p int
        fmt.Fscan(reader, &d, &f, &p)
        totalFish += f
    }
    // As a fallback, output total observed fish
    fmt.Fprintln(writer, totalFish)
}
