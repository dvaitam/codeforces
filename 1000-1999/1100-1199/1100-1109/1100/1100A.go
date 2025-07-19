package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    var n, k int
    if _, err := fmt.Fscan(reader, &n, &k); err != nil {
        return
    }
    a := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &a[i])
    }
    maxAbs := 0
    for offset := 0; offset < k; offset++ {
        sum := 0
        for j := 0; j < n; j++ {
            if j%k == offset {
                continue
            }
            sum += a[j]
        }
        if sum < 0 {
            sum = -sum
        }
        if sum > maxAbs {
            maxAbs = sum
        }
    }
    fmt.Println(maxAbs)
}
