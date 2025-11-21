package main

import (
    "bufio"
    "fmt"
    "os"
)

type factor struct {
    p int
    e int
}

func gcd(a, b int) int {
    for b != 0 {
        a, b = b, a%b
    }
    if a < 0 {
        return -a
    }
    return a
}

func buildSPF(limit int) []int {
    spf := make([]int, limit+1)
    for i := 2; i <= limit; i++ {
        if spf[i] == 0 {
            spf[i] = i
            if i*i <= limit {
                for j := i * i; j <= limit; j += i {
                    if spf[j] == 0 {
                        spf[j] = i
                    }
                }
            }
        }
    }
    spf[1] = 1
    return spf
}

func factorize(x int, spf []int) []factor {
    if x == 1 {
        return nil
    }
    res := make([]factor, 0)
    for x > 1 {
        p := spf[x]
        cnt := 0
        for x%p == 0 {
            x /= p
            cnt++
        }
        res = append(res, factor{p, cnt})
    }
    return res
}

func powMod(base, exp, mod int) int {
    if mod == 1 {
        return 0
    }
    result := 1
    b := base % mod
    for exp > 0 {
        if exp&1 == 1 {
            result = int((int64(result) * int64(b)) % int64(mod))
        }
        b = int((int64(b) * int64(b)) % int64(mod))
        exp >>= 1
    }
    return result
}

func eIn(num, prime int) int {
    cnt := 0
    for num%prime == 0 {
        num /= prime
        cnt++
    }
    return cnt
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    fmt.Fscan(in, &t)
    tests := make([][2]int, t)
    maxVal := 0
    for i := 0; i < t; i++ {
        fmt.Fscan(in, &tests[i][0], &tests[i][1])
        if tests[i][0] > maxVal {
            maxVal = tests[i][0]
        }
        if tests[i][1] > maxVal {
            maxVal = tests[i][1]
        }
    }

    spf := buildSPF(maxVal)

    for _, test := range tests {
        b, n := test[0], test[1]
        bestType := 0
        bestK := 0

        factorsN := factorize(n, spf)

        // Type 1
        k1 := 0
        possible := true
        for _, f := range factorsN {
            eb := eIn(b, f.p)
            if eb == 0 {
                possible = false
                break
            }
            need := (f.e + eb - 1) / eb
            if need > k1 {
                k1 = need
            }
        }
        if possible && k1 > 0 {
            bestType = 1
            bestK = k1
        }

        // Type 2 and 3 require gcd = 1
        if gcd(b, n) == 1 {
            // compute order
            phi := n
            for _, f := range factorsN {
                phi = phi / f.p * (f.p - 1)
            }

            order := phi
            factorsPhi := factorize(order, spf)
            for _, f := range factorsPhi {
                for order%f.p == 0 {
                    if powMod(b, order/f.p, n) == 1 {
                        order /= f.p
                    } else {
                        break
                    }
                }
            }

            // Type 2
            if order > 0 && (bestType == 0 || order < bestK) {
                bestType = 2
                bestK = order
            }

            // Type 3
            if order%2 == 0 {
                if powMod(b, order/2, n) == (n-1)%n {
                    k3 := order / 2
                    if bestType == 0 || k3 < bestK || (k3 == bestK && 3 < bestType) {
                        bestType = 3
                        bestK = k3
                    }
                }
            }
        }

        if bestType == 0 {
            fmt.Fprintln(out, 0)
        } else {
            fmt.Fprintf(out, "%d %d\n", bestType, bestK)
        }
    }
}
