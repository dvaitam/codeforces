package main

import (
    "bufio"
    "fmt"
    "math/bits"
    "os"
)

type signed128 struct {
    sign int
    hi   uint64
    lo   uint64
}

func absInt64(x int64) uint64 {
    mask := x >> 63
    return uint64((x ^ mask) - mask)
}

func mulSigned(a, b int64) signed128 {
    if a == 0 || b == 0 {
        return signed128{sign: 0}
    }
    sign := 1
    if (a < 0) != (b < 0) {
        sign = -1
    }
    ua := absInt64(a)
    ub := absInt64(b)
    hi, lo := bits.Mul64(ua, ub)
    return signed128{sign: sign, hi: hi, lo: lo}
}

func cmpFrac(num1, den1, num2, den2 int64) int {
    // compare num1/den1 ? num2/den2
    prod1 := mulSigned(num1, den2)
    prod2 := mulSigned(num2, den1)
    if prod1.sign != prod2.sign {
        if prod1.sign < prod2.sign {
            return -1
        }
        return 1
    }
    if prod1.sign == 0 {
        return 0
    }
    if prod1.sign > 0 {
        if prod1.hi != prod2.hi {
            if prod1.hi < prod2.hi {
                return -1
            }
            return 1
        }
        if prod1.lo != prod2.lo {
            if prod1.lo < prod2.lo {
                return -1
            }
            return 1
        }
        return 0
    }
    // both negative: larger magnitude means smaller value
    if prod1.hi != prod2.hi {
        if prod1.hi > prod2.hi {
            return -1
        }
        return 1
    }
    if prod1.lo != prod2.lo {
        if prod1.lo > prod2.lo {
            return -1
        }
        return 1
    }
    return 0
}

func floorDiv(a, b int64) int64 {
    if a >= 0 {
        return a / b
    }
    return -((-a + b - 1) / b)
}

func ceilDiv(a, b int64) int64 {
    return -floorDiv(-a, b)
}

func feasible(pref []int64, n int, d int64) bool {
    total := pref[n]
    ubNum := pref[1]
    ubDen := int64(1)
    for k := 2; k <= n; k++ {
        num := pref[k]
        den := int64(k)
        if cmpFrac(num, den, ubNum, ubDen) < 0 {
            ubNum = num
            ubDen = den
        }
    }
    lbNum := total - int64(n)*d
    lbDen := int64(n)
    for k := 1; k < n; k++ {
        num := total - pref[k] - d*int64(n-k)
        den := int64(n - k)
        if cmpFrac(num, den, lbNum, lbDen) > 0 {
            lbNum = num
            lbDen = den
        }
    }
    low := ceilDiv(lbNum, lbDen)
    high := floorDiv(ubNum, ubDen)
    return low <= high
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    fmt.Fscan(in, &t)
    for ; t > 0; t-- {
        var n int
        fmt.Fscan(in, &n)
        a := make([]int64, n)
        var mn, mx int64
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &a[i])
            if i == 0 || a[i] < mn {
                mn = a[i]
            }
            if i == 0 || a[i] > mx {
                mx = a[i]
            }
        }
        pref := make([]int64, n+1)
        for i := 0; i < n; i++ {
            pref[i+1] = pref[i] + a[i]
        }
        lo, hi := int64(0), mx-mn
        for lo < hi {
            mid := (lo + hi) / 2
            if feasible(pref, n, mid) {
                hi = mid
            } else {
                lo = mid + 1
            }
        }
        fmt.Fprintln(out, lo)
    }
}
