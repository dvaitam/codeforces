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

    var n int
    if _, err := fmt.Fscan(in, &n); err != nil {
        return
    }

    values := make([]int, 0)
    down := make(map[int]int)

    for i := 0; i < n; i++ {
        var k int
        fmt.Fscan(in, &k)
        tower := make([]int, k)
        for j := 0; j < k; j++ {
            fmt.Fscan(in, &tower[j])
            values = append(values, tower[j])
            if j > 0 {
                down[tower[j-1]] = tower[j]
            }
        }
    }

    sort.Ints(values)

    breaks := 0
    for i := 0; i+1 < len(values); i++ {
        if next, ok := down[values[i]]; !ok || next != values[i+1] {
            breaks++
        }
    }

    segments := breaks + 1
    splits := segments - n
    combines := segments - 1
    fmt.Fprintf(out, "%d %d\n", splits, combines)
}
