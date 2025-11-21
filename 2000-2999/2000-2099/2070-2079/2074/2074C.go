package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    fmt.Fscan(in, &T)
    for ; T > 0; T-- {
        var x int64
        fmt.Fscan(in, &x)
        if isPowerOfTwo(x) || isAllOnes(x) {
            fmt.Fprintln(out, -1)
            continue
        }
        msb := bitLen(x) - 1
        var s, q int
        foundS := false
        foundQ := false
        for i := 0; i < msb; i++ {
            if ((x>>int64(i))&1) == 1 {
                s = i
                foundS = true
                break
            }
        }
        for i := 0; i < msb; i++ {
            if ((x>>int64(i))&1) == 0 {
                q = i
                foundQ = true
                break
            }
        }
        if !foundS || !foundQ || s == q {
            fmt.Fprintln(out, -1)
            continue
        }
        y := (int64(1) << s) | (int64(1) << q)
        fmt.Fprintln(out, y)
    }
}

func isPowerOfTwo(x int64) bool {
    return x > 0 && (x&(x-1)) == 0
}

func isAllOnes(x int64) bool {
    return ((x + 1) & x) == 0
}

func bitLen(x int64) int {
    l := 0
    for x > 0 {
        l++
        x >>= 1
    }
    return l
}
