package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    if _, err := fmt.Fscan(in, &T); err != nil {
        return
    }
    for ; T > 0; T-- {
        var n, m int
        fmt.Fscan(in, &n, &m)
        // read columns as slices
        cols := make([][]int, m)
        for j := 0; j < m; j++ {
            cols[j] = make([]int, n)
        }
        for i := 0; i < n; i++ {
            for j := 0; j < m; j++ {
                fmt.Fscan(in, &cols[j][i])
            }
        }
        var total int64
        for j := 0; j < m; j++ {
            col := cols[j]
            sort.Ints(col)
            var prefix int64
            for i := 0; i < n; i++ {
                total += int64(col[i])*int64(i) - prefix
                prefix += int64(col[i])
            }
        }
        fmt.Fprintln(out, total)
    }
}
