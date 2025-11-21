package main

import (
    "bufio"
    "fmt"
    "math"
    "os"
    "strings"
)

func countLE(x, r int, pref []int, n int) int {
    total := 0
    for start := 0; start <= n; start += x {
        l := start
        if l < 1 {
            l = 1
        }
        rb := start + r
        blockEnd := start + x - 1
        if rb > blockEnd {
            rb = blockEnd
        }
        if rb > n {
            rb = n
        }
        if rb >= 1 && l <= rb {
            total += pref[rb] - pref[l-1]
        }
    }
    return total
}

func solveLarge(x int, pref []int, n int, target int) int {
    lo, hi := 0, x-1
    ans := 0
    for lo <= hi {
        mid := (lo + hi) / 2
        if countLE(x, mid, pref, n) >= target {
            ans = mid
            hi = mid - 1
        } else {
            lo = mid + 1
        }
    }
    return ans
}

func main() {
    in := bufio.NewReader(os.Stdin)
    var T int
    fmt.Fscan(in, &T)
    var sb strings.Builder
    for ; T > 0; T-- {
        var n, q int
        fmt.Fscan(in, &n, &q)
        a := make([]int, n)
        freq := make([]int, n+1)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &a[i])
            freq[a[i]]++
        }
        pref := make([]int, n+1)
        for i := 1; i <= n; i++ {
            pref[i] = pref[i-1] + freq[i]
        }
        target := n/2 + 1
        B := int(math.Sqrt(float64(n))) + 1
        ansSmall := make([]int, B+1)
        for x := 1; x <= B; x++ {
            counts := make([]int, x)
            for _, v := range a {
                counts[v%x]++
            }
            sum := 0
            res := 0
            for r := 0; r < x; r++ {
                sum += counts[r]
                if sum >= target {
                    res = r
                    break
                }
            }
            ansSmall[x] = res
        }
        cache := make(map[int]int)
        for i := 0; i < q; i++ {
            var x int
            fmt.Fscan(in, &x)
            var ans int
            if x <= B {
                ans = ansSmall[x]
            } else {
                if val, ok := cache[x]; ok {
                    ans = val
                } else {
                    val := solveLarge(x, pref, n, target)
                    cache[x] = val
                    ans = val
                }
            }
            sb.WriteString(fmt.Sprintf("%d ", ans))
        }
        sb.WriteString("\n")
    }
    fmt.Print(sb.String())
}
