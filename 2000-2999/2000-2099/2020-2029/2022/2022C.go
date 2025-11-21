package main

import (
    "bufio"
    "fmt"
    "os"
)

type cell struct{ r, c int }

const (
    windowWidth = 3
    maskSize    = 1 << (2 * windowWidth)
    negInf      = -1 << 60
)

var shapeList = [][][3]cell{
    {
        {{0, 0}, {0, 1}, {0, 2}},
        {{0, 0}, {1, 0}, {0, 1}},
        {{0, 0}, {1, 0}, {1, 1}},
        {{0, 0}, {0, 1}, {1, 1}},
    },
    {
        {{1, 0}, {1, 1}, {1, 2}},
        {{1, 0}, {0, 0}, {1, 1}},
        {{1, 0}, {0, 0}, {0, 1}},
        {{1, 0}, {1, 1}, {0, 1}},
    },
}

func solveCase(n int, top, bottom string) int {
    grid := [2][]byte{[]byte(top), []byte(bottom)}
    blocked := make([]int, n+windowWidth)
    for i := 0; i < len(blocked); i++ {
        mask := 0
        for off := 0; off < windowWidth; off++ {
            if i+off >= n {
                mask |= 3 << (2 * off)
            }
        }
        blocked[i] = mask
    }

    dp := make([]int, maskSize)
    for i := range dp {
        dp[i] = negInf
    }
    dp[0] = 0

    var dfs func(idx int, mask int, gain int, next []int)
    dfs = func(idx int, mask int, gain int, next []int) {
        mask |= blocked[idx]
        if mask&3 == 3 {
            nextMask := (mask >> 2) & (maskSize - 1)
            if gain > next[nextMask] {
                next[nextMask] = gain
            }
            return
        }
        row := 0
        if mask&1 != 0 {
            row = 1
        }
        for _, shape := range shapeList[row] {
            valid := true
            newMask := mask
            cnt := 0
            for _, cell := range shape {
                col := idx + cell.c
                if col >= n {
                    valid = false
                    break
                }
                bit := 1 << (2*cell.c + cell.r)
                if newMask&bit != 0 {
                    valid = false
                    break
                }
                newMask |= bit
                if grid[cell.r][col] == 'A' {
                    cnt++
                }
            }
            if !valid {
                continue
            }
            add := 0
            if cnt >= 2 {
                add = 1
            }
            dfs(idx, newMask, gain+add, next)
        }
    }

    for col := 0; col < n; col++ {
        next := make([]int, maskSize)
        for i := range next {
            next[i] = negInf
        }
        for mask, val := range dp {
            if val == negInf {
                continue
            }
            dfs(col, mask, val, next)
        }
        dp = next
    }

    ans := 0
    for _, val := range dp {
        if val > ans {
            ans = val
        }
    }
    return ans
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
        var s1, s2 string
        fmt.Fscan(in, &s1)
        fmt.Fscan(in, &s2)
        fmt.Fprintln(out, solveCase(n, s1, s2))
    }
}
