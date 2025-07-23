package main

import (
    "bufio"
    "fmt"
    "os"
)

// find returns the first time tFirst when sequence reaches a, and the cycle length after that.
// If a is never reached, tFirst=-1. If a is reached only once, cycle=-1.
func find(h0, a, x, y, m int64) (tFirst, cycle int64) {
    tFirst = -1
    cycle = -1
    h := h0
    // simulate up to 2*m steps to find first occurrence and next for cycle
    for t := int64(0); t < 2*m; t++ {
        if h == a {
            if tFirst == -1 {
                tFirst = t
            } else {
                cycle = t - tFirst
                break
            }
        }
        h = (x*h + y) % m
    }
    return
}

// extgcd returns g=gcd(a,b) and x,y such that a*x+b*y=g.
func extgcd(a, b int64) (g, x, y int64) {
    if b == 0 {
        return a, 1, 0
    }
    g, x1, y1 := extgcd(b, a%b)
    return g, y1, x1 - (a/b)*y1
}

// crt solves x ≡ a1 mod n1, x ≡ a2 mod n2. Returns smallest x>=0 mod lcm and lcm, ok if solvable.
func crt(a1, n1, a2, n2 int64) (x, lcm int64, ok bool) {
    g, m1, _ := extgcd(n1, n2)
    if (a2-a1)%g != 0 {
        return 0, 0, false
    }
    lcm = n1 / g * n2
    // compute k = ((a2-a1)/g) * inv(n1/g mod n2/g) mod (n2/g)
    mod := n2 / g
    k := (a2 - a1) / g % mod
    inv := m1 % mod
    if inv < 0 {
        inv += mod
    }
    k = k * inv % mod
    if k < 0 {
        k += mod
    }
    x = a1 + n1*k
    x %= lcm
    if x < 0 {
        x += lcm
    }
    return x, lcm, true
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    var m int64
    fmt.Fscan(reader, &m)
    var h1, a1, x1, y1 int64
    fmt.Fscan(reader, &h1, &a1)
    fmt.Fscan(reader, &x1, &y1)
    var h2, a2, x2, y2 int64
    fmt.Fscan(reader, &h2, &a2)
    fmt.Fscan(reader, &x2, &y2)

    t1, c1 := find(h1, a1, x1, y1, m)
    t2, c2 := find(h2, a2, x2, y2, m)
    if t1 == -1 || t2 == -1 {
        fmt.Println(-1)
        return
    }
    // both reach only once
    if c1 == -1 && c2 == -1 {
        if t1 == t2 {
            fmt.Println(t1)
        } else {
            fmt.Println(-1)
        }
        return
    }
    // one unique, one cyclic
    if c1 == -1 {
        if t1 >= t2 && c2 > 0 && (t1-t2)%c2 == 0 {
            fmt.Println(t1)
        } else {
            fmt.Println(-1)
        }
        return
    }
    if c2 == -1 {
        if t2 >= t1 && c1 > 0 && (t2-t1)%c1 == 0 {
            fmt.Println(t2)
        } else {
            fmt.Println(-1)
        }
        return
    }
    // both cyclic
    x0, lcm, ok := crt(t1, c1, t2, c2)
    if !ok {
        fmt.Println(-1)
        return
    }
    ans := x0
    mx := t1
    if t2 > mx {
        mx = t2
    }
    if ans < mx {
        k := (mx - ans + lcm - 1) / lcm
        ans += k * lcm
    }
    fmt.Println(ans)
}
