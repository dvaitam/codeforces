package main

import (
    "bufio"
    "fmt"
    "os"
)

const inf int64 = 1<<62 - 1

func minCost(num int, costs []int64, banSame, banPlus, banMinus []bool) (int64, bool) {
    dp0, dp1 := int64(0), costs[0]
    for i := 0; i < num-1; i++ {
        n0, n1 := inf, inf
        if !banSame[i] && dp0 < inf {
            if dp0 < n0 {
                n0 = dp0
            }
        }
        if !banMinus[i] && dp1 < inf {
            if dp1 < n0 {
                n0 = dp1
            }
        }
        if !banPlus[i] && dp0 < inf {
            if dp0 < n1 {
                n1 = dp0
            }
        }
        if !banSame[i] && dp1 < inf {
            if dp1 < n1 {
                n1 = dp1
            }
        }
        if n1 < inf {
            n1 += costs[i+1]
        }
        dp0, dp1 = n0, n1
    }
    res := dp0
    if dp1 < res {
        res = dp1
    }
    if res >= inf {
        return 0, false
    }
    return res, true
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    fmt.Fscan(in, &T)
    for ; T > 0; T-- {
        var n int
        fmt.Fscan(in, &n)
        h := make([][]int, n)
        for i := 0; i < n; i++ {
            h[i] = make([]int, n)
            for j := 0; j < n; j++ {
                fmt.Fscan(in, &h[i][j])
            }
        }
        rowCost := make([]int64, n)
        colCost := make([]int64, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &rowCost[i])
        }
        for j := 0; j < n; j++ {
            fmt.Fscan(in, &colCost[j])
        }

        banSameRow := make([]bool, n-1)
        banPlusRow := make([]bool, n-1)
        banMinusRow := make([]bool, n-1)
        for i := 0; i < n-1; i++ {
            for j := 0; j < n; j++ {
                diff := h[i][j] - h[i+1][j]
                if diff == 0 {
                    banSameRow[i] = true
                } else if diff == 1 {
                    banPlusRow[i] = true
                } else if diff == -1 {
                    banMinusRow[i] = true
                }
            }
        }

        banSameCol := make([]bool, n-1)
        banPlusCol := make([]bool, n-1)
        banMinusCol := make([]bool, n-1)
        for j := 0; j < n-1; j++ {
            for i := 0; i < n; i++ {
                diff := h[i][j] - h[i][j+1]
                if diff == 0 {
                    banSameCol[j] = true
                } else if diff == 1 {
                    banPlusCol[j] = true
                } else if diff == -1 {
                    banMinusCol[j] = true
                }
            }
        }

        rowAns, ok1 := minCost(n, rowCost, banSameRow, banPlusRow, banMinusRow)
        colAns, ok2 := minCost(n, colCost, banSameCol, banPlusCol, banMinusCol)
        if !ok1 || !ok2 {
            fmt.Fprintln(out, -1)
        } else {
            fmt.Fprintln(out, rowAns+colAns)
        }
    }
}
