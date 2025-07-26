package main

import (
    "bufio"
    "fmt"
    "os"
)

func solve(r *bufio.Reader, w *bufio.Writer) {
    var n int
    fmt.Fscan(r, &n)
    freq := make([]int, 101)
    maxVal := 0
    for i := 0; i < n; i++ {
        var x int
        fmt.Fscan(r, &x)
        if x > maxVal {
            maxVal = x
        }
        if x >= len(freq) {
            tmp := make([]int, x+1)
            copy(tmp, freq)
            freq = tmp
        }
        freq[x]++
    }
    ok := true
    for i := 1; i <= maxVal; i++ {
        if freq[i] > freq[i-1] {
            ok = false
            break
        }
    }
    if ok {
        fmt.Fprintln(w, "YES")
    } else {
        fmt.Fprintln(w, "NO")
    }
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var t int
    if _, err := fmt.Fscan(reader, &t); err != nil {
        return
    }
    for ; t > 0; t-- {
        solve(reader, writer)
    }
}
