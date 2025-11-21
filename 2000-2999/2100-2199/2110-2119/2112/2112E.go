package main

import (
    "bufio"
    "fmt"
    "os"
)

const maxM = 500000
const INF int64 = 1<<60

func main() {
    divs := make([][]int, maxM+1)
    for d := 3; d <= maxM; d++ {
        for mult := d; mult <= maxM; mult += d {
            divs[mult] = append(divs[mult], d)
        }
    }

    dp := make([]int64, maxM+1)
    h := make([]int64, maxM+1)
    for i := 0; i <= maxM; i++ {
        dp[i] = INF
        h[i] = INF
    }
    dp[1] = 1
    h[1] = 0

    for m := 2; m <= maxM; m++ {
        best := INF
        for _, d := range divs[m] {
            child := d - 2
            if child <= 0 || child > maxM {
                continue
            }
            if dp[child] == INF {
                continue
            }
            rem := m / d
            if h[rem] == INF {
                continue
            }
            cand := dp[child] + h[rem]
            if cand < best {
                best = cand
            }
        }
        if best < INF {
            h[m] = best
            dp[m] = best + 1
        }
    }

    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    fmt.Fscan(in, &t)
    for ; t > 0; t-- {
        var m int
        fmt.Fscan(in, &m)
        if m <= maxM && dp[m] < INF {
            fmt.Fprintln(out, dp[m])
        } else {
            fmt.Fprintln(out, -1)
        }
    }
}
