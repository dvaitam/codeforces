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

    var t int
    if _, err := fmt.Fscan(reader, &t); err != nil {
        return
    }
    for ; t > 0; t-- {
        var n int
        fmt.Fscan(reader, &n)
        // first line: n '(' then n ')'
        for i := 0; i < n; i++ {
            writer.WriteByte('(')
        }
        for i := 0; i < n; i++ {
            writer.WriteByte(')')
        }
        writer.WriteByte('\n')

        // subsequent lines
        for i := 0; i < n-1; i++ {
            // prefix pairs () repeated i+1 times
            for j := 0; j <= i; j++ {
                writer.WriteString("()")
            }
            // remaining parentheses
            rem := n - (i + 1)
            for j := 0; j < rem; j++ {
                writer.WriteByte('(')
            }
            for j := 0; j < rem; j++ {
                writer.WriteByte(')')
            }
            writer.WriteByte('\n')
        }
    }
}
