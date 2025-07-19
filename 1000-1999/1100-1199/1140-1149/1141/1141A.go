package main

import "fmt"

func main() {
    var n, m int64
    if _, err := fmt.Scan(&n, &m); err != nil {
        return
    }
    const maxOps = 1<<60
    ans := int64(maxOps)
    for i := int64(0); i <= 30; i++ {
        for j := int64(0); j <= 30; j++ {
            temp := n
            for k := int64(0); k < i && temp < m; k++ {
                temp *= 2
            }
            for k := int64(0); k < j && temp < m; k++ {
                temp *= 3
            }
            if temp == m && i+j < ans {
                ans = i + j
            }
        }
    }
    if ans == maxOps {
        fmt.Println(-1)
    } else {
        fmt.Println(ans)
    }
}
