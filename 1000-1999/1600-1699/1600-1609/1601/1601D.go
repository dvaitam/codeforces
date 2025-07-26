package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var n int
    var d int
    if _, err := fmt.Fscan(reader, &n, &d); err != nil {
        return
    }
    type pair struct{ s, a int }
    arr := make([]pair, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &arr[i].s, &arr[i].a)
    }

    used := make([]bool, n)
    cur := d
    ans := 0
    for {
        best := -1
        bestA := 0
        for i, p := range arr {
            if used[i] {
                continue
            }
            if cur <= p.s {
                if best == -1 || p.a < bestA {
                    best = i
                    bestA = p.a
                }
            }
        }
        if best == -1 {
            break
        }
        used[best] = true
        if arr[best].a > cur {
            cur = arr[best].a
        }
        ans++
    }
    fmt.Fprintln(writer, ans)
}
