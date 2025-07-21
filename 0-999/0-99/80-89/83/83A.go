package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    var n int
    if _, err := fmt.Fscan(reader, &n); err != nil {
        return
    }
    a := make([]int64, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &a[i])
    }

    var ans int64
    var run int64 = 1
    for i := 1; i < n; i++ {
        if a[i] == a[i-1] {
            run++
        } else {
            ans += run * (run + 1) / 2
            run = 1
        }
    }
    ans += run * (run + 1) / 2

    writer := bufio.NewWriter(os.Stdout)
    fmt.Fprint(writer, ans)
    writer.Flush()
}
