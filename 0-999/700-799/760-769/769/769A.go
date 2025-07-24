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
    years := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &years[i])
    }
    sort.Ints(years)
    median := years[n/2]
    fmt.Fprintln(writer, median)
}
