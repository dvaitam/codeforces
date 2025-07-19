package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    var s string
    fmt.Fscan(reader, &s)
    n := len(s)
    k := n / 20
    if n%20 != 0 {
        k++
    }
    a := make([]int, k)
    total := 0
    for i := 0; i < k; i++ {
        a[i] = n / k
        total += a[i]
    }
    mx := n / k
    idxInc := 0
    for total < n {
        a[idxInc]++
        total++
        if a[idxInc] > mx {
            mx = a[idxInc]
        }
        idxInc++
        if idxInc >= k {
            idxInc = 0
        }
    }
    fmt.Println(k, mx)
    start := 0
    for i := 0; i < k; i++ {
        end := start + a[i]
        if end > n {
            end = n
        }
        part := s[start:end]
        start = end
        for len(part) < mx {
            part += "*"
        }
        fmt.Println(part)
    }
}
