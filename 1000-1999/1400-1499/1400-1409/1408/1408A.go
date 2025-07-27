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
    fmt.Fscan(reader, &t)
    for tc := 0; tc < t; tc++ {
        var n int
        fmt.Fscan(reader, &n)
        a := make([]int, n)
        b := make([]int, n)
        c := make([]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(reader, &a[i])
        }
        for i := 0; i < n; i++ {
            fmt.Fscan(reader, &b[i])
        }
        for i := 0; i < n; i++ {
            fmt.Fscan(reader, &c[i])
        }
        p := make([]int, n)
        // choose first element
        p[0] = a[0]
        // choose middle elements ensuring no conflict with previous
        for i := 1; i < n-1; i++ {
            if a[i] != p[i-1] {
                p[i] = a[i]
            } else if b[i] != p[i-1] {
                p[i] = b[i]
            } else {
                p[i] = c[i]
            }
        }
        // choose last element ensuring no conflict with previous and first
        for _, v := range []int{a[n-1], b[n-1], c[n-1]} {
            if v != p[n-2] && v != p[0] {
                p[n-1] = v
                break
            }
        }
        // output the sequence
        for i, v := range p {
            if i > 0 {
                fmt.Fprint(writer, " ")
            }
            fmt.Fprint(writer, v)
        }
        fmt.Fprintln(writer)
    }
}
