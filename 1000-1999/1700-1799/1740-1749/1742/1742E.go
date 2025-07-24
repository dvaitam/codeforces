package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var t int
    if _, err := fmt.Fscan(reader, &t); err != nil {
        return
    }
    for ; t > 0; t-- {
        var n, q int
        fmt.Fscan(reader, &n, &q)
        a := make([]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(reader, &a[i])
        }
        prefMax := make([]int, n)
        prefSum := make([]int64, n)
        curMax := 0
        var curSum int64
        for i, v := range a {
            if v > curMax {
                curMax = v
            }
            curSum += int64(v)
            prefMax[i] = curMax
            prefSum[i] = curSum
        }
        for i := 0; i < q; i++ {
            var k int
            fmt.Fscan(reader, &k)
            idx := sort.Search(n, func(j int) bool { return prefMax[j] > k })
            var ans int64
            if idx > 0 {
                ans = prefSum[idx-1]
            }
            if i > 0 {
                writer.WriteByte(' ')
            }
            fmt.Fprint(writer, ans)
        }
        writer.WriteByte('\n')
    }
}
