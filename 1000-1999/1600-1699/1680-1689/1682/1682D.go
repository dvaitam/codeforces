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

    var T int
    fmt.Fscan(reader, &T)
    for T > 0 {
        T--
        var N int
        fmt.Fscan(reader, &N)
        var s string
        fmt.Fscan(reader, &s)
        odd := 0
        for i := 0; i < N; i++ {
            if s[i] == '1' {
                odd++
            }
        }
        if odd == 0 || odd%2 == 1 {
            writer.WriteString("NO\n")
            continue
        }
        writer.WriteString("YES\n")
        dx := 0
        // rotate so that last position has '1'
        if s[N-1] != '1' {
            for s[dx] != '1' {
                dx++
            }
            dx++
        }
        // build edges
        for i := 0; i < N; i++ {
            j := i
            if i != 0 {
                u := dx + 1
                v := (i+dx)%N + 1
                fmt.Fprintln(writer, u, v)
            }
            for s[(j+dx)%N] == '0' {
                u := (j+dx)%N + 1
                v := (j+dx+1)%N + 1
                fmt.Fprintln(writer, u, v)
                j++
            }
            i = j
        }
    }
}
