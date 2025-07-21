package main

import "fmt"

func main() {
    var n1, n2, k1, k2 int
    if _, err := fmt.Scan(&n1, &n2, &k1, &k2); err != nil {
        return
    }
    const mod = 100000000
    dpF := make([][]int, n1+1)
    dpH := make([][]int, n1+1)
    for i := 0; i <= n1; i++ {
        dpF[i] = make([]int, n2+1)
        dpH[i] = make([]int, n2+1)
    }
    dpF[0][0] = 1
    dpH[0][0] = 1
    for i := 0; i <= n1; i++ {
        for j := 0; j <= n2; j++ {
            if i == 0 && j == 0 {
                continue
            }
            if i > 0 {
                sum := 0
                for x := 1; x <= k1 && x <= i; x++ {
                    sum += dpH[i-x][j]
                    if sum >= mod {
                        sum -= mod
                    }
                }
                dpF[i][j] = sum
            }
            if j > 0 {
                sum := 0
                for y := 1; y <= k2 && y <= j; y++ {
                    sum += dpF[i][j-y]
                    if sum >= mod {
                        sum -= mod
                    }
                }
                dpH[i][j] = sum
            }
        }
    }
    ans := dpF[n1][n2] + dpH[n1][n2]
    ans %= mod
    fmt.Println(ans)
}
