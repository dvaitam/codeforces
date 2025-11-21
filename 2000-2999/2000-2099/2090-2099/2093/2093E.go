package main

import (
    "bufio"
    "fmt"
    "os"
)

func feasible(a []int, k int, x int) bool {
    if x == 0 {
        return true
    }
    n := len(a)
    cnt := make([]int, x)
    segments := 0
    have := 0
    j := 0
    for i := 0; i < n && segments < k; i++ {
        for j < n && have < x {
            val := a[j]
            if val < x {
                if cnt[val] == 0 {
                    have++
                }
                cnt[val]++
            }
            j++
        }
        if have == x {
            segments++
            for t := 0; t < x; t++ {
                cnt[t] = 0
            }
            have = 0
        } else {
            break
        }
    }
    return segments >= k
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    fmt.Fscan(in, &t)
    for ; t > 0; t-- {
        var n, k int
        fmt.Fscan(in, &n, &k)
        a := make([]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &a[i])
        }

        low, high := 0, n
        for low < high {
            mid := (low + high + 1) >> 1
            if feasible(a, k, mid) {
                low = mid
            } else {
                high = mid - 1
            }
        }
        fmt.Fprintln(out, low)
    }
}
