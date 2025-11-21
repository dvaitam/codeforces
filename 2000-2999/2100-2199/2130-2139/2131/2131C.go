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
    fmt.Fscan(in, &T)
    for ; T > 0; T-- {
        var n int
        fmt.Fscan(in, &n)
        s := make([]int, n)
        tArr := make([]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &s[i])
        }
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &tArr[i])
        }
        sort.Ints(s)
        sort.Ints(tArr)
        same := true
        for i := 0; i < n; i++ {
            if s[i] != tArr[i] {
                same = false
                break
            }
        }
        if same {
            fmt.Fprintln(out, "YES")
        } else {
            fmt.Fprintln(out, "NO")
        }
    }
}
