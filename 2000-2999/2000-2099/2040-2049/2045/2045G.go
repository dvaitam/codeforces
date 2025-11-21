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

    var R, C, X int
    fmt.Fscan(in, &R, &C, &X)
    grid := make([][]int, R)
    for i := 0; i < R; i++ {
        var line string
        fmt.Fscan(in, &line)
        row := make([]int, C)
        for j := 0; j < C; j++ {
            row[j] = int(line[j] - '0')
        }
        grid[i] = row
    }

    wDiff := make([]int64, 19)
    for d := -9; d <= 9; d++ {
        val := int64(1)
        for i := 0; i < X; i++ {
            val *= int64(abs(d))
        }
        if d < 0 {
            val = -val
        }
        wDiff[d+9] = val
    }

    prefDown := make([][]int64, R)
    for i := 0; i < R; i++ {
        prefDown[i] = make([]int64, C)
    }
    for r := 1; r < R; r++ {
        row := prefDown[r]
        prev := prefDown[r-1]
        for c := 0; c < C; c++ {
            diff := grid[r-1][c] - grid[r][c]
            row[c] = prev[c] + wDiff[diff+9]
        }
    }

    prefRight := make([][]int64, R)
    for r := 0; r < R; r++ {
        prefRight[r] = make([]int64, C)
        for c := 1; c < C; c++ {
            diff := grid[r][c-1] - grid[r][c]
            prefRight[r][c] = prefRight[r][c-1] + wDiff[diff+9]
        }
    }

    invalid := false
    if R > 1 && C > 1 {
        for r := 0; r < R-1 && !invalid; r++ {
            for c := 0; c < C-1; c++ {
                tl := grid[r][c]
                tr := grid[r][c+1]
                bl := grid[r+1][c]
                br := grid[r+1][c+1]
                val := wDiff[tl-bl+9] + wDiff[bl-br+9] + wDiff[br-tr+9] + wDiff[tr-tl+9]
                if val < 0 {
                    invalid = true
                    break
                }
            }
        }
    }

    var Q int
    fmt.Fscan(in, &Q)

    vertCost := func(r1, c, r2 int) int64 {
        if r1 == r2 {
            return 0
        }
        if r1 < r2 {
            return prefDown[r2][c] - prefDown[r1][c]
        }
        return -(prefDown[r1][c] - prefDown[r2][c])
    }

    horizCost := func(r, c1, c2 int) int64 {
        if c1 == c2 {
            return 0
        }
        if c1 < c2 {
            return prefRight[r][c2] - prefRight[r][c1]
        }
        return -(prefRight[r][c1] - prefRight[r][c2])
    }

    for ; Q > 0; Q-- {
        var rs, cs, rf, cf int
        fmt.Fscan(in, &rs, &cs, &rf, &cf)
        rs--
        cs--
        rf--
        cf--
        if invalid {
            fmt.Fprintln(out, "INVALID")
            continue
        }
        if rs == rf && cs == cf {
            fmt.Fprintln(out, 0)
            continue
        }
        path1 := vertCost(rs, cs, rf) + horizCost(rf, cs, cf)
        path2 := horizCost(rs, cs, cf) + vertCost(rs, cf, rf)
        if path1 < path2 {
            fmt.Fprintln(out, path1)
        } else {
            fmt.Fprintln(out, path2)
        }
    }
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}
