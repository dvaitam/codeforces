package main

import (
    "bufio"
    "fmt"
    "os"
)

const SQRTN = 1000

var dp [SQRTN][SQRTN]int64
var binom [SQRTN][SQRTN]int64

func h(x int) int64 {
    return binom[x+3][3]
}

func g(num int, siz int, mod int64) int64 {
    if num == 0 {
        return 1
    }
    if siz == SQRTN {
        return 0
    }
    if dp[num][siz] != -1 {
        return dp[num][siz]
    }
    var res int64
    for i := 0; i*siz <= num; i++ {
        res = (res + g(num-i*siz, siz+1, mod)*h(i)) % mod
    }
    dp[num][siz] = res
    return res
}

func f(x, y, n int64, mod int64) int64 {
    num := int(x*y - n)
    return g(num, 1, mod)
}

func solve(w *bufio.Writer, r *bufio.Reader, m int64, u int) {
    var n int64
    fmt.Fscan(r, &n)
    low, up := int64(0), int64(1500)
    for up-low > 1 {
        mid := (up + low) / 2
        a := mid / 2
        b := mid - a
        if a*b >= n {
            up = mid
        } else {
            low = mid
        }
    }
    a := up / 2
    b := up - a
    if u == 1 {
        fmt.Fprintf(w, "%d %d\n", a, b)
        rem := n
        for i := int64(0); i < a; i++ {
            for j := int64(0); j < b; j++ {
                if rem > 0 {
                    w.WriteByte('#')
                    rem--
                } else {
                    w.WriteByte('.')
                }
            }
            w.WriteByte('\n')
        }
    } else {
        // u == 2
        fmt.Fprintf(w, "%d ", (a+b)*2)
        var ans int64
        for i := a; (up-i) >= 1 && i*(up-i) >= n; i++ {
            ans = (ans + f(i, up-i, n, m)) % m
        }
        for i := a - 1; i >= 1 && i*(up-i) >= n; i-- {
            ans = (ans + f(i, up-i, n, m)) % m
        }
        fmt.Fprintf(w, "%d\n", ans)
    }
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()
    var t int
    var u int
    fmt.Fscan(reader, &t, &u)
    var m int64
    if u == 2 {
        fmt.Fscan(reader, &m)
    }
    // initialize dp with -1
    for i := 0; i < SQRTN; i++ {
        for j := 0; j < SQRTN; j++ {
            dp[i][j] = -1
        }
    }
    // initialize binom if needed
    if u == 2 {
        for i := 0; i < SQRTN; i++ {
            for j := 0; j <= i; j++ {
                if j == 0 || j == i {
                    binom[i][j] = 1
                } else {
                    binom[i][j] = (binom[i-1][j] + binom[i-1][j-1]) % m
                }
            }
        }
    }
    for i := 0; i < t; i++ {
        solve(writer, reader, m, u)
    }
}
