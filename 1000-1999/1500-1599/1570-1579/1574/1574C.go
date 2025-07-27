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

    var n int
    if _, err := fmt.Fscan(reader, &n); err != nil {
        return
    }
    a := make([]int64, n)
    var sum int64
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &a[i])
        sum += a[i]
    }
    sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })

    var m int
    fmt.Fscan(reader, &m)
    for ; m > 0; m-- {
        var x, y int64
        fmt.Fscan(reader, &x, &y)
        idx := sort.Search(n, func(i int) bool { return a[i] >= x })
        best := int64(1<<63 - 1)
        if idx < n {
            cost := int64(0)
            rem := sum - a[idx]
            if rem < y {
                cost += y - rem
            }
            if cost < best {
                best = cost
            }
        }
        if idx > 0 {
            hero := a[idx-1]
            cost := x - hero
            rem := sum - hero
            if rem < y {
                cost += y - rem
            }
            if cost < best {
                best = cost
            }
        }
        fmt.Fprintln(writer, best)
    }
}
