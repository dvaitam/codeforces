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
    tasks := make([]int64, n)
    var sum int64
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &tasks[i])
        sum += tasks[i]
    }

    // average load each server should have
    q := sum / int64(n)
    r := int(sum % int64(n))

    sort.Slice(tasks, func(i, j int) bool {
        return tasks[i] > tasks[j]
    })

    var diff int64
    for i := 0; i < n; i++ {
        target := q
        if i < r {
            target = q + 1
        }
        if tasks[i] > target {
            diff += tasks[i] - target
        } else {
            diff += target - tasks[i]
        }
    }

    fmt.Fprintln(writer, diff/2)
}
