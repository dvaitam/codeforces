package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    var n int
    if _, err := fmt.Fscan(reader, &n); err != nil {
        return
    }
    home := make([]int, n)
    guest := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &home[i], &guest[i])
    }
    count := 0
    for i := 0; i < n; i++ {
        for j := 0; j < n; j++ {
            if home[i] == guest[j] {
                count++
            }
        }
    }
    fmt.Println(count)
}
