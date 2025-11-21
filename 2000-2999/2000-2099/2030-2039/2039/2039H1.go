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
    fmt.Fscan(in, &T)
    for ; T > 0; T-- {
        var n int
        fmt.Fscan(in, &n)
        a := make([]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &a[i])
        }
        sorted := make([]int, n)
        copy(sorted, a)
        sortInts(sorted)
        if equal(a, sorted) {
            fmt.Fprintln(out, 0)
            continue
        }
        fmt.Fprintln(out, 1)
        path := make([]byte, 0, 2*n-2)
        for i := 0; i < n-1; i++ {
            path = append(path, 'D')
        }
        for i := 0; i < n-1; i++ {
            path = append(path, 'R')
        }
        fmt.Fprintln(out, string(path))
    }
}

func equal(a, b []int) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}

func sortInts(a []int) {
    if len(a) <= 1 {
        return
    }
    quickSort(a, 0, len(a)-1)
}

func quickSort(a []int, l, r int) {
    for l < r {
        i, j := l, r
        pivot := a[(l+r)/2]
        for i <= j {
            for a[i] < pivot {
                i++
            }
            for a[j] > pivot {
                j--
            }
            if i <= j {
                a[i], a[j] = a[j], a[i]
                i++
                j--
            }
        }
        if j-l < r-i {
            if l < j {
                quickSort(a, l, j)
            }
            l = i
        } else {
            if i < r {
                quickSort(a, i, r)
            }
            r = j
        }
    }
}
