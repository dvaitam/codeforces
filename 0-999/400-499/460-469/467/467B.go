package main

import (
    "bufio"
    "fmt"
    "math/bits"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    var n, m, k int
    // Read n (number of bits), m (number of other players), k (max differing bits)
    fmt.Fscan(reader, &n, &m, &k)
    armies := make([]int, m+1)
    for i := 0; i < m+1; i++ {
        fmt.Fscan(reader, &armies[i])
    }
    fedor := armies[m]
    friends := 0
    for i := 0; i < m; i++ {
        if bits.OnesCount(uint(fedor^armies[i])) <= k {
            friends++
        }
    }
    fmt.Println(friends)
}
