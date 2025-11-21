package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

func gcd(a, b int64) int64 {
    for b != 0 {
        a, b = b, a%b
    }
    if a < 0 {
        return -a
    }
    return a
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    fmt.Fscan(in, &T)
    for ; T > 0; T-- {
        var n int
        var a, b int64
        fmt.Fscan(in, &n, &a, &b)
        g := gcd(a, b)
        residues := make([]int64, n)
        for i := 0; i < n; i++ {
            var c int64
            fmt.Fscan(in, &c)
            if g == 0 {
                residues[i] = 0
            } else {
                residues[i] = c % g
            }
        }
        if g == 0 {
            fmt.Fprintln(out, 0)
            continue
        }
        sort.Slice(residues, func(i, j int) bool { return residues[i] < residues[j] })
        maxGap := int64(0)
        for i := 0; i+1 < n; i++ {
            gap := residues[i+1] - residues[i]
            if gap > maxGap {
                maxGap = gap
            }
        }
        wrap := g - (residues[n-1] - residues[0])
        if wrap > maxGap {
            maxGap = wrap
        }
        ans := g - maxGap
        fmt.Fprintln(out, ans)
    }
}
