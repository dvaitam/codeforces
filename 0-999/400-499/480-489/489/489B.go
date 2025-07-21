package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    var n, m int
    fmt.Fscan(reader, &n)
    boys := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &boys[i])
    }
    fmt.Fscan(reader, &m)
    girls := make([]int, m)
    for i := 0; i < m; i++ {
        fmt.Fscan(reader, &girls[i])
    }
    sort.Ints(boys)
    sort.Ints(girls)
    i, j, ans := 0, 0, 0
    for i < n && j < m {
        if abs(boys[i]-girls[j]) <= 1 {
            ans++
            i++
            j++
        } else if boys[i] < girls[j] {
            i++
        } else {
            j++
        }
    }
    fmt.Println(ans)
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}
