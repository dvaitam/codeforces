package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    var n int
    var m int64
    if _, err := fmt.Fscan(in, &n, &m); err != nil {
        return
    }
    vectors := make([][2]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(in, &vectors[i][0], &vectors[i][1])
    }
    // TODO: implement full solution
    fmt.Println(0)
}

