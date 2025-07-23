package main

import (
    "bufio"
    "fmt"
    "os"
)

func exgcd(a, b int64) (g, x, y int64) {
    if b == 0 {
        return a, 1, 0
    }
    g, x1, y1 := exgcd(b, a%b)
    x = y1
    y = x1 - y1*(a/b)
    return
}

// solve t ≡ a1 (mod m1), t ≡ a2 (mod m2)
func crt(a1, m1, a2, m2 int64) (t int64, ok bool) {
    g, x, _ := exgcd(m1, m2)
    if (a2-a1)%g != 0 {
        return 0, false
    }
    lcm := m1 / g * m2
    mul := ((a2 - a1) / g * x) % (m2 / g)
    if mul < 0 {
        mul += m2 / g
    }
    t = (a1 + mul*m1) % lcm
    if t < 0 {
        t += lcm
    }
    return t, true
}

func mod(a, m int64) int64 {
    a %= m
    if a < 0 {
        a += m
    }
    return a
}

func main() {
    in := bufio.NewReader(os.Stdin)
    var n, m, x, y int64
    var vx, vy int64
    if _, err := fmt.Fscan(in, &n, &m, &x, &y, &vx, &vy); err != nil {
        return
    }

    // handle cases where velocity is zero on an axis
    if vx == 0 {
        if x != 0 && x != n {
            fmt.Println(-1)
            return
        }
        if vy == 1 {
            fmt.Printf("%d %d\n", x, m)
        } else {
            fmt.Printf("%d %d\n", x, 0)
        }
        return
    }
    if vy == 0 {
        if y != 0 && y != m {
            fmt.Println(-1)
            return
        }
        if vx == 1 {
            fmt.Printf("%d %d\n", n, y)
        } else {
            fmt.Printf("%d %d\n", 0, y)
        }
        return
    }

    bestT := int64(-1)
    ansX, ansY := int64(-1), int64(-1)

    rxOptions := []int64{0, n}
    ryOptions := []int64{0, m}
    for _, rx := range rxOptions {
        for _, ry := range ryOptions {
            // times modulo 2n and 2m
            t1 := int64(0)
            if vx == 1 {
                t1 = mod(rx-x, 2*n)
            } else { // vx == -1
                t1 = mod(x-rx, 2*n)
            }
            t2 := int64(0)
            if vy == 1 {
                t2 = mod(ry-y, 2*m)
            } else { // vy == -1
                t2 = mod(y-ry, 2*m)
            }
            t, ok := crt(t1, 2*n, t2, 2*m)
            if !ok {
                continue
            }
            if bestT == -1 || t < bestT {
                bestT = t
                ansX = rx
                ansY = ry
            }
        }
    }

    if bestT == -1 {
        fmt.Println(-1)
    } else {
        fmt.Printf("%d %d\n", ansX, ansY)
    }
}
