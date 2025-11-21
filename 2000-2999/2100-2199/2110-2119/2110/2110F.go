package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

const K = 2000

func maxBeauty(arr []int) int {
    n := len(arr)
    limit := n
    if limit > K {
        limit = K
    }
    top := make([]int, limit)
    copy(top, arr[n-limit:])
    best := 0
    for _, x := range top {
        for _, y := range arr {
            if x == 0 || y == 0 {
                continue
            }
            val := x%y + y%x
            if val > best {
                best = val
            }
        }
    }
    return best
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
        arr := make([]int, n)
        for i := range arr {
            fmt.Fscan(in, &arr[i])
        }
        sort.Ints(arr)

        res := make([]int, n)
        for i := 0; i < n; i++ {
            currArr := arr[:i+1]
            res[i] = maxBeauty(currArr)
        }

        for i := 0; i < n; i++ {
            if i > 0 {
                fmt.Fprint(out, " ")
            }
            fmt.Fprint(out, res[i])
        }
        fmt.Fprintln(out)
    }
}
