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
        var n int
        fmt.Fscan(reader, &n)
        arr := make([]int64, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(reader, &arr[i])
        }
        a1 := arr[0]
        if n > 1 {
            rest := arr[1:]
            sort.Slice(rest, func(i, j int) bool { return rest[i] < rest[j] })
            for _, h := range rest {
                if h > a1 {
                    diff := h - a1
                    a1 += (diff + 1) / 2
                }
            }
        }
        fmt.Fprintln(writer, a1)
    }
}
