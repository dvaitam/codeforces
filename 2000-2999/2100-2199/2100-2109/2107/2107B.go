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

    var T int
    if _, err := fmt.Fscan(in, &T); err != nil {
        return
    }
    for ; T > 0; T-- {
        var n int
        var k int64
        fmt.Fscan(in, &n, &k)
        var mn int64 = 1<<63 - 1
        var mx int64
        var cMax int64
        var sum int64
        for i := 0; i < n; i++ {
            var v int64
            fmt.Fscan(in, &v)
            sum += v
            if v > mx {
                mx = v
                cMax = 1
            } else if v == mx {
                cMax++
            }
            if v < mn {
                mn = v
            }
        }
        diff := mx - mn
        win := false
        if diff > k+1 {
            win = false
        } else if diff == k+1 {
            if cMax >= 2 {
                win = false
            } else {
                win = (sum%2 == 1)
            }
        } else { // diff <= k
            win = (sum%2 == 1)
        }
        if win {
            fmt.Fprintln(out, "Tom")
        } else {
            fmt.Fprintln(out, "Jerry")
        }
    }
}

