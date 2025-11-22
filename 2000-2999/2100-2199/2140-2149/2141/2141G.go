package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

func countRunRectangles(h []int) int64 {
    k := len(h)
    if k == 0 {
        return 0
    }
    prev := make([]int, k)
    stack := make([]int, 0, k)
    for i := 0; i < k; i++ {
        for len(stack) > 0 && h[stack[len(stack)-1]] >= h[i] {
            stack = stack[:len(stack)-1]
        }
        if len(stack) == 0 {
            prev[i] = -1
        } else {
            prev[i] = stack[len(stack)-1]
        }
        stack = append(stack, i)
    }
    next := make([]int, k)
    stack = stack[:0]
    for i := k - 1; i >= 0; i-- {
        for len(stack) > 0 && h[stack[len(stack)-1]] > h[i] {
            stack = stack[:len(stack)-1]
        }
        if len(stack) == 0 {
            next[i] = k
        } else {
            next[i] = stack[len(stack)-1]
        }
        stack = append(stack, i)
    }
    var totalGe2WidthGe1 int64
    var width1Sum int64
    for i := 0; i < k; i++ {
        if h[i] >= 2 {
            val := int64(h[i] - 1)
            totalGe2WidthGe1 += val * int64(i-prev[i]) * int64(next[i]-i)
            width1Sum += val
        }
    }
    return totalGe2WidthGe1 - width1Sum // width >= 2, height >= 2
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    if _, err := fmt.Fscan(in, &T); err != nil {
        return
    }

    for ; T > 0; T-- {
        var n int
        fmt.Fscan(in, &n)

        rowMap := make(map[int64][]int64)
        colMap := make(map[int64][]int64)
        xsAll := make([]int64, 0, n)

        for i := 0; i < n; i++ {
            var x, y int64
            fmt.Fscan(in, &x, &y)
            rowMap[y] = append(rowMap[y], x)
            colMap[x] = append(colMap[x], y)
        }

        // collect unique x values
        for x := range colMap {
            xsAll = append(xsAll, x)
        }
        sort.Slice(xsAll, func(i, j int) bool { return xsAll[i] < xsAll[j] })

        var ans int64

        // horizontal segments (same row)
        for _, xs := range rowMap {
            sort.Slice(xs, func(i, j int) bool { return xs[i] < xs[j] })
            cnt := 1
            for i := 1; i < len(xs); i++ {
                if xs[i] == xs[i-1]+1 {
                    cnt++
                } else {
                    ans += int64(cnt) * int64(cnt-1)
                    cnt = 1
                }
            }
            ans += int64(cnt) * int64(cnt-1)
        }

        // vertical segments (same column)
        for _, ys := range colMap {
            sort.Slice(ys, func(i, j int) bool { return ys[i] < ys[j] })
            cnt := 1
            for i := 1; i < len(ys); i++ {
                if ys[i] == ys[i-1]+1 {
                    cnt++
                } else {
                    ans += int64(cnt) * int64(cnt-1)
                    cnt = 1
                }
            }
            ans += int64(cnt) * int64(cnt-1)
        }

        // rectangles with height >= 2 and width >= 2
        var rectCount int64

        // process blocks of consecutive x values
        i := 0
        for i < len(xsAll) {
            j := i
            for j+1 < len(xsAll) && xsAll[j+1] == xsAll[j]+1 {
                j++
            }
            blockXs := xsAll[i : j+1]
            w := len(blockXs)
            // map from y to list of column indices in this block
            rowCols := make(map[int64][]int)
            for idx, x := range blockXs {
                ys := colMap[x]
                for _, y := range ys {
                    rowCols[y] = append(rowCols[y], idx)
                }
            }
            // sorted y keys
            yKeys := make([]int64, 0, len(rowCols))
            for y := range rowCols {
                yKeys = append(yKeys, y)
            }
            sort.Slice(yKeys, func(a, b int) bool { return yKeys[a] < yKeys[b] })

            heights := make([]int, w)
            lastRow := make([]int64, w)
            for k := range lastRow {
                lastRow[k] = int64(1<<62) * -1 // effectively very small
            }

            for _, y := range yKeys {
                cols := rowCols[y]
                sort.Ints(cols)
                for _, c := range cols {
                    if lastRow[c] == y-1 {
                        heights[c]++
                    } else {
                        heights[c] = 1
                    }
                    lastRow[c] = y
                }
                // process runs of consecutive columns present in this row
                start := 0
                for start < len(cols) {
                    end := start
                    for end+1 < len(cols) && cols[end+1] == cols[end]+1 {
                        end++
                    }
                    run := cols[start : end+1]
                    hRun := make([]int, len(run))
                    for idx, c := range run {
                        hRun[idx] = heights[c]
                    }
                    rectCount += countRunRectangles(hRun)
                    start = end + 1
                }
            }

            i = j + 1
        }

        ans += rectCount * 4
        fmt.Fprintln(out, ans)
    }
}

