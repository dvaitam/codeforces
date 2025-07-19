package main

import (
    "fmt"
    "sort"
)

var a []int

// gen generates all lucky numbers composed of digits 4 and 7 up to 1e9
func gen(x int) {
    if x > 1000000000 {
        return
    }
    if x != 0 {
        a = append(a, x)
    }
    gen(x*10 + 4)
    gen(x*10 + 7)
}

// get returns the fraction of overlap between [l,r] and [a,b]
func get(l, r, a, b int) float64 {
    if l < a {
        l = a
    }
    if r > b {
        r = b
    }
    if l > r {
        return 0.0
    }
    return float64(r-l+1) / float64(b-a+1)
}

func main() {
    var pl, pr, vl, vr, k int
    if _, err := fmt.Scan(&pl, &pr, &vl, &vr, &k); err != nil {
        return
    }
    gen(0)
    sort.Ints(a)
    n := len(a)
    var res float64
    for i := 0; i < n; i++ {
        j := i + k - 1
        if j >= n {
            break
        }
        // define intervals [l1,r1] around a[i] and [l2,r2] around a[j]
        l1 := 1
        if i > 0 {
            l1 = a[i-1] + 1
        }
        r1 := a[i]
        l2 := a[j]
        r2 := 1000000000
        if j < n-1 {
            r2 = a[j+1] - 1
        }
        // accumulate probability contributions
        res += get(l1, r1, pl, pr) * get(l2, r2, vl, vr)
        res += get(l2, r2, pl, pr) * get(l1, r1, vl, vr)
        // subtract double counted overlap if intervals touch
        if r1 == l2 {
            res -= get(r1, l2, pl, pr) * get(r1, l2, vl, vr)
        }
    }
    fmt.Printf("%.10f\n", res)
}
