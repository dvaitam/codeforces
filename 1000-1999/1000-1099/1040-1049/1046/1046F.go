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
    var x, k int64
    fmt.Fscan(reader, &x, &k)
    var ans int64
    t := x + k
    for _, ai := range a {
        quotient := ai / t
        ans += k * quotient
        rem := ai % t
        if rem > x {
            ans += k
        }
    }
    fmt.Println(ans)
}
