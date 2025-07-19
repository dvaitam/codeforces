package main

import "fmt"

const maxN = 1005

var f [maxN][maxN]float64
var flagArr [maxN][maxN]bool

// dfs computes the probability using memoized recursion
func dfs(x, y int) float64 {
    if x == 0 || y == 0 {
        return 1.0 / float64(y+1)
    }
    if flagArr[x][y] {
        return f[x][y]
    }
    flagArr[x][y] = true
    a := 1 - dfs(y, x-1)
    b := float64(y) / float64(y+1) * (1 - dfs(y-1, x))
    c := b + 1.0/float64(y+1)
    p := (c - b) / (1 - a - b + c)
    f[x][y] = p*(a-c) + c
    return f[x][y]
}

func main() {
    var n, m int
    if _, err := fmt.Scan(&n, &m); err != nil {
        return
    }
    res := dfs(n, m)
    fmt.Printf("%.10f %.10f\n", res, 1-res)
}
