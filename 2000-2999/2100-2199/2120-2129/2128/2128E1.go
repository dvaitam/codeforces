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

    var t int
    fmt.Fscan(in, &t)
    for ; t > 0; t-- {
        var n, k int
        fmt.Fscan(in, &n, &k)
        a := make([]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &a[i])
        }
        sorted := make([]int, n)
        copy(sorted, a)
        sort.Ints(sorted)

        bestV := sorted[0]
        bestL, bestR := 1, k
        for _, v := range sorted {
            found := false
            sum := make(map[int]int)
            sum[0] = 1
            cur := 0
            firstOcc := map[int]int{0: 0}
            lastOcc := map[int]int{}
            for i := 0; i < n; i++ {
                if a[i] >= v {
                    cur++
                } else {
                    cur--
                }
                if _, ok := firstOcc[cur]; !ok {
                    firstOcc[cur] = i + 1
                }
                lastOcc[cur] = i + 1
                if prev, ok := firstOcc[cur]; ok {
                    if lastOcc[cur]-prev >= k {
                        bestV = v
                        bestL = prev + 1
                        bestR = lastOcc[cur]
                        found = true
                        break
                    }
                }
            }
            if found {
                break
            }
        }
        fmt.Fprintf(out, "%d %d %d\n", bestV, bestL, bestR)
    }
}

