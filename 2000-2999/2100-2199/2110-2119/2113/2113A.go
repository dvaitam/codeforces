package main

import (
    "bufio"
    "fmt"
    "os"
)

func max(a, b int64) int64 {
    if a > b {
        return a
    }
    return b
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    fmt.Fscan(in, &T)
    for ; T > 0; T-- {
        var k, a, b, x, y int64
        fmt.Fscan(in, &k, &a, &b, &x, &y)

        if a > b {
            a, b = b, a
            x, y = y, x
        }

        if k < a {
            fmt.Fprintln(out, 0)
            continue
        }

        var best int64
        var limit int64
        if x < y {
            limit = min((k-a)/x+1, int64(200000))
        } else {
            limit = min((k-b)/y+1, int64(200000))
        }

        for i := int64(0); i <= limit; i++ {
            temp := k - i*x
            if temp < 0 {
                break
            }
            if temp < a {
                break
            }
            var cnt int64
            cnt = i
            temp2 := temp
            for temp2 >= b {
                temp2 -= y
                cnt++
            }
            for temp2 >= a {
                temp2 -= x
                cnt++
            }
            if cnt > best {
                best = cnt
            }
        }
        fmt.Fprintln(out, best)
    }
}

func min(a, b int64) int64 {
    if a < b {
        return a
    }
    return b
}
