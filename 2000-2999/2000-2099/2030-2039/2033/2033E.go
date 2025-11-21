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

    var t int
    fmt.Fscan(in, &t)
    for ; t > 0; t-- {
        var n int
        fmt.Fscan(in, &n)
        p := make([]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &p[i])
            p[i]--
        }
        visited := make([]bool, n)
        var ans int64
        for i := 0; i < n; i++ {
            if visited[i] {
                continue
            }
            length := 0
            cur := i
            for !visited[cur] {
                visited[cur] = true
                cur = p[cur]
                length++
            }
            ans += int64((length - 1) / 2)
        }
        fmt.Fprintln(out, ans)
    }
}
