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
    if _, err := fmt.Fscan(in, &t); err != nil {
        return
    }
    for ; t > 0; t-- {
        var n int
        fmt.Fscan(in, &n)
        prefix := int64(0)
        seen := map[int64]bool{0: true}
        ok := false
        for i := 1; i <= n; i++ {
            var x int64
            fmt.Fscan(in, &x)
            if i%2 == 1 {
                prefix += x
            } else {
                prefix -= x
            }
            if seen[prefix] {
                ok = true
            }
            seen[prefix] = true
        }
        if ok {
            fmt.Fprintln(out, "YES")
        } else {
            fmt.Fprintln(out, "NO")
        }
    }
}
