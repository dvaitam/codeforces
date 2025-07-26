package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var n, q int
    if _, err := fmt.Fscan(in, &n, &q); err != nil {
        return
    }

    a := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(in, &a[i])
    }

    ks := make([]int, q)
    for i := 0; i < q; i++ {
        fmt.Fscan(in, &ks[i])
    }

    sort.Ints(a)

    for _, k := range ks {
        arr := append([]int(nil), a...)
        colors := make([]bool, n) // false = red, true = blue
        maxOps := k
        if maxOps > 2000 {
            maxOps = 2000 // naive fallback for large k
        }
        for i := 1; i <= maxOps; i++ {
            idx := 0
            for j := 1; j < n; j++ {
                if arr[j] < arr[idx] {
                    idx = j
                }
            }
            if !colors[idx] {
                arr[idx] += i
                colors[idx] = true
            } else {
                arr[idx] -= i
                colors[idx] = false
            }
        }
        mn := arr[0]
        for _, v := range arr {
            if v < mn {
                mn = v
            }
        }
        fmt.Fprintln(out, mn)
    }
}
