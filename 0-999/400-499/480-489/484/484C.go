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

    var s string
    fmt.Fscan(reader, &s)
    S := []byte(s)
    n := len(S)
    var m int
    fmt.Fscan(reader, &m)
    tmp := make([]byte, n)
    for qi := 0; qi < m; qi++ {
        var k, d int
        fmt.Fscan(reader, &k, &d)
        // apply d-sorting to each substring of length k
        for i := 0; i <= n-k; i++ {
            p := 0
            // collect characters in groups by modulo
            for g := 0; g < d; g++ {
                for j := i + g; j < i+k; j += d {
                    tmp[p] = S[j]
                    p++
                }
            }
            // write back sorted substring
            for t := 0; t < p; t++ {
                S[i+t] = tmp[t]
            }
        }
        // print current string
        writer.Write(S)
        writer.WriteByte('\n')
    }
}
