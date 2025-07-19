package main

import (
    "bufio"
    "fmt"
    "os"
)

var (
    n, m, x int
    a, b [8]int
    res float64
)

// cnt explores outcomes for current b[] distribution
func cnt(l, d int, s, p float64, r *float64) {
    if l == n {
        if d*2 > n {
            *r += p
        } else {
            *r += p * float64(x) / (float64(x) + s)
        }
        return
    }
    // successful kill branch
    if b[l] > 0 {
        cnt(l+1, d+1, s, p*0.1*float64(b[l]), r)
    }
    // miss branch
    if b[l] < 10 {
        cnt(l+1, d, s+float64(a[l]), p*0.1*float64(10-b[l]), r)
    }
}

// rec distributes remaining candies c across positions starting at l
func rec(l, c int) {
    if l == n {
        // evaluate this distribution
        var r float64
        cnt(0, 0, 0, 1.0, &r)
        if r > res {
            res = r
        }
        return
    }
    // try giving i candies (units) to b[l]
    for i := 0; i <= c && b[l]+i <= 10; i++ {
        b[l] += i
        rec(l+1, c-i)
        b[l] -= i
    }
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    if _, err := fmt.Fscan(reader, &n, &m, &x); err != nil {
        return
    }
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &a[i], &b[i])
        b[i] /= 10
    }
    rec(0, m)
    // print result with 8 decimal places
    fmt.Printf("%.8f\n", res)
}
