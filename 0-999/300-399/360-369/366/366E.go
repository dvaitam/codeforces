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

    var n, m, k, s int
    fmt.Fscan(reader, &n, &m, &k, &s)
    // stats per note: index 1..k
    inf := int(1e9)
    maxS := make([]int, k+1)
    minS := make([]int, k+1)
    maxD := make([]int, k+1)
    minD := make([]int, k+1)
    for t := 1; t <= k; t++ {
        maxS[t] = -inf
        minS[t] = inf
        maxD[t] = -inf
        minD[t] = inf
    }
    // read guitar
    for i := 1; i <= n; i++ {
        for j := 1; j <= m; j++ {
            var t int
            fmt.Fscan(reader, &t)
            S := i + j
            D := i - j
            if S > maxS[t] {
                maxS[t] = S
            }
            if S < minS[t] {
                minS[t] = S
            }
            if D > maxD[t] {
                maxD[t] = D
            }
            if D < minD[t] {
                minD[t] = D
            }
        }
    }
    // read song
    Q := make([]int, s)
    for i := 0; i < s; i++ {
        fmt.Fscan(reader, &Q[i])
    }
    // precompute dist
    dist := make([][]int, k+1)
    for t := 0; t <= k; t++ {
        dist[t] = make([]int, k+1)
    }
    for t1 := 1; t1 <= k; t1++ {
        for t2 := 1; t2 <= k; t2++ {
            var d int
            if t1 == t2 {
                // max distance within same note
                d1 := maxS[t1] - minS[t1]
                d2 := maxD[t1] - minD[t1]
                if d2 > d1 {
                    d1 = d2
                }
                d = d1
            } else {
                // max over |S1 - S2| and |D1 - D2|
                d1 := maxS[t1] - minS[t2]
                if maxS[t2]-minS[t1] > d1 {
                    d1 = maxS[t2] - minS[t1]
                }
                d2 := maxD[t1] - minD[t2]
                if maxD[t2]-minD[t1] > d2 {
                    d2 = maxD[t2] - minD[t1]
                }
                if d2 > d1 {
                    d = d2
                } else {
                    d = d1
                }
            }
            dist[t1][t2] = d
        }
    }
    // compute answer
    ans := 0
    for i := 0; i+1 < s; i++ {
        t1 := Q[i]
        t2 := Q[i+1]
        if dist[t1][t2] > ans {
            ans = dist[t1][t2]
        }
    }
    fmt.Fprintln(writer, ans)
}
