package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var n int
    if _, err := fmt.Fscan(in, &n); err != nil {
        return
    }

    a := make([]int, n)
    freq := make(map[int]int)
    total := 0
    for i := 0; i < n; i++ {
        fmt.Fscan(in, &a[i])
        total += a[i]
        freq[a[i]]++
    }

    result := make([]int, 0)
    for i := 0; i < n; i++ {
        val := a[i]
        freq[val]--
        remaining := total - val
        if remaining%2 == 0 {
            target := remaining / 2
            if freq[target] > 0 {
                result = append(result, i+1)
            }
        }
        freq[val]++
    }

    fmt.Fprintln(out, len(result))
    if len(result) > 0 {
        for i, idx := range result {
            if i > 0 {
                fmt.Fprint(out, " ")
            }
            fmt.Fprint(out, idx)
        }
        fmt.Fprintln(out)
    }
}
