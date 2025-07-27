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
        var x int
        fmt.Fscan(reader, &x)
        a[i] = int64(x)
    }

    pref := make([]int64, n+1)
    for i := 0; i < n; i++ {
        pref[i+1] = pref[i] + a[i]
    }

    const limit = 130
    count := 0
    for l := 0; l < n; l++ {
        for r := l + 2; r < n && r-l <= limit; r++ {
            sumInside := pref[r] - pref[l+1]
            if (a[l]^a[r]) == sumInside {
                count++
            }
        }
    }

    fmt.Println(count)
}
