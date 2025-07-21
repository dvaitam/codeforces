package main

import (
    "fmt"
)

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    // Read friends' preference lists and build position lookup
    posF := make([][]int, n+1)
    for i := 1; i <= n; i++ {
        posF[i] = make([]int, n+1)
        for j := 1; j <= n; j++ {
            var x int
            fmt.Scan(&x)
            posF[i][x] = j
        }
    }
    // Read Alexander's preference list
    a := make([]int, n+1)
    for i := 1; i <= n; i++ {
        fmt.Scan(&a[i])
    }
    // Precompute for each k the top two available cards by Alexander
    b1 := make([]int, n+1)
    b2 := make([]int, n+1)
    for k := 1; k <= n; k++ {
        cnt := 0
        for t := 1; t <= n; t++ {
            card := a[t]
            if card <= k {
                cnt++
                if cnt == 1 {
                    b1[k] = card
                } else if cnt == 2 {
                    b2[k] = card
                    break
                }
            }
        }
    }
    // Determine best sending time for each friend
    res := make([]int, n+1)
    for j := 1; j <= n; j++ {
        bestRank := n + 1
        bestK := 1
        for k := 1; k <= n; k++ {
            // Determine card sent to friend j at time k
            cj := b1[k]
            if cj == j {
                cj = b2[k]
            }
            if cj == 0 {
                continue
            }
            // Friend j's preference rank for this card
            if posF[j][cj] < bestRank {
                bestRank = posF[j][cj]
                bestK = k
            }
        }
        res[j] = bestK
    }
    // Output result
    for i := 1; i <= n; i++ {
        if i > 1 {
            fmt.Print(" ")
        }
        fmt.Print(res[i])
    }
    fmt.Println()
}
