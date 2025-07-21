package main

import (
    "bufio"
    "fmt"
    "os"
)

var n int
var a []int
var P [][][2]int
var lastUse []int
var M int

// dfs tries to assign pairs for generating a[k]..a[n] with register limit M
func dfs(k int) bool {
    if k > n {
        return true
    }
    for _, pair := range P[k] {
        i, j := pair[0], pair[1]
        oldI, oldJ := lastUse[i], lastUse[j]
        lastUse[i], lastUse[j] = k, k
        // compute live count at step k
        cnt := 0
        for t := 1; t < k; t++ {
            if lastUse[t] >= k {
                cnt++
            }
        }
        if cnt <= M {
            if dfs(k + 1) {
                return true
            }
        }
        lastUse[i], lastUse[j] = oldI, oldJ
    }
    return false
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    if _, err := fmt.Fscan(reader, &n); err != nil {
        return
    }
    a = make([]int, n+1)
    for i := 1; i <= n; i++ {
        fmt.Fscan(reader, &a[i])
    }
    // precompute possible summand pairs for each k
    P = make([][][2]int, n+1)
    for k := 2; k <= n; k++ {
        for i := 1; i < k; i++ {
            for j := i; j < k; j++ {
                if a[i]+a[j] == a[k] {
                    P[k] = append(P[k], [2]int{i, j})
                }
            }
        }
        if len(P[k]) == 0 {
            fmt.Println(-1)
            return
        }
    }
    // try minimal m from 1 to n
    for m := 1; m <= n; m++ {
        M = m
        lastUse = make([]int, n+1)
        if dfs(2) {
            fmt.Println(m)
            return
        }
    }
    fmt.Println(-1)
}
